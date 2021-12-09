package eth

import (
	"fmt"

	"github.com/HAL-xyz/web3-multicall-go/multicall"
	"github.com/plusar-it/hal/pkg/command"
)

// Multicaller supports batch RPC calls
type Multicaller struct {
	ethMulticall multicall.Multicall
	ethClient    *EthClient
}

// NewMulticaller is a constructor of Multicaller
func NewMulticaller(client *EthClient) (*Multicaller, error) {
	m, err := multicall.New(client)
	if err != nil {
		return nil, err
	}
	return &Multicaller{ethMulticall: m, ethClient: client}, nil
}

// Call makes the batch multi call
func (m *Multicaller) Call(triggers []command.Trigger) (uint64, error) {
	vcs := multicall.ViewCalls{}

	for i, trigger := range triggers {
		key := fmt.Sprintf("key-%d", i)
		vcs = append(vcs, multicall.NewViewCall(
			key,
			trigger.ContractAddress,
			trigger.Method,
			[]interface{}{trigger.UserAddress},
		))
	}

	var lastBlockNo, err = m.ethClient.ethRpcClient.EthBlockNumber()
	if err != nil {
		return 0, err
	}

	result, err := m.ethMulticall.Call(vcs, fmt.Sprintf("%x", lastBlockNo))
	if err != nil {
		return 0, err
	}
	return result.BlockNumber, nil
}
