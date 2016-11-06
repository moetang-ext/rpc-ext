package main

import (
	"fmt"

	"github.com/moetang-ext/rpc-ext/implements/redis/client"
	"import.moetang.info/go/nekoq-api/rpc"
)

func main() {
	rpc.InitClient()

	c, err := rpc.GetClient(client.REDIS_SERVICE_NAME)
	if err != nil {
		panic(err)
	}

	var result = new(client.BulkString)
	var req rpc.Param
	req.Method = client.REDIS_METHOD_INFO
	timeout, err := c.Call(req, result)
	fmt.Println(result)
	fmt.Println(string(result.Data()))
	fmt.Println(timeout, err)
	fmt.Println(len(result.Data()))

	timeout, err = c.Call(req, result)
	fmt.Println(len(result.Data()))

	timeout, err = c.Call(req, result)
	fmt.Println(len(result.Data()))

	timeout, err = c.Call(req, result)
	fmt.Println(len(result.Data()))

}
