package client

import (
	"reflect"

	"import.moetang.info/go/nekoq-api/future"
	"import.moetang.info/go/nekoq-api/rpc"
)

var _ rpc.ClientFactory = &clientFactoryImpl{}

type clientFactoryImpl struct {
}

func (this *clientFactoryImpl) CreateClient(conn map[string]string) (rpc.Client, error) {
	return nil, nil
}

func (this *clientFactoryImpl) PreRegisterMethod(methodName string, in reflect.Type, out reflect.Type) error {
	return nil
}

var _ rpc.Client = &clientImpl{}

type clientImpl struct {
}

func (this *clientImpl) AsyncCall(param rpc.Param, resultPtr interface{}) (future.Future, error) {
	return nil, nil
}

func (this *clientImpl) Call(param rpc.Param, resultPtr interface{}) (bool, error) {
	return false, nil
}
