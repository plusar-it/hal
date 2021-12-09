package main

import (
	"fmt"

	"github.com/plusar-it/hal/pkg/command"
	"github.com/plusar-it/hal/pkg/config"
	"github.com/plusar-it/hal/pkg/consumer"
	"github.com/plusar-it/hal/pkg/db"
	"github.com/plusar-it/hal/pkg/eth"
	"github.com/plusar-it/hal/pkg/producer"
)

func main() {

	// load the initial configuration settings
	configuration := config.LoadConfig()

	// ethereum client used by the multicaller
	ethClient := eth.NewEthClient(configuration.FromAddress, configuration.NetworkAddress)

	// multicaller supports calling triggers in batches
	multicaller, err := eth.NewMulticaller(ethClient)
	if err != nil {
		fmt.Println(fmt.Errorf("cannot initialise multicaller: %w", err))
		return
	}

	// triggers provider provides the triggers to be produced
	triggerProvider := command.NewTriggerProvider(configuration.TriggersSourceFilePath)
	triggers, err := triggerProvider.GetTriggers()
	if err != nil {
		fmt.Println(fmt.Errorf("Error getting triggers from provider: %w", err))
		return
	}

	producer := producer.NewProducer(multicaller, configuration.NoOfProducers, configuration.BatchSize)
	consumer := consumer.NewConsumer(&db.Store{}, producer.BlockNumbers)

	// start the consumer that will listen on the BlockNumbers channel, for messages from the producer
	// the consumer will save all received messages in the db
	consumer.Execute()

	// the producer will start the producing gouritines and waits (blocks) until they have finished
	// after the work is done, the channel is closed by the producer
	producer.Execute(triggers)

	// wait for the consumer to finish
	// the consumer returns as soon as the channel is closed by the producer
	<-consumer.Done
}
