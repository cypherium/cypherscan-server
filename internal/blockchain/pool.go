package blockchain

import (
	"context"
	"sync"
	"time"

	"github.com/cypherium/cypherscan-server/internal/util"
	log "github.com/sirupsen/logrus"
)

// BlockChainClientPooledObject is the client stored in pool
type BlockChainClientPooledObject struct {
	Client *BlockChain
	Err    *util.MyError
}

// Pool is the object pool
type Pool struct {
	ctx             context.Context
	maxSize         int
	currentSize     int
	inUseSize       int
	buffer          chan *BlockChainClientPooledObject
	borrowTimeoutMs int
	mux             *sync.Mutex
	nodesUrls       []string
	cond            *sync.Cond
}

// NewPoolOptions is the options when new a pool
type NewPoolOptions struct {
	MaxSize int
	// BorrowTimeout is the timeout when try to borrow object from pool
	BorrowTimeoutMs int
	NodesUrls       []string
}

// NewPool is to create the pool
func NewPool(ctx context.Context, options *NewPoolOptions) (*Pool, chan interface{}) {
	c := make(chan *BlockChainClientPooledObject, options.MaxSize)
	mux := &sync.Mutex{}
	cond := sync.NewCond(mux)
	pool := &Pool{
		ctx:             ctx,
		buffer:          c,
		maxSize:         options.MaxSize,
		nodesUrls:       options.NodesUrls,
		borrowTimeoutMs: options.BorrowTimeoutMs,
		inUseSize:       0,
		mux:             mux,
		cond:            cond,
	}
	poolTerminated := make(chan interface{})
	go func() {
		select {
		case <-ctx.Done():
			cond.L.Lock()
			for pool.inUseSize > 0 {
				cond.Wait()
			}
			cond.L.Unlock()
			close(pool.buffer)
			for client := range pool.buffer {
				client.Client.Close()
			}
			close(poolTerminated)
		}
	}()
	return pool, poolTerminated
}

// Borrow is to borrow an object from the
func (p *Pool) Borrow() (*BlockChainClientPooledObject, error) {
	for {
		if p.currentSize < p.maxSize {
			select {
			case client := <-p.buffer:
				if client.Err != nil {
					client.Client.Close()
					p.mux.Lock()
					p.currentSize--
					p.mux.Unlock()
					break
				}
				p.mux.Lock()
				p.inUseSize++
				p.mux.Unlock()
				return client, nil
			default:
				client, err := Dial(p.ctx, p.nodesUrls[0])
				log.Infof("Connected to blockchain nodes, %d of %d", p.currentSize, p.maxSize)
				if err != nil {
					return nil, util.NewError(err, "Connect to node failed when borrowing")
				}
				p.buffer <- &BlockChainClientPooledObject{client, nil}
				p.mux.Lock()
				p.currentSize++
				p.mux.Unlock()
			}
		} else {
			select {
			case client := <-p.buffer:
				if client.Err != nil {
					client.Client.Close()
					p.mux.Lock()
					p.currentSize--
					p.mux.Unlock()
					break
				}
				p.mux.Lock()
				p.inUseSize++
				p.mux.Unlock()
				return client, nil
			case <-time.NewTimer(time.Duration(p.borrowTimeoutMs) * time.Millisecond).C:
				return nil, util.NewError(nil, "Borrow timeout")
			}
		}
	}
}

// Return will return the borrowed object back to pool
func (p *Pool) Return(client *BlockChainClientPooledObject) {
	p.buffer <- client
	p.mux.Lock()
	p.inUseSize--
	p.mux.Unlock()
	p.cond.Signal()
}
