// Copyright 2017 The cypherBFT Authors
// This file is part of the cypherBFT library.
//
// The cypherBFT library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The cypherBFT library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the cypherBFT library. If not, see <http://www.gnu.org/licenses/>.

// +build linux darwin

package cphash

/*
#cgo CFLAGS: -I./randomX
#cgo LDFLAGS: -L./randomX -lrandomx -lstdc++
#include <randomx.h>
*/
import "C"
import (
	"bytes"
	"encoding/binary"
	"errors"
	"math/big"
	"runtime"
	"unsafe"

	//"github.com/cypherium/cypherBFT-P/go-cypherium/common/math"

	"github.com/cypherium/cypherBFT-P/go-cypherium/core/types"
	"github.com/cypherium/cypherBFT-P/go-cypherium/log"
	"github.com/cypherium/cypherBFT-P/go-cypherium/params"
	//set "gopkg.in/fatih/set.v0"
)

// Various error messages to mark blocks invalid. These should be private to
// prevent engine specific errors from being referenced in the remainder of the
// codebase, inherently breaking if the engine is swapped out. Please put common
// error types into the pow package.
var (
	errLargeBlockTime    = errors.New("timestamp too big")
	errZeroBlockTime     = errors.New("timestamp equals parent's")
	errTooManyUncles     = errors.New("too many uncles")
	errDuplicateUncle    = errors.New("duplicate uncle")
	errUncleIsAncestor   = errors.New("uncle is ancestor")
	errDanglingUncle     = errors.New("uncle's parent is not ancestor")
	errInvalidDifficulty = errors.New("non-positive difficulty")
	errInvalidMixDigest  = errors.New("invalid mix digest")
	errInvalidPoW        = errors.New("invalid proof-of-work")
)

// CalcKeyBlockDifficulty is the difficulty adjustment algorithm. It returns
// the difficulty that a new block should have when created at time
// given the parent block's time and difficulty.
func (cphash *Cphash) CalcKeyBlockDifficulty(chain types.KeyChainReader, time uint64, parent *types.KeyBlockHeader) *big.Int {
	return calcKeyBlockDifficultyByzantium(time, parent)
}

const (
	MinFoundedSeconds = 5
)

// Some weird constants to avoid constant memory allocs for them.
var (
	expDiffPeriod = big.NewInt(1000)
	big1          = big.NewInt(1)
	big2          = big.NewInt(2)
	big9          = big.NewInt(9)
	big10         = big.NewInt(10)
	bigMinus99    = big.NewInt(-99)
	big2999999    = big.NewInt(2999999)
)

func calcKeyBlockDifficultyByzantium(time uint64, parent *types.KeyBlockHeader) *big.Int {
	// https://github.com/cypherium/EIPs/issues/100.
	// algorithm:
	// diff = (parent_diff +
	//         (parent_diff / 2048 * max((2 if len(parent.uncles) else 1) - ((timestamp - parent.timestamp) // 9), -99))
	//        ) + 2^(periodCount - 2)

	bigTime := new(big.Int).SetUint64(time)
	bigParentTime := new(big.Int).Set(parent.Time)

	// holds intermediate values to make the algo easier to read & audit
	x := new(big.Int)
	y := new(big.Int)

	// (2 if len(parent_uncles) else 1) - (block_timestamp - parent_timestamp) // 9
	x.Sub(bigTime, bigParentTime)
	x.Div(x, big9)
	x.Sub(big1, x)

	// max((2 if len(parent_uncles) else 1) - (block_timestamp - parent_timestamp) // 9, -99)
	if x.Cmp(bigMinus99) < 0 {
		x.Set(bigMinus99)
	}
	// parent_diff + (parent_diff / 2048 * max((2 if len(parent.uncles) else 1) - ((timestamp - parent.timestamp) // 9), -99))
	y.Div(parent.Difficulty, params.DifficultyBoundDivisor)
	x.Mul(y, x)
	x.Add(parent.Difficulty, x)

	// minimum difficulty can ever be (before exponential factor)
	if x.Cmp(params.MinimumDifficulty) < 0 {
		x.Set(params.MinimumDifficulty)
	}
	// calculate a fake block number for the ice-age delay:
	//   https://github.com/cypherium/EIPs/pull/669
	//   fake_block_number = min(0, block.number - 3_000_000
	fakeBlockNumber := new(big.Int)
	if parent.Number.Cmp(big2999999) >= 0 {
		fakeBlockNumber = fakeBlockNumber.Sub(parent.Number, big2999999) // Note, parent is 1 less than the actual block number
	}
	// for the exponential factor
	periodCount := fakeBlockNumber
	periodCount.Div(periodCount, expDiffPeriod)

	// the exponential factor, commonly referred to as "the bomb"
	// diff = diff + 2^(periodCount - 2)
	if periodCount.Cmp(big1) > 0 {
		y.Sub(periodCount, big2)
		y.Exp(big2, y, nil)
		x.Add(x, y)
	}
	return x
}

