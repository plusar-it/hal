package consumer

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

type mockStore struct {
	storedItems []uint64
}

func newmockStore() *mockStore {
	return &mockStore{storedItems: make([]uint64, 0)}
}

func (s *mockStore) Save(value uint64) error {
	s.storedItems = append(s.storedItems, value)
	return nil
}

func Test_consumer_execute(t *testing.T) {
	tests := []struct {
		inputValues []uint64
		wantValues  []uint64
	}{
		{
			inputValues: []uint64{4, 6, 8},
			wantValues:  []uint64{4, 6, 8},
		},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			mockStore := newmockStore()
			blockNumbers := make(chan uint64)
			consumer := NewConsumer(mockStore, blockNumbers)

			consumer.Execute()

			for _, v := range tt.inputValues {
				blockNumbers <- v
			}

			close(blockNumbers)

			require.Equal(t, tt.wantValues, mockStore.storedItems)
		})
	}

}
