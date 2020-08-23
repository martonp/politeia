package paywall

import (
	"fmt"
	"testing"
	"time"

	pstypes "github.com/decred/dcrdata/pubsub/types/v3"
	client "github.com/decred/dcrdata/pubsub/v4/psclient"

	"github.com/decred/politeia/politeiawww/wsdcrdata"

	"github.com/decred/politeia/util/txfetcher"
)

type AddressPaywall struct {
	Address    string
	AmountPaid uint64
	Fulfilled  bool
}

func createUpdatePaywallsCallback(paywalls *[]*AddressPaywall) Callback {

	return func(entry *Entry, txs []txfetcher.TxDetails, fulfilled bool) error {
		var foundPaywall bool

		for _, paywall := range *paywalls {

			if paywall.Address == entry.address {
				foundPaywall = true
			} else {
				continue
			}

			var amount uint64
			for i := 0; i < len(txs); i++ {
				amount += txs[i].Amount
			}

			paywall.AmountPaid = amount
			paywall.Fulfilled = fulfilled

			if fulfilled && amount < entry.amount {
				return fmt.Errorf("fulfilled but not enough paid")
			}

			break
		}

		if !foundPaywall {
			return fmt.Errorf("Callback called with address without entry")
		}

		return nil
	}
}

func TestTransactionsFulfillPaywall(t *testing.T) {
	testEntries := []Entry{
		{
			address:     "1",
			amount:      1000,
			txNotBefore: 100,
		},
		{
			address:     "2",
			amount:      1000,
			txNotBefore: 100,
		},
		{
			address:     "3",
			amount:      1000,
			txNotBefore: 100,
		},
		{
			address:     "4",
			amount:      1000,
			txNotBefore: 100,
		},
	}

	testTXs := []txfetcher.TxDetails{
		txfetcher.TxDetails{
			Address:   "1",
			TxID:      "1",
			Amount:    1000,
			Timestamp: 80,
		},
		txfetcher.TxDetails{
			Address:   "2",
			TxID:      "2",
			Amount:    1000,
			Timestamp: 120,
		},
		txfetcher.TxDetails{
			Address:   "3",
			TxID:      "3",
			Amount:    800,
			Timestamp: 130,
		},
		txfetcher.TxDetails{
			Address:   "4",
			TxID:      "4",
			Amount:    500,
			Timestamp: 120,
		},
		txfetcher.TxDetails{
			Address:   "4",
			TxID:      "5",
			Amount:    700,
			Timestamp: 120,
		},
	}

	expectedResults := []AddressPaywall{
		AddressPaywall{
			Address:    "1",
			AmountPaid: 0,
			Fulfilled:  false,
		},
		AddressPaywall{
			Address:    "2",
			AmountPaid: 1000,
			Fulfilled:  true,
		},
		AddressPaywall{
			Address:    "3",
			AmountPaid: 800,
			Fulfilled:  false,
		},
		AddressPaywall{
			Address:    "4",
			AmountPaid: 1200,
			Fulfilled:  true,
		},
	}

	paywalls := make([]*AddressPaywall, 0)

	txFetcher := txfetcher.NewTestTxFetcher()
	wsDcrdata := wsdcrdata.NewTestWSDcrdata()

	paywallManager := NewDcrdataManager(wsDcrdata, txFetcher)

	for _, entry := range testEntries {
		paywallManager.RegisterPaywall(entry)
		paywalls = append(paywalls, &AddressPaywall{
			Address: entry.address,
		})
	}

	callback := createUpdatePaywallsCallback(&paywalls)
	paywallManager.SetCallback(callback)

	for _, tx := range testTXs {
		txFetcher.InsertTx(tx)

		wsDcrdata.SendMessage(
			&client.ClientMessage{
				EventId: "",
				Message: &pstypes.AddressMessage{
					Address: tx.Address,
					TxHash:  tx.TxID,
				},
			})
	}

	time.Sleep(1 * time.Second)

	if len(expectedResults) != len(paywalls) {
		t.Fatal("results and expected results have different lengths")
	}

	for i := 0; i < len(expectedResults); i++ {
		if expectedResults[i].Address != paywalls[i].Address {
			t.Fatal("results and expected results are in the wrong order")
		}
		if expectedResults[i].AmountPaid != paywalls[i].AmountPaid {
			t.Fatal("results and expected results amount paid does not match")
		}
		if expectedResults[i].Fulfilled != paywalls[i].Fulfilled {
			t.Fatal("results and expected results fulfilled does not match")
		}
	}
}
