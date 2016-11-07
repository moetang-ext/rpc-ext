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

	//methodInfo(c)
	//methodSet(c)
	methodGet(c)
}

func methodGet(c rpc.Client) {
	var result = new(client.BulkString)
	var req rpc.Param
	req.Request = []byte("bbbb")
	req.Method = client.REDIS_METHOD_GET
	timeout, err := c.Call(req, result)
	fmt.Println(result)
	fmt.Println(string(result.Data()))
	fmt.Println(timeout, err)
}

func methodSet(c rpc.Client) {
	var result string
	var req rpc.Param
	req.Request = client.SetReq{
		Key:        []byte("aaaa"),
		Value:      []byte("value_a"),
		ExpireMode: client.SET_EXPIRE_MILLIS_SECOND,
		ExpireCnt:  5,
		SetMode:    client.SET_NX,
	}
	req.Method = client.REDIS_METHOD_SET
	timeout, err := c.Call(req, &result)
	fmt.Println(result)
	fmt.Println(timeout, err)
}

func methodInfo(c rpc.Client) {
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
