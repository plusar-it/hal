package producer

import (
	"sort"
	"strconv"
	"testing"

	"github.com/plusar-it/hal/pkg/command"

	"github.com/stretchr/testify/require"
)

type mockMulticaler struct {
	storedItems []command.Trigger
	blockNumber uint64
}

func newMockMulticaler(blockNumber uint64) *mockMulticaler {
	return &mockMulticaler{storedItems: make([]command.Trigger, 0), blockNumber: blockNumber}
}

func (s *mockMulticaler) Call(triggers []command.Trigger) (uint64, error) {
	s.storedItems = append(s.storedItems, triggers...)

	return s.blockNumber, nil
}

func (s *mockMulticaler) sortStoredItems() {
	sort.Slice(s.storedItems, func(i, j int) bool {
		return s.storedItems[i].TriggerName < s.storedItems[j].TriggerName
	})
}

type mockConsumer struct {
	storedItems []uint64
}

func newMockConsumer() *mockConsumer {
	return &mockConsumer{storedItems: make([]uint64, 0)}
}

func (c *mockConsumer) Consume(blockNumbers chan uint64) {
	for {
		blockNumber, more := <-blockNumbers
		// exit when the channel is closed
		if !more {
			return
		}
		c.storedItems = append(c.storedItems, blockNumber)
	}
}

func Test_producer_execute(t *testing.T) {
	tests := []struct {
		inputValues          []command.Trigger
		wantValues           []command.Trigger
		inputBlockNumber     uint64
		consumedBlockNumbers []uint64
	}{
		{
			// each routine will handle 1 trigger; the batch size is 2, we ensure scenario when noOf triggers lower than batch size
			inputValues: []command.Trigger{
				{TriggerName: "BAT", UserAddress: "0x0000000000000", ContractAddress: "0x0d8775f6484", Method: "balanceOf(address)(uint256)"},
				{TriggerName: "BC", UserAddress: "0x0000000000000", ContractAddress: "0x21ab6c9fac80", Method: "balanceOf(address)(uint256)"},
				{TriggerName: "ETHLAND", UserAddress: "0x0000000000000", ContractAddress: "0x80fB784B7eD66", Method: "balanceOf(address)(uint256)"},
			},
			wantValues: []command.Trigger{
				{TriggerName: "BAT", UserAddress: "0x0000000000000", ContractAddress: "0x0d8775f6484", Method: "balanceOf(address)(uint256)"},
				{TriggerName: "BC", UserAddress: "0x0000000000000", ContractAddress: "0x21ab6c9fac80", Method: "balanceOf(address)(uint256)"},
				{TriggerName: "ETHLAND", UserAddress: "0x0000000000000", ContractAddress: "0x80fB784B7eD66", Method: "balanceOf(address)(uint256)"},
			},
			inputBlockNumber:     2,                 // even block number
			consumedBlockNumbers: []uint64{2, 2, 2}, // 3 triggers, each returns an even number
		},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			mockMulticaler := newMockMulticaler(tt.inputBlockNumber)

			// the producer will use 3 routines, batch calls of 2 requests
			producer := NewProducer(mockMulticaler, 3, 2)

			consumer := newMockConsumer()
			go consumer.Consume(producer.BlockNumbers)

			producer.Execute(tt.inputValues)

			// we sort before we compare, the order is not maintained by the producer, the handling being done on different go routines
			mockMulticaler.sortStoredItems()

			require.Equal(t, tt.wantValues, mockMulticaler.storedItems)
			require.Equal(t, tt.consumedBlockNumbers, consumer.storedItems)
		})
	}
}
