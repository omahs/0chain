package miner

import (
	"context"
	"sync"

	"0chain.net/block"
	"0chain.net/chain"
	"0chain.net/common"
	"0chain.net/datastore"
	"0chain.net/memorystore"
	"0chain.net/node"
	"0chain.net/round"
)

//RoundMismatch - to indicate an error where the current round and the given round don't match
const RoundMismatch = "round_mismatch"

//ErrRoundMismatch - an error object for mismatched round error
var ErrRoundMismatch = common.NewError(RoundMismatch, "Current round number of the chain doesn't match the block generation round")

var minerChain = &Chain{}

/*SetupMinerChain - setup the miner's chain */
func SetupMinerChain(c *chain.Chain) {
	minerChain.Chain = *c
	minerChain.rounds = make(map[int64]*Round)
	minerChain.roundsMutex = &sync.Mutex{}
	minerChain.BlockMessageChannel = make(chan *BlockMessage, 25)
}

/*GetMinerChain - get the miner's chain */
func GetMinerChain() *Chain {
	return minerChain
}

/*Chain - A miner chain to manage the miner activities */
type Chain struct {
	chain.Chain
	BlockMessageChannel chan *BlockMessage
	roundsMutex         *sync.Mutex
	rounds              map[int64]*Round
	DiscoverClients     bool
}

/*GetBlockMessageChannel - get the block messages channel */
func (mc *Chain) GetBlockMessageChannel() chan *BlockMessage {
	return mc.BlockMessageChannel
}

/*SetupGenesisBlock - setup the genesis block for this chain */
func (mc *Chain) SetupGenesisBlock(hash string) *block.Block {
	gr, gb := mc.GenerateGenesisBlock(hash)
	if gr == nil || gb == nil {
		panic("Genesis round/block canot be null")
	}
	mgr := mc.CreateRound(gr)
	mc.AddRound(mgr)
	mc.AddGenesisBlock(gb)
	return gb
}

/*CreateRound - create a round */
func (mc *Chain) CreateRound(r *round.Round) *Round {
	var mr Round
	r.ComputeRanks(mc.Miners.Size(), mc.Sharders.Size())
	mr.Round = r
	mr.blocksToVerifyChannel = make(chan *block.Block, mc.NumGenerators)
	return &mr
}

/*AddRound - Add Round to the block */
func (mc *Chain) AddRound(r *Round) bool {
	mc.roundsMutex.Lock()
	defer mc.roundsMutex.Unlock()
	_, ok := mc.rounds[r.Number]
	if ok {
		return false
	}
	r.ComputeRanks(mc.Miners.Size(), mc.Sharders.Size())
	mc.rounds[r.Number] = r
	if r.Number > mc.CurrentRound {
		mc.CurrentRound = r.Number
	}
	return true
}

/*SetLatestFinalizedBlock - Set latest finalized block */
func (mc *Chain) SetLatestFinalizedBlock(ctx context.Context, b *block.Block) {
	mc.AddBlock(b)
	mc.LatestFinalizedBlock = b
	var r *round.Round = datastore.GetEntityMetadata("round").Instance().(*round.Round)
	r.Number = b.Round
	r.RandomSeed = b.RoundRandomSeed
	mr := mc.CreateRound(r)
	mc.AddRound(mr)
	mc.AddNotarizedBlock(ctx, r, b)
}

/*GetRound - get a round */
func (mc *Chain) GetRound(roundNumber int64) *Round {
	mc.roundsMutex.Lock()
	defer mc.roundsMutex.Unlock()
	round, ok := mc.rounds[roundNumber]
	if !ok {
		return nil
	}
	return round
}

/*DeleteRound - delete a round and associated block data */
func (mc *Chain) DeleteRound(ctx context.Context, r *round.Round) {
	mc.roundsMutex.Lock()
	defer mc.roundsMutex.Unlock()
	delete(mc.rounds, r.Number)
}

/*DeleteRoundsBelow - delete rounds below */
func (mc *Chain) DeleteRoundsBelow(ctx context.Context, round int64) {
	mc.roundsMutex.Lock()
	defer mc.roundsMutex.Unlock()
	rounds := make([]*Round, 0, 1)
	for _, r := range mc.rounds {
		if r.Number < round {
			rounds = append(rounds, r)
		}
	}
	for _, r := range rounds {
		r.Clear()
		delete(mc.rounds, r.Number)
	}
}

/*CancelRoundsBelow - delete rounds below */
func (mc *Chain) CancelRoundsBelow(ctx context.Context, round int64) {
	mc.roundsMutex.Lock()
	defer mc.roundsMutex.Unlock()
	for _, r := range mc.rounds {
		if r.Number < round {
			r.CancelVerification()
		}
	}
}

func (mc *Chain) deleteTxns(txns []datastore.Entity) error {
	transactionMetadataProvider := datastore.GetEntityMetadata("txn")
	ctx := memorystore.WithEntityConnection(common.GetRootContext(), transactionMetadataProvider)
	defer memorystore.Close(ctx)
	return transactionMetadataProvider.GetStore().MultiDelete(ctx, transactionMetadataProvider, txns)
}

/*SetPreviousBlock - set the previous block */
func (mc *Chain) SetPreviousBlock(ctx context.Context, r *round.Round, b *block.Block, pb *block.Block) {
	if r == nil {
		mr := mc.GetRound(b.Round)
		if mr != nil {
			r = mr.Round
		} else {
			r = datastore.GetEntityMetadata("round").Instance().(*round.Round)
			r.Number = b.Round
			r.RandomSeed = b.RoundRandomSeed
		}
	}
	b.SetPreviousBlock(pb)
	b.RoundRandomSeed = r.RandomSeed
	bNode := node.GetNode(b.MinerID)
	b.RoundRank = r.GetMinerRank(bNode.SetIndex)
	b.ComputeChainWeight()
}
