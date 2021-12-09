package producer

import (
	"fmt"
	"sync"

	"github.com/plusar-it/hal/pkg/command"
)

// Producer executes the triggers (commands) and pushes results in the BlockNumbers channel
type Producer struct {
	Multicaller   Multicaller
	BlockNumbers  chan uint64
	noOfProducers int
	batchSize     int
	wg            sync.WaitGroup
}

type Multicaller interface {
	Call(triggers []command.Trigger) (uint64, error)
}

// NewProducer is a constructor for the Producer
func NewProducer(multicaller Multicaller, noOfProducers, batchSize int) *Producer {
	return &Producer{
		Multicaller:   multicaller,
		BlockNumbers:  make(chan uint64),
		noOfProducers: noOfProducers,
		batchSize:     batchSize,
		wg:            sync.WaitGroup{},
	}
}

// Execute starts the producer
func (p *Producer) Execute(triggers []command.Trigger) error {
	// map used to split the triggers processed by each routine
	triggerMap := p.getTriggerMap(triggers)

	for i := 0; i < p.noOfProducers; i++ {
		p.wg.Add(1)
		go p.processTriggers(triggerMap[i])
	}

	p.wg.Wait()
	close(p.BlockNumbers)

	return nil
}

// processTriggers prepares the batch request and calls the multicaller;
// if the resulted block number is even , push it to the sync channel
func (p *Producer) processTriggers(triggers []command.Trigger) {
	defer p.wg.Done()

	for i := 0; i < len(triggers); i += p.batchSize {
		//process the slice in batches of "batchSize"
		j := i + p.batchSize
		if j > len(triggers) {
			j = len(triggers)
		}

		blockNumber, err := p.Multicaller.Call(triggers[i:j])
		if err != nil {
			fmt.Printf("Failed multicaller call, but we continue")
		}

		//if blockNumer is an even number, publish to channel
		if blockNumber%2 == 0 {
			p.BlockNumbers <- blockNumber
		}
	}
}

// getTriggerMap creates a map of trigger slices, each slice will be processed by a dedicated go routine
func (p *Producer) getTriggerMap(triggers []command.Trigger) map[int][]command.Trigger {
	// map splits in slices the triggers processed by each routine
	triggerMap := make(map[int][]command.Trigger)

	for i, trigger := range triggers {
		index := i % p.noOfProducers
		triggerMap[index] = append(triggerMap[index], trigger)
	}

	return triggerMap
}