func (cphash *Cphash) verifyRangeCandidate(chain types.KeyChainReader, candidate *types.Candidate) error {
	// If we're running a fake PoW, accept any seal as valid
	if cphash.config.PowMode == ModeFake || cphash.config.PowMode == ModeFullFake {
		// time.Sleep(cphash.fakeDelay)
		if cphash.fakeFail == candidate.KeyCandidate.Number.Uint64() {
			return errInvalidPoW
		}
		return nil
	}

	// If we're running a shared PoW, delegate verification to it
	if cphash.shared != nil {
		return cphash.shared.VerifyCandidate(chain, candidate)
	}
	// Ensure that we have a valid difficulty for the block
	if candidate.KeyCandidate.Difficulty.Sign() <= 0 {
		return errInvalidDifficulty
	}
	// Recompute the digest and PoW value and verify against the header
	number := candidate.KeyCandidate.Number.Uint64()

	cache := cphash.cache(number)
	size := datasetSize(number)
	if cphash.config.PowMode == ModeTest {
		size = 32 * 1024
	}
	digest, result := hashimotoLight(size, cache.cache, candidate.HashNoNonce().Bytes(), candidate.KeyCandidate.Nonce.Uint64())
	// Caches are unmapped in a finalizer. Ensure that the cache stays live
	// until after the call to hashimotoLight so it's not unmapped while being used.
	runtime.KeepAlive(cache)

	if !bytes.Equal(candidate.KeyCandidate.MixDigest[:], digest) {
		return errInvalidMixDigest
	}

	target := new(big.Int).Div(maxUint256, candidate.KeyCandidate.Difficulty)
	if new(big.Int).SetBytes(result).Cmp(target) > 0 {
		return errInvalidPoW
	}
	return nil
}

////////////////////////////////////////////////////////////////////////////////////////////////
// VerifyCandidate implements pow.Engine, checking whether the given candidate satisfies
// the PoW difficulty requirements.
func (cphash *Cphash) verifyCandidate(chain types.KeyChainReader, candidate *types.Candidate) error {
	// If we're running a fake PoW, accept any seal as valid
	if cphash.config.PowMode == ModeFake || cphash.config.PowMode == ModeFullFake {
		// time.Sleep(cphash.fakeDelay)
		if cphash.fakeFail == candidate.KeyCandidate.Number.Uint64() {
			return errInvalidPoW
		}
		return nil
	}
	defer func() {
		log.Debug("verify candidate end")
	}()

	log.Debug("verify candidate")

	// If we're running a shared PoW, delegate verification to it
	//if cphash.shared != nil {
	//	return cphash.shared.VerifyCandidate(chain, candidate)
	//}
	// Ensure that we have a valid difficulty for the block
	if candidate.KeyCandidate.Difficulty.Sign() <= 0 {
		return errInvalidDifficulty
	}
	// Recompute the digest and PoW value and verify against the header
	//number := candidate.KeyCandidate.Number.Uint64()

	//cache := cphash.cache(number)
	//size := datasetSize(number)
	//if cphash.config.PowMode == ModeTest {
	//	size = 32 * 1024
	//}
	//digest, result := hashimotoLight(size, cache.cache, candidate.HashNoNonce().Bytes(), candidate.KeyCandidate.Nonce.Uint64())
	cphash.lock.Lock()

	var result [C.RANDOMX_HASH_SIZE]byte

	hash := candidate.HashNoNonce().Bytes()
	inputWithNoce := make([]byte, len(hash)+8)
	copy(inputWithNoce[0:], hash)

	nonce := candidate.KeyCandidate.Nonce.Uint64()

	binary.LittleEndian.PutUint64(inputWithNoce[len(hash):], uint64(nonce))
	C.randomx_calculate_hash(cphash.vvm, unsafe.Pointer(&inputWithNoce[0]), (C.ulong)(len(inputWithNoce)), unsafe.Pointer(&result[0]))

	cphash.lock.Unlock()
	// Caches are unmapped in a finalizer. Ensure that the cache stays live
	// until after the call to hashimotoLight so it's not unmapped while being used.
	//runtime.KeepAlive(cache)

	if !bytes.Equal(candidate.KeyCandidate.MixDigest[:], result[:]) {
		return errInvalidMixDigest
	}

	target := new(big.Int).Div(maxUint256, candidate.KeyCandidate.Difficulty)
	if new(big.Int).SetBytes(result[:]).Cmp(target) > 0 {
		return errInvalidPoW
	}
	return nil
}

