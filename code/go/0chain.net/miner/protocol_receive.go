package miner

import (
	"context"
	"time"

	"0chain.net/chaincore/block"
	"0chain.net/core/common"
	"0chain.net/core/logging"
	"go.uber.org/zap"
)

// HandleVRFShare - handles the vrf share.
func (mc *Chain) HandleVRFShare(ctx context.Context, msg *BlockMessage) {

	var mr = mc.getOrStartRoundNotAhead(ctx, msg.VRFShare.Round)
	if mr == nil {
		return
	}

	// add the VRFS
	mc.AddVRFShare(ctx, mr, msg.VRFShare)
}

// HandleVerifyBlockMessage - handles the verify block message.
func (mc *Chain) HandleVerifyBlockMessage(ctx context.Context,
	msg *BlockMessage) {

	var b = msg.Block

	var mr, pr = mc.GetMinerRound(b.Round), mc.GetMinerRound(b.Round - 1)
	if pr == nil {
		logging.Logger.Error("handle verify block -- no previous round (ignore)",
			zap.Int64("round", b.Round), zap.Int64("prev_round", b.Round-1))
		return
	}

	// return if the block already in local chain and its previous block is notarized
	_, err := mc.GetBlock(ctx, b.Hash)
	if err == nil { // block already exist in local chain
		// check if previous block exist and is notarized
		pb, err := mc.GetBlock(ctx, b.PrevHash)
		if err == nil && pb != nil && pb.IsBlockNotarized() {
			logging.Logger.Debug("handle verify block - block already exist, ignore",
				zap.Int64("round", b.Round))
			return
		}
	}

	logging.Logger.Debug("verify block handler",
		zap.Int64("round", b.Round),
		zap.String("block", b.Hash))

	if err := b.Validate(ctx); err != nil {
		logging.Logger.Debug("verify block handler -- can't validate",
			zap.Int64("round", b.Round), zap.Error(err))
		return
	}

	if b.Round < mc.GetCurrentRound()-1 {
		logging.Logger.Debug("verify block (round mismatch)",
			zap.Int64("current_round", mc.GetCurrentRound()),
			zap.Int64("block_round", b.Round))
		return
	}

	// get previous block notarization tickets, and update local prev block if exist
	if b.Round > 1 {
		// TODO: run in gorountine for debug and test purpose
		// do not run this in goroutine
		//
		// put into a goroutine so that tickets verification would not affect the
		// new round RRS generation
		go func() {
			// TODO: check if the block's prev notarized block reached the notarization threshold
			if err := mc.updatePreviousBlockNotarization(ctx, b, pr); err != nil {
				return
			}
		}()
	}

	if mr == nil {
		logging.Logger.Error("handle verify block -- got block proposal before starting round",
			zap.Int64("round", b.Round), zap.String("block", b.Hash),
			zap.String("miner", b.MinerID))

		mr = mc.getOrStartRoundNotAhead(ctx, b.Round)
		if mr == nil {
			logging.Logger.Error("handle verify block -- can't start new round",
				zap.Int64("round", b.Round))
			return
		}

		mc.startRound(ctx, mr, b.GetRoundRandomSeed())

		mc.AddToRoundVerification(ctx, mr, b)
		return
	}

	if !mr.IsVRFComplete() {
		logging.Logger.Info("handle verify block - got block proposal before VRF is complete",
			zap.Int64("round", b.Round), zap.String("block", b.Hash),
			zap.String("miner", b.MinerID))

		if mr.GetTimeoutCount() < b.RoundTimeoutCount {
			logging.Logger.Info("Insync ignoring handle verify block - got block proposal before VRF is complete",
				zap.Int64("round", b.Round), zap.String("block", b.Hash),
				zap.String("miner", b.MinerID),
				zap.Int("round_toc", mr.GetTimeoutCount()),
				zap.Int("round_toc", b.RoundTimeoutCount))
			return
		}

		if b.GetRoundRandomSeed() != mr.GetRandomSeed() {
			logging.Logger.Info("handle verify block - got block with different RRS",
				zap.Int64("round", b.Round),
				zap.Int64("block RRS", b.GetRoundRandomSeed()),
				zap.Int64("round RRS", mr.GetRandomSeed()))
			mc.startRound(ctx, mr, b.GetRoundRandomSeed())
		}
	}

	vts := mr.GetVerificationTickets(b.Hash)
	if len(vts) == 0 {
		mc.AddToRoundVerification(ctx, mr, b)
		return
	}

	// TODO: mc.MergeVerificationTickets does not verify block's own tickets, might be a problem!
	mc.MergeVerificationTickets(b, vts)
	if !b.IsBlockNotarized() {
		mc.AddToRoundVerification(ctx, mr, b)
		return
	}

	if mr.GetRandomSeed() == b.GetRoundRandomSeed() {
		b = mc.AddRoundBlock(mr, b)
		mc.checkBlockNotarization(ctx, mr, b)
		return
	}

	/* Since this is a notarized block, we are accepting it. */
	b1, r1, err := mc.AddNotarizedBlockToRound(mr, b)
	if err != nil {
		logging.Logger.Error("handle verify block failed",
			zap.Int64("round", b.Round),
			zap.String("block", b.Hash),
			zap.String("miner", b.MinerID),
			zap.Error(err))
		return
	}

	b = b1
	mr = r1.(*Round)
	logging.Logger.Info("Added a notarizedBlockToRound - got notarized block with different RRS",
		zap.Int64("round", b.Round),
		zap.String("block", b.Hash),
		zap.String("miner", b.MinerID),
		zap.Int("round_toc", mr.GetTimeoutCount()),
		zap.Int("round_toc", b.RoundTimeoutCount))

	mc.checkBlockNotarization(ctx, mr, b)
}

