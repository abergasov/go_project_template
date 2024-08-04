package logger

import (
	_ "embed"
	"fmt"
)

func WithMethod(method string) StringWith {
	return StringWith{Key: "_method", Val: method}
}

func WithNetworkName(network string) StringWith {
	return StringWith{Key: "_network_name", Val: network}
}

func WithETHNetwork() StringWith {
	return StringWith{Key: "_network", Val: "eth"}
}

func WithTransaction(txHash string) StringWith {
	return StringWith{Key: "_transaction_hash", Val: txHash}
}

func WithBlockNumber(blockNumber uint64) StringWith {
	return StringWith{Key: "_block_id", Val: fmt.Sprint(blockNumber)}
}

func WithService(serviceName string) StringWith {
	return StringWith{Key: "_service", Val: serviceName}
}
