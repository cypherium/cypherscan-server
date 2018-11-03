package home

import (
  "context"
  "fmt"
  "github.com/ethereum/go-ethereum/rpc"
  "github.com/jinzhu/gorm"
  "gitlab.com/ron-liu/cypherscan-server/internal/env"
  "gitlab.com/ron-liu/cypherscan-server/internal/txblock"
  "gitlab.com/ron-liu/cypherscan-server/internal/util"
  // "math/big"
  "time"
)

type blockInfo struct {
  ParentHash   string
  UncleHash    string
  Coinbase     string
  Root         string
  TxHash       string
  ReceiptHash  string
  Bloom        string
  Difficulty   string
  Number       string
  GasLimit     string
  GasUsed      string
  Time         string
  Extra        string
  MixDigest    string
  Nonce        string
  Transactions []transactionInfo
}

type transactionInfo struct {
  AccountNonce     string
  Price            string
  GasLimit         string
  Recipient        string
  Amount           string
  Payload          string
  V                string
  R                string
  S                string
  Hash             string
  BlockNumber      string
  BlockHash        string
  From             string
  TransactionIndex string
}

// SubscribeNewBlock is to subscribe new block
func SubscribeNewBlock() {
  client, err := rpc.Dial(env.Env.TsBlockChainWsURL)
  if err != nil {
    fmt.Println(err)
    return
  }
  subch := make(chan blockInfo)

  // Ensure that subch receives the latest block.
  go func() {
    for {
      subscribeBlocks(client, subch)
      time.Sleep(2 * time.Second)
    }
  }()

  // Print events from the subscription as they arrive.
  for block := range subch {
    fmt.Printf("pre block: %+v\n", block)
    txBlock := transformToTxBlock(block)
    fmt.Printf("post block: %#v", txBlock)
    util.Run(func(db *gorm.DB) error {
      db.NewRecord(txBlock)
      db.Create(&txBlock)
      return nil
    })
    // fmt.Println("latest block:", block.Number, "diff", block.Difficulty, "txHash", block.TxHash, "gas limit", block.GasLimit, "gs used", block.GasUsed, block.Nonce)
  }
}

func transformToTxBlock(b blockInfo) *txblock.TxBlock {
  return &txblock.TxBlock{
    Number:      util.Parse(b.Number, util.UInt64Type).(uint64),
    ParentHash:  util.Parse(b.ParentHash, util.BytesType).([]byte),
    UncleHash:   util.Parse(b.UncleHash, util.BytesType).([]byte),
    Coinbase:    util.Parse(b.Coinbase, util.BytesType).([]byte),
    Root:        util.Parse(b.Root, util.BytesType).([]byte),
    TxHash:      util.Parse(b.TxHash, util.BytesType).([]byte),
    ReceiptHash: util.Parse(b.ReceiptHash, util.BytesType).([]byte),
    Bloom:       util.Parse(b.Bloom, util.BytesType).([]byte),
    Difficulty:  util.Parse(b.Difficulty, util.BytesType).([]byte),
    GasLimit:    util.Parse(b.GasLimit, util.BytesType).([]byte),
    GasUsed:     util.Parse(b.GasUsed, util.BytesType).([]byte),
    // // Time:        util.Parse(b.Time, util.TimeType).(time.Time),
    // // Extra       string
    MixDigest: util.Parse(b.MixDigest, util.BytesType).([]byte),
    Nonce:     util.Parse(b.Nonce, util.BytesType).([]byte),
    Transactions: func(ts []transactionInfo) []txblock.Transaction {
      transactions := make([]txblock.Transaction, len(ts))
      for i, t := range b.Transactions {
        transactions[i] = txblock.Transaction{
          AccountNonce:     util.Parse(t.AccountNonce, util.UInt64Type).(uint64),
          Price:            util.Parse(t.Price, util.BytesType).([]byte),
          GasLimit:         util.Parse(t.GasLimit, util.UInt64Type).(uint64),
          Recipient:        util.Parse(t.Recipient, util.BytesType).([]byte),
          Amount:           util.Parse(t.Amount, util.BytesType).([]byte),
          Payload:          util.Parse(t.Payload, util.BytesType).([]byte),
          V:                util.Parse(t.V, util.BytesType).([]byte),
          R:                util.Parse(t.R, util.BytesType).([]byte),
          S:                util.Parse(t.S, util.BytesType).([]byte),
          Hash:             util.Parse(t.Hash, util.BytesType).([]byte),
          BlockNumber:      util.Parse(t.BlockNumber, util.UInt64Type).(uint64),
          BlockHash:        util.Parse(t.BlockHash, util.BytesType).([]byte),
          From:             util.Parse(t.From, util.BytesType).([]byte),
          TransactionIndex: util.Parse(t.TransactionIndex, util.UInt32Type).(uint32),
        }
      }
      return transactions
    }(b.Transactions),
  }
}