func (cphash *Cphash) VerifyCandidate(chain types.KeyChainReader, candidate *types.Candidate) error {
	if cphash.config.PowRangeMode == 0 {
		return cphash.verifyCandidate(chain, candidate)
	} else {
		return cphash.verifyRangeCandidate(chain, candidate)
	}
}

// Prepare implements pow.Engine, initializing the difficulty field of a
// candidate to conform to the cphash protocol. The changes are done inline.
func (cphash *Cphash) PrepareRangeCandidate(chain types.KeyChainReader, candidate *types.Candidate, committeeSize int) error {
	log.Debug("prepare candidate from header", "hash", candidate.KeyCandidate.ParentHash, "number", candidate.KeyCandidate.Number.Uint64()-1)

	parent := chain.GetHeader(candidate.KeyCandidate.ParentHash, candidate.KeyCandidate.Number.Uint64()-1)
	if parent == nil {
		return types.ErrUnknownAncestor
	}

	candidate.KeyCandidate.Difficulty = calcRangeCandidateDifficulty(candidate.KeyCandidate.Time.Uint64(), parent, committeeSize)
	log.Info("PrepareCandidate", "parent difficulty", parent.Difficulty, "current difficulty", candidate.KeyCandidate.Difficulty, "minus value", candidate.KeyCandidate.Difficulty.Int64()-parent.Difficulty.Int64(), "committeeSize", committeeSize)
	return nil
}

// calcCandidateDifficulty is the difficulty adjustment algorithm. It returns
// the difficulty that a new candidate should have when created at time
// given the keyblock's time and difficulty.
func calcRangeCandidateDifficulty(time uint64, parent *types.KeyBlockHeader, committeeSize int) *big.Int {
	// algorithm:
	// diff = (parent_diff +
	//         (parent_diff / 2048 * max(1 - (block_timestamp - parent_timestamp) // 10, -99))
	//        ) + 2^(periodCount - 2)

	bigTime := new(big.Int).SetUint64(time)
	bigParentTime := new(big.Int).Set(parent.Time)

	// holds intermediate values to make the algo easier to read & audit
	x := new(big.Int)
	y := new(big.Int)

	// 1 - (block_timestamp - parent_timestamp) // 10
	x.Sub(bigTime, bigParentTime)
	x.Div(x, big.NewInt(50))
	x.Sub(big1, x)

	// max(1 - (block_timestamp - parent_timestamp) // 10, -99)
	if x.Cmp(bigMinus99) < 0 {
		x.Set(bigMinus99)
	}
	// (parent_diff + (parent_diff // 2048) * max(1 - (block_timestamp - parent_timestamp) // 10, -99))
	y.Div(parent.Difficulty, params.DifficultyBoundDivisor)
	x.Mul(y, x)
	x.Add(parent.Difficulty, x)

	// minimum difficulty can ever be (before exponential factor)
	if x.Cmp(params.MinimumDifficulty) < 0 {
		x.Set(params.MinimumDifficulty)
	}
	// for the exponential factor
	periodCount := big.NewInt(int64(committeeSize))
	periodCount.Div(periodCount, expDiffPeriod)

	// the exponential factor, commonly referred to as "the bomb"
	// diff = diff + 2^(periodCount - 2)
	if periodCount.Cmp(big1) > 0 {
		y.Sub(periodCount, big2)
		y.Exp(big2, y, nil)
		x.Add(x, y)
	}
	return x
}

