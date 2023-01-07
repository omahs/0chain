package event

import (
	"fmt"
	"time"

	"github.com/0chain/common/core/logging"

	"0chain.net/smartcontract/stakepool/spenum"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"0chain.net/smartcontract/dbs"
)

type providerRewardsDelegates struct {
	rewards        []ProviderRewards
	delegatePools  []DelegatePool
	desc           [][]string
	rewardProvider []RewardProvider
	rewardDelegate []RewardDelegate
}

func aggregateProviderRewards(
	spus []dbs.StakePoolReward, round int64,
) (*providerRewardsDelegates, error) {
	var (
		rewards        = make([]ProviderRewards, 0, len(spus))
		delegatePools  = make([]DelegatePool, 0, len(spus))
		descs          = make([][]string, 0, len(spus))
		rewardProvider = make([]RewardProvider, 0, len(spus))
		rewardDelegate = make([]RewardDelegate, 0, len(spus))
	)
	rewardProvider = rewardProvider
	rewardDelegate = rewardDelegate
	for i, sp := range spus {
		if sp.Reward != 0 {
			rewards = append(rewards, ProviderRewards{
				ProviderID:       sp.ProviderId,
				Rewards:          sp.Reward,
				TotalRewards:     sp.Reward,
				RoundLastUpdated: round,
			})
		}

		// merge delegate rewards and penalties
		for k, v := range spus[i].DelegateRewards {
			delegatePools = append(delegatePools, DelegatePool{
				ProviderID:       sp.ProviderId,
				ProviderType:     sp.ProviderType,
				PoolID:           k,
				Reward:           v,
				TotalReward:      v,
				TotalPenalty:     spus[i].DelegatePenalties[k],
				RoundLastUpdated: round,
			})
		}

		// append remaining penalties if any
		for k, v := range spus[i].DelegatePenalties {
			if _, ok := sp.DelegateRewards[k]; !ok {
				delegatePools = append(delegatePools, DelegatePool{
					ProviderID:       sp.ProviderId,
					ProviderType:     sp.ProviderType,
					PoolID:           k,
					TotalPenalty:     v,
					RoundLastUpdated: round,
				})
			}
		}
	}

	return &providerRewardsDelegates{
		rewards:       rewards,
		delegatePools: delegatePools,
		desc:          descs,
	}, nil
}

func mergeStakePoolRewardsEvents() *eventsMergerImpl[dbs.StakePoolReward] {
	return newEventsMerger[dbs.StakePoolReward](TagStakePoolReward, withProviderRewardsPenaltiesAdded())
}

// withProviderRewardsPenaltiesAdded is an event merger middleware that merge two
// StakePoolRewards
func withProviderRewardsPenaltiesAdded() eventMergeMiddleware {
	return withEventMerge(func(a, b *dbs.StakePoolReward) (*dbs.StakePoolReward, error) {
		a.Reward += b.Reward

		// merge delegate pool rewards
		for k, v := range b.DelegateRewards {
			_, ok := a.DelegateRewards[k]
			if !ok {
				a.DelegateRewards[k] = v
				continue
			}

			a.DelegateRewards[k] += v
		}

		// merge delegate pool penalties
		for k, v := range b.DelegatePenalties {
			_, ok := a.DelegatePenalties[k]
			if !ok {
				a.DelegatePenalties[k] = v
				continue
			}

			a.DelegatePenalties[k] += v
		}

		return a, nil
	})
}

