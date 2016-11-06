package client_test

import "import.moetang.info/go/nekoq-api/rpc"

func ExampleDemo() {
	rpc.RegisterClientFactory()
	rpc.RegisterMethodFactory()
}
