package logger

import (
	_ "embed"
)

func WithMethod(method string) Field {
	return WithString("method", method)
}

func WithNetworkName(network string) Field {
	return WithString("network_name", network)
}

func WithETHNetwork() Field {
	return WithNetworkName("eth")
}

func WithTransaction(txHash string) Field {
	return WithString("_transaction_hash", txHash)
}

func WithBlockNumber(blockNumber uint64) Field {
	return WithUnt64("block_number", blockNumber)
}

func WithService(serviceName string) Field {
	return WithString("service", serviceName)
}