func (mc *Chain) verifyTicketsWithRetry(ctx context.Context,
	r int64, block string, bvts []*block.VerificationTicket, retryN int) error {
	for i := 0; i < retryN; i++ {
		err := func() error {
			logging.Logger.Debug("verification ticket",
				zap.Int64("round", r),
				zap.String("block", block),
				zap.Int("retry", i))
			cctx, cancel := context.WithTimeout(ctx, time.Second)
			defer cancel()
			return mc.VerifyTickets(cctx, block, bvts, r)
		}()

		switch err {
		case nil:
			return nil
		case context.DeadlineExceeded:
			if mc.GetCurrentRound() > r {
				return common.NewErrorf("verify_tickets_timeout", "chain moved on, round: %d", r)
			}
		default:
			logging.Logger.Error("verification ticket failed",
				zap.Int64("round", r),
				zap.Error(err))
			return err
		}
	}

	return common.NewErrorf("verify_tickets_timeout", "ticket timeout with retry, round: %d", r)
}

// HandleVerificationTicketMessage - handles the verification ticket message.
func (mc *Chain) HandleVerificationTicketMessage(ctx context.Context,
	msg *BlockMessage) {

	var (
		bvt = msg.BlockVerificationTicket
		rn  = bvt.Round
	)

	logging.Logger.Debug("handle vt. msg - verification ticket",
		zap.Int64("round", bvt.Round),
		zap.String("block", bvt.BlockID))

	var mr = mc.getOrStartRoundNotAhead(ctx, rn)
	if mr == nil {
		logging.Logger.Error("handle vt. msg -- ahead of sharders or no pr",
			zap.Int64("round", rn))
		return
	}

	if mc.GetMinerRound(rn-1) == nil {
		logging.Logger.Error("handle vt. msg -- no previous round (ignore)",
			zap.Int64("round", rn), zap.Int64("pr", rn-1))
		return
	}

	if mr.isVerificationComplete() {
		logging.Logger.Error("handle vt. msg -- round verification completed", zap.Int64("round", rn))
		return
	}

	// check if the ticket has already verified
	if mr.IsTicketCollected(&bvt.VerificationTicket) {
		logging.Logger.Error("handle vt. msg -- ticket already collected",
			zap.Int64("round", rn), zap.String("block", bvt.BlockID))
		return
	}

	err := mc.verifyTicketsWithRetry(ctx, rn, bvt.BlockID, []*block.VerificationTicket{&bvt.VerificationTicket}, 3)
	if err != nil {
		logging.Logger.Error("handle vt. msg - verification ticket failed", zap.Error(err))
		return
	}

	b, err := mc.GetBlock(ctx, bvt.BlockID)
	if err != nil {
		logging.Logger.Debug("handle vt. msg - block does not exist, collect tickets though",
			zap.Int64("round", bvt.Round),
			zap.String("block", bvt.BlockID))
		mr.AddVerificationTickets([]*block.BlockVerificationTicket{bvt})
		return
	}

	var lfb = mc.GetLatestFinalizedBlock()
	if b.Round < lfb.Round {
		logging.Logger.Debug("verification message (round mismatch)",
			zap.Int64("round", b.Round), zap.String("block", b.Hash),
			zap.Int64("lfb", lfb.Round))
		return
	}

	mc.ProcessVerifiedTicket(ctx, mr, b, &bvt.VerificationTicket)
}