func subscribeBlocks(client *rpc.Client, subch chan blockInfo) {
  ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
  defer cancel()

  // Subscribe to new blocks.
  sub, err := client.EthSubscribe(ctx, subch, "newHeads")
  if err != nil {
    fmt.Println("subscribe error:", err)
    return
  }

  // The connection is established now.
  // Update the channel with the current block.
  var lastBlock blockInfo
  if err := client.CallContext(ctx, &lastBlock, "eth_getBlockByNumber", "latest", true); err != nil {
    fmt.Println("can't get latest block:", err)
    return
  }
  subch <- lastBlock

  // The subscription will deliver events to the channel. Wait for the
  // subscription to end for any reason, then loop around to re-establish
  // the connection.
  fmt.Println("connection lost: ", <-sub.Err())
}

/* data recieved from rpc
block:
map[
  parentHash:0x676f8750affdc977b9fd5e2729966167d6ac985f306ed8282edca152f4cc1294
  extraData:0x73656f34
  gasUsed:0x7a0ddc
  miner:0xb2930b35844a230f00e51431acae96fe543a0347
  mixHash:0x2ab7b824cb0c2897a3c30584bf756484d63c99965e43195a64e256a07a3560a6
  logsBloom:0x5380010534c080413142c84020520810d68e12886904000a005442484320017a41328048001024220080068120822c25063001c1c9616ad066d6010084240000c1c010332210412d0a02ab184848520274604c594011a12048102280080c0013101444c846418080044c1940d2805801700130041d501242062240195082ac6410408054868080706220405340f5020b1a4731c162009004c3040006ab0082011b546c81210092b844321f2c96888482b036a442a94032a905b0a01a06400080140ab89318988900a290000314c0034640114012811b39000020033158f026920111350580800422008318a6486aa01d100800081a00a1b50811021884810182
  stateRoot:0xa52ebb5e5328ff89587ece3d30799d6093ac4e66016e4069240782ac168421a6
  difficulty:0xac0364880c423
  gasLimit:0x7a121d
  hash:0xf0688febaff9c24ba6234b5e2687d9be4b416a9b0d22c22e7864ab33413d8597
  nonce:0x47526d9c0a7b56b1
  number:0x6530fa
  transactionsRoot:0xffdb2b0e3cd21165916b441e6477d9ce075b685722bf474b200035ba808fe08f
  timestamp:0x5bdcb1a1
  receiptsRoot:0xcecafe1acd3b5f44f6ebd2163e959cf95b4a324c3e298a9b9e9afc483c02551c

  uncles:[]
  totalDifficulty:0x19b5186b4063fd60e07
  sha3Uncles:0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347
  size:0x3f86
  transactions:[
    map[
      hash:0xaf45dac06635e12b88aa5752a25d454dfea038100ecb502cf70e416079e03b1d
      v:0x26
      from:0x3f4c913d6ac0926941b849d6d9305d5284fae899
      to:0x5bfe5f416f7ea1d0e95601700476ad140b6b7490
      transactionIndex:0x0
      value:0x24f6f618cc1800
      blockHash:0xf0688febaff9c24ba6234b5e2687d9be4b416a9b0d22c22e7864ab33413d8597
      blockNumber:0x6530fa
      gas:0x5208
      nonce:0x5
      gasPrice:0x98bca5a00
      input:0x
      r:0x1023d243e032bf5a805b2f845c10b34ef4a03452cf4d9ee030d25be18148784f
      s:0x4e57aad06427421370a3817290dc1adc4b8a14f903415a6ab9be2131825525c0
    ]
  ]
]
*/
