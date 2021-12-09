package eth

import (
	"fmt"

	"github.com/HAL-xyz/ethrpc"
)

// EthClient supports RPC calls to the ethereum network
type EthClient struct {
	ethRpcClient *ethrpc.EthRPC
	fromAddress  string
}

// NewEthClient is a constructor for EthClient
func NewEthClient(from, networkAddress string) *EthClient {
	return &EthClient{
		ethRpcClient: ethrpc.New(networkAddress),
		fromAddress:  from,
	}
}

// MakeEthRpcCall makes the RPC call to ethereum
func (c *EthClient) MakeEthRpcCall(contractAddress, data string, blockNumber int) (string, error) {
	params := ethrpc.T{
		To:   contractAddress,
		From: c.fromAddress,
		Data: data,
	}
	hexBlockNo := fmt.Sprintf("0x%x", blockNumber)
	return c.ethRpcClient.EthCall(params, hexBlockNo)
}