// HandleNotarizationMessage - handles the block notarization message.
func (mc *Chain) HandleNotarizationMessage(ctx context.Context, msg *BlockMessage) {

	var (
		lfb = mc.GetLatestFinalizedBlock()
		not = msg.Notarization
	)

	if not.Round < lfb.Round {
		logging.Logger.Debug("handle notarization message",
			zap.Int64("round", not.Round),
			zap.Int64("finalized_round", lfb.Round),
			zap.String("block", not.BlockID))
		return
	}

	var r = mc.GetMinerRound(not.Round)
	if r == nil {
		if msg.ShouldRetry() {
			logging.Logger.Error("handle notarization message (round not started yet) retrying",
				zap.Int64("round", not.Round),
				zap.String("block", not.BlockID),
				zap.Int8("retry_count", msg.RetryCount))

			msg.Retry(mc.blockMessageChannel)
		} else {
			logging.Logger.Error("handle notarization message (round not started yet)",
				zap.Int64("round", not.Round),
				zap.String("block", not.BlockID),
				zap.Int8("retry_count", msg.RetryCount))
		}
		return
	}

	if mc.GetMinerRound(not.Round-1) == nil {
		logging.Logger.Error("handle notarization message -- no previous round",
			zap.Int64("round", not.Round),
			zap.Int64("prev_round", not.Round-1))
		return
	}

	msg.Round = r

	var b, err = mc.GetBlock(ctx, not.BlockID)
	if err != nil {
		// TODO: add this back
		logging.Logger.Debug("handle notarization message -- not exist, try to fetch notarized block",
			zap.Int64("round", not.Round),
			zap.String("block", not.BlockID))
		go mc.GetNotarizedBlock(ctx, not.BlockID, not.Round)
		return
	}

	if !b.IsBlockNotarized() {
		var vts = b.UnknownTickets(not.VerificationTickets)
		if len(vts) == 0 {
			return
		}

		logging.Logger.Debug("handle notarization message -- merge notarization block",
			zap.Int64("round", b.Round),
			zap.String("block", b.Hash))
		go mc.MergeNotarization(ctx, r, b, vts)
	}
}

// HandleNotarizedBlockMessage - handles a notarized block for a previous round.
func (mc *Chain) HandleNotarizedBlockMessage(ctx context.Context,
	msg *BlockMessage) {

	nb := msg.Block

	//var mc = GetMinerChain()
	if nb.Round < mc.GetCurrentRound()-1 {
		logging.Logger.Debug("notarized block handler (round older than the current round)",
			zap.String("block", nb.Hash), zap.Any("round", nb.Round))
		return
	}

	var lfb = mc.GetLatestFinalizedBlock()
	if nb.Round <= lfb.Round {
		return // doesn't need the not. block
	}

	if mc.GetMinerRound(nb.Round-1) == nil {
		logging.Logger.Error("not. block handler -- no previous round (ignore)",
			zap.Int64("round", nb.Round), zap.Int64("prev_round", nb.Round-1))
		return // no previous round
	}

	var mr = mc.getOrStartRoundNotAhead(ctx, nb.Round)
	if mr == nil {
		logging.Logger.Debug("notarized block handler -- is ahead or no pr",
			zap.String("block", nb.Hash), zap.Any("round", nb.Round),
			zap.Bool("has_pr", mc.GetMinerRound(nb.Round-1) != nil))
		return // can't handle yet
	}

	if mr.GetRandomSeed() == 0 {
		mc.SetRandomSeed(mr, nb.GetRoundRandomSeed())
	}

	if mr.IsFinalizing() || mr.IsFinalized() {
		return // doesn't need a not. block
	}

	if mr.IsVerificationComplete() {
		return // verification for the round complete
	}

	for _, blk := range mr.GetNotarizedBlocks() {
		if blk.Hash == nb.Hash {
			return // already have
		}
	}

	if err := mc.VerifyNotarization(ctx, nb, nb.GetVerificationTickets(), mr.GetRoundNumber()); err != nil {
		logging.Logger.Error("not. block handler -- verify notarization failed",
			zap.Int64("round", nb.Round),
			zap.String("block", nb.Hash),
			zap.Error(err))
		return
	}

	if !mr.IsVRFComplete() {
		mc.startRound(ctx, mr, nb.GetRoundRandomSeed())
	}

	var b = mc.AddRoundBlock(mr, nb)
	if !mc.AddNotarizedBlock(ctx, mr, b) {
		return
	}

	if mc.isAheadOfSharders(ctx, nb.Round+1) {
		logging.Logger.Error("handle not. block -- ahead of sharders",
			zap.Int64("round", nb.Round+1))
		return
	}

	mc.StartNextRound(ctx, mr) // start next or skip
}