func (edb *EventDb) rewardUpdate(spus []dbs.StakePoolReward, round int64) error {
	if len(spus) == 0 {
		return nil
	}

	ts := time.Now()
	rewards, err := aggregateProviderRewards(spus, round)
	if err != nil {
		return err
	}
	logging.Logger.Info("rewardUpdate",
		zap.Bool("debug", edb.Debug()),
		zap.Int64("round", round),
		zap.Any("rewards", rewards))
	defer func() {
		du := time.Since(ts)
		n := len(rewards.rewards) + len(rewards.delegatePools)
		if du > 50*time.Millisecond {
			logging.Logger.Debug("event db - update reward slow",
				zap.Int64("round", round),
				zap.Any("duration", du),
				zap.Int("total update items", n),
				zap.Int("rewards num", len(rewards.rewards)),
				zap.Int("delegate pools num", len(rewards.delegatePools)),
				zap.Any("desc", rewards.desc))
		}
	}()

	if len(rewards.rewards) > 0 {
		if err := edb.rewardProviders(rewards.rewards); err != nil {
			return fmt.Errorf("could not rewards providers: %v", err)
		}
	}

	rpdu := time.Since(ts)
	if rpdu.Milliseconds() > 50 {
		logging.Logger.Debug("event db - reward providers slow",
			zap.Any("duration", rpdu),
			zap.Int64("round", round))
	}

	if len(rewards.delegatePools) > 0 {
		if err := rewardProviderDelegates(edb, rewards.delegatePools); err != nil {
			return fmt.Errorf("could not rewards delegate pool: %v", err)
		}
	}

	if edb.Debug() {
		if err := edb.insertProviderReward(spus, round); err != nil {
			return err
		}
		if err := edb.insertDelegateReward(spus, round); err != nil {
			return err
		}
	}

	return nil
}

func (edb *EventDb) rewardProviders(rewards []ProviderRewards) error {
	logging.Logger.Info("piers rewardProviders", zap.Any("rewards", rewards))
	vs := map[string]interface{}{
		"rewards":            gorm.Expr("provider_rewards.rewards + excluded.rewards"),
		"total_rewards":      gorm.Expr("provider_rewards.total_rewards + excluded.total_rewards"),
		"round_last_updated": gorm.Expr("excluded.round_last_updated"),
	}

	return edb.Store.Get().Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "provider_id"}},
		DoUpdates: clause.Assignments(vs),
	}).Create(rewards).Error
}

func rewardProviderDelegates(edb *EventDb, rewards []DelegatePool) error {
	vs := map[string]interface{}{
		"reward":        gorm.Expr("delegate_pools.reward + excluded.reward"),
		"total_reward":  gorm.Expr("delegate_pools.total_reward + excluded.total_reward"),
		"total_penalty": gorm.Expr("delegate_pools.total_penalty + excluded.total_penalty"),
	}

	return edb.Store.Get().Clauses(clause.OnConflict{
		Where: clause.Where{
			Exprs: []clause.Expression{gorm.Expr("delegate_pools.status != ?", spenum.Deleted)},
		},
		Columns: []clause.Column{
			{Name: "provider_type"},
			{Name: "provider_id"},
			{Name: "pool_id"},
		},
		DoUpdates: clause.Assignments(vs),
	}).Create(&rewards).Error
}

func (edb *EventDb) rewardProvider(spu dbs.StakePoolReward) error { //nolint: unused
	if spu.Reward == 0 {
		return nil
	}

	var provider interface{}
	switch spenum.Provider(spu.ProviderType) {
	case spenum.Blobber:
		provider = &Blobber{Provider: Provider{ID: spu.ProviderId}}
	case spenum.Validator:
		provider = &Validator{Provider: Provider{ID: spu.ProviderId}}
	case spenum.Miner:
		provider = &Miner{Provider: Provider{ID: spu.ProviderId}}
	case spenum.Sharder:
		provider = &Sharder{Provider: Provider{ID: spu.ProviderId}}
	default:
		return fmt.Errorf("not implented provider type %v", spu.ProviderType)
	}

	vs := map[string]interface{}{
		"rewards":      gorm.Expr("rewards + ?", spu.Reward),
		"total_reward": gorm.Expr("total_reward + ?", spu.Reward),
	}

	return edb.Store.Get().Model(provider).Where(provider).Updates(vs).Error
}
