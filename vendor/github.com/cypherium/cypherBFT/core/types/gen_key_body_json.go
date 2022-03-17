// Code generated by github.com/fjl/gencodec. DO NOT EDIT.

package types

import (
	"encoding/json"
	"errors"

	"github.com/cypherium/cypherBFT/common/hexutil"
)

var _ = (*keyBlockBodyMarshaling)(nil)

// MarshalJSON marshals as JSON.
func (k KeyBlockBody) MarshalJSON() ([]byte, error) {
	type KeyBlockBody struct {
		LeaderPubKey  string `json:"leaderPubKey"           gencodec:"required"`
		LeaderAddress string
		InPubKey      string        `json:"inPubKey"            	gencodec:"required"`
		InAddress     string        `json:"inAddress"            gencodec:"required"`
		OutPubKey     string        `json:"outPubKey"            	gencodec:"required"`
		OutAddress    string        `json:"outAddress"            gencodec:"required"`
		Signatrue     hexutil.Bytes `json:"signature"`
		Exceptions    hexutil.Bytes `json:"exceptions"`
	}
	var enc KeyBlockBody
	enc.Signatrue = k.Signatrue
	enc.Exceptions = k.Exceptions
	enc.InPubKey = k.InPubKey
	enc.InAddress = k.InAddress
	enc.LeaderPubKey = k.LeaderPubKey
	enc.LeaderAddress = k.LeaderAddress
	enc.OutPubKey = k.OutPubKey
	enc.OutAddress = k.OutAddress
	return json.Marshal(&enc)
}

// UnmarshalJSON unmarshals from JSON.
func (k *KeyBlockBody) UnmarshalJSON(input []byte) error {
	type KeyBlockBody struct {
		LeaderPubKey  *string `json:"leaderPubKey"           gencodec:"required"`
		LeaderAddress *string
		InPubKey      *string        `json:"inPubKey"            	gencodec:"required"`
		InAddress     *string        `json:"inAddress"            gencodec:"required"`
		OutPubKey     *string        `json:"outPubKey"            	gencodec:"required"`
		OutAddress    *string        `json:"outAddress"            gencodec:"required"`
		Signatrue     *hexutil.Bytes `json:"signature"`
		Exceptions    *hexutil.Bytes `json:"exceptions"`
	}
	var dec KeyBlockBody
	if err := json.Unmarshal(input, &dec); err != nil {
		return err
	}
	if dec.Signatrue != nil {
		k.Signatrue = *dec.Signatrue
	}
	if dec.Exceptions != nil {
		k.Exceptions = *dec.Exceptions
	}
	if dec.InPubKey != nil {
		k.InPubKey = *dec.InPubKey
	}
	if dec.InAddress == nil {
		return errors.New("missing required field 'inAddress' for KeyBlockBody")
	}
	k.InAddress = *dec.InAddress
	if dec.LeaderPubKey == nil {
		return errors.New("missing required field 'leaderPubKey' for KeyBlockBody")
	}
	k.LeaderPubKey = *dec.LeaderPubKey
	if dec.LeaderAddress == nil {
		return errors.New("missing required field 'LeaderAddress' for KeyBlockBody")
	}
	k.LeaderAddress = *dec.LeaderAddress

	if dec.OutPubKey != nil {
		k.OutPubKey = *dec.OutPubKey
	}
	if dec.OutAddress == nil {
		return errors.New("missing required field 'outAddress' for KeyBlockBody")
	}
	k.OutAddress = *dec.OutAddress
	return nil
}