package command

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

// TriggerProvider provides the triggers to be processed from an underlying data source
type TriggerProvider struct {
	SourceFilePath string
}

// Trigger holds the details about a command (trigger)
type Trigger struct {
	TriggerName     string `json:"TriggerName"`
	UserAddress     string `json:"UserAddress"`
	ContractAddress string `json:"ContractAddress"`
	Method          string `json:"Method"`
}

// NewTriggerProvider is a constructor for the TriggerProvider
func NewTriggerProvider(sourceFilePath string) *TriggerProvider {
	return &TriggerProvider{SourceFilePath: sourceFilePath}
}

func (t *TriggerProvider) GetTriggers() ([]Trigger, error) {
	fileContent, err := ioutil.ReadFile(t.SourceFilePath)
	if err != nil {
		fmt.Println(fmt.Errorf("Error reading triggers file: %w", err))
		return nil, err
	}

	importedTriggers := make([]Trigger, 0)
	err = json.Unmarshal([]byte(fileContent), &importedTriggers)
	if err != nil {
		fmt.Println(fmt.Errorf("Error deserializing triggers file: %w", err))
		return nil, err
	}

	return removeDuplicates(importedTriggers), nil
}

// removes the duplicate triggers by contract address (assuming the method is always the same, balance)
func removeDuplicates(triggers []Trigger) []Trigger {
	result := make([]Trigger, 0)

	triggerMap := make(map[string]Trigger)
	for _, trigger := range triggers {
		if _, ok := triggerMap[trigger.ContractAddress]; !ok {
			triggerMap[trigger.ContractAddress] = trigger
			result = append(result, trigger)
		}
	}

	return result
}