// Prepare implements pow.Engine, initializing the difficulty field of a
// candidate to conform to the cphash protocol. The changes are done inline.
func (cphash *Cphash) prepareCandidate(chain types.KeyChainReader, candidate *types.Candidate, committeeSize int) error {
	log.Debug("prepare candidate ", "hash", candidate.KeyCandidate.ParentHash, "number", candidate.KeyCandidate.Number.Uint64())

	parent := chain.GetHeader(candidate.KeyCandidate.ParentHash, candidate.KeyCandidate.Number.Uint64()-1)
	if parent == nil {
		return types.ErrUnknownAncestor
	}

	//candidate.KeyCandidate.Difficulty = calcCandidateDifficulty(candidate.KeyCandidate.Time.Uint64(), parent, big.NewInt(3000000))
	candidate.KeyCandidate.Difficulty = new(big.Int).SetUint64(0x5ff00)
	//log.Info("PrepareCandidate", "parent difficulty", parent.Difficulty, "current difficulty", candidate.KeyCandidate.Difficulty, "minus value", candidate.KeyCandidate.Difficulty.Int64()-parent.Difficulty.Int64(), "committeeSize", committeeSize)
	return nil
}

// calcCandidateDifficulty is the difficulty adjustment algorithm. It returns
// the difficulty that a new candidate should have when created at time
// given the keyblock's time and difficulty.
func calcCandidateDifficulty(time uint64, parent *types.KeyBlockHeader, bombDelay *big.Int) *big.Int {
	// algorithm:
	// diff = (parent_diff +
	//         (parent_diff / 2048 * max(1 - (block_timestamp - parent_timestamp) // 10, -99))
	//        ) + 2^(periodCount - 2)

	bigTime := new(big.Int).SetUint64(time)
	bigParentTime := new(big.Int).Set(parent.Time)

	// holds intermediate values to make the algo easier to read & audit
	x := new(big.Int)
	y := new(big.Int)

	// 1 - (block_timestamp - parent_timestamp) // 10
	x.Sub(bigTime, bigParentTime)
	x.Div(x, big.NewInt(50))
	x.Sub(big1, x)

	// max(1 - (block_timestamp - parent_timestamp) // 10, -99)
	if x.Cmp(bigMinus99) < 0 {
		x.Set(bigMinus99)
	}
	// (parent_diff + (parent_diff // 2048) * max(1 - (block_timestamp - parent_timestamp) // 10, -99))
	y.Div(parent.Difficulty, params.DifficultyBoundDivisor)
	x.Mul(y, x)
	x.Add(parent.Difficulty, x)

	// minimum difficulty can ever be (before exponential factor)
	if x.Cmp(params.MinimumDifficulty) < 0 {
		x.Set(params.MinimumDifficulty)
	}

	bombDelayFromParent := new(big.Int).Sub(bombDelay, big1)

	// calculate a fake block number for the ice-age delay

	fakeBlockNumber := new(big.Int)
	if parent.Number.Cmp(bombDelayFromParent) >= 0 {
		fakeBlockNumber = fakeBlockNumber.Sub(parent.Number, bombDelayFromParent)
	}
	// for the exponential factor
	periodCount := fakeBlockNumber
	periodCount.Div(periodCount, expDiffPeriod)
	//
	//// for the exponential factor
	//periodCount := big.NewInt(int64(committeeSize))
	//periodCount.Div(periodCount, expDiffPeriod)

	// the exponential factor, commonly referred to as "the bomb"
	// diff = diff + 2^(periodCount - 2)
	if periodCount.Cmp(big1) > 0 {
		y.Sub(periodCount, big2)
		y.Exp(big2, y, nil)
		x.Add(x, y)
	}
	return x
}

func (cphash *Cphash) PrepareCandidate(chain types.KeyChainReader, candidate *types.Candidate, committeeSize int) error {
	if cphash.config.PowRangeMode == 0 {
		return cphash.prepareCandidate(chain, candidate, committeeSize)
	} else {
		return cphash.PrepareRangeCandidate(chain, candidate, committeeSize)
	}
}

func (cphash *Cphash) PowMode() uint {
	return uint(cphash.config.PowMode)
}

func (cphash *Cphash) PowRangeMode() uint {
	return uint(cphash.config.PowRangeMode)
}
