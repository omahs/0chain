package storagesc

import (
	"fmt"

	cstate "0chain.net/chaincore/chain/state"
	"0chain.net/chaincore/transaction"
	"0chain.net/core/common"
	"0chain.net/smartcontract/provider"
)

func (ssc *StorageSmartContract) blobberHealthCheck(
	t *transaction.Transaction,
	_ []byte,
	balances cstate.StateContextI,
) (string, error) {
	var blobber *StorageNode
	var err error
	if blobber, err = ssc.getBlobber(t.ClientID, balances); err != nil {
		return "", common.NewError("blobber_health_check_failed",
			"can't get the blobber "+t.ClientID+": "+err.Error())
	}

	conf, err := getConfig(balances)
	if err != nil {
		return "", common.NewError("blobber_health_check_failed",
			"can't get configs"+err.Error())
	}

	blobber.HealthCheck(t.CreationDate, conf.HealthCheckPeriod, balances)

	if _, err = balances.InsertTrieNode(blobber.GetKey(ssc.ID), blobber); err != nil {
		return "", common.NewError("blobber_health_check_failed",
			"can't save blobber: "+err.Error())
	}

	return string(blobber.Encode()), nil
}

func (ssc *StorageSmartContract) validatorHealthCheck(
	t *transaction.Transaction,
	_ []byte,
	balances cstate.StateContextI,
) (string, error) {
	var validator = newValidatorNode(t.ClientID)
	if err := balances.GetTrieNode(validator.GetKey(ssc.ID), validator); err != nil {
		return "", common.NewError("blobber_health_check_failed",
			"can't get the blobber "+t.ClientID+": "+err.Error())
	}

	if err := healthCheck(validator, t.CreationDate, balances); err != nil {
		return "", common.NewError("validator_health_check_failed", err.Error())
	}

	if _, err := balances.InsertTrieNode(validator.GetKey(ssc.ID), validator); err != nil {
		return "", common.NewError("blobber_health_check_failed",
			"can't save blobber: "+err.Error())
	}

	return "", nil
}

func healthCheck(provider provider.Provider, now common.Timestamp, balances cstate.StateContextI) error {
	conf, err := getConfig(balances)
	if err != nil {
		return fmt.Errorf("can't get configs: %v", err)
	}

	provider.HealthCheck(now, conf.HealthCheckPeriod, balances)
}
