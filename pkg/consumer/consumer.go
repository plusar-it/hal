package consumer

import (
	"fmt"
)

// Consumer consumes the block number channel and stores the values in the db
type Consumer struct {
	Done         chan bool
	store        Store
	blockNumbers chan uint64
}

type Store interface {
	Save(value uint64) error
}

// NewConsumer is a constructor for the consumer
func NewConsumer(store Store, blockNumbers chan uint64) *Consumer {
	return &Consumer{blockNumbers: blockNumbers, store: store, Done: make(chan bool)}
}

// Execute starts the consumer process
func (c *Consumer) Execute() {
	go c.consume()
}

// consumes the block number channel and stores the values in the db
func (c *Consumer) consume() {
	for {
		blockNumber, more := <-c.blockNumbers
		// exit when the channel is closed
		if !more {
			c.Done <- true
		}
		err := c.store.Save(blockNumber)
		if err != nil {
			fmt.Printf("Error saving consumed block number \"%v\", but continue consuming \n", blockNumber)
		}
	}
}
