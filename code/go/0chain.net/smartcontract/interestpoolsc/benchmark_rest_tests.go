package interestpoolsc

import (
	"context"
	"net/url"

	cstate "0chain.net/chaincore/chain/state"
	"0chain.net/chaincore/smartcontract"
	sci "0chain.net/chaincore/smartcontractinterface"
	"0chain.net/chaincore/transaction"
	bk "0chain.net/smartcontract/benchmark"
)

type RestBenchTest struct {
	name     string
	endpoint func(
		context.Context,
		url.Values,
		cstate.StateContextI,
	) (interface{}, error)
	params url.Values
}

func (rbt RestBenchTest) Name() string {
	return rbt.name
}

func (bt RestBenchTest) Transaction() *transaction.Transaction {
	return &transaction.Transaction{}
}

func (rbt RestBenchTest) Run(balances cstate.StateContextI) {
	_, err := rbt.endpoint(context.TODO(), rbt.params, balances)
	if err != nil {
		panic(err)
	}
}

func BenchmarkRestTests(
	data bk.BenchData, sigScheme bk.SignatureScheme,
) bk.TestSuit {
	var isc = InterestPoolSmartContract{
		SmartContract: sci.NewSC(ADDRESS),
	}
	isc.setSC(isc.SmartContract, &smartcontract.BCContext{})
	var tests = []RestBenchTest{
		{
			name:     "interest_pool_rest.getPoolsStats",
			endpoint: isc.getPoolsStats,
			params: func() url.Values {
				var values url.Values = make(map[string][]string)
				values.Set("client_id", data.Clients[0])
				return values
			}(),
		},
		{
			name:     "interest_pool_rest.getLockConfig",
			endpoint: isc.getLockConfig,
		},
	}
	var testsI []bk.BenchTestI
	for _, test := range tests {
		testsI = append(testsI, test)
	}
	return bk.TestSuit{
		Source:     bk.InterestPoolRest,
		Benchmarks: testsI,
	}
}
