package client

import (
	"bufio"
	"net"
	"reflect"
	"sync"
	"time"

	"import.moetang.info/go/nekoq-api/errorutil"
	"import.moetang.info/go/nekoq-api/future"
	"import.moetang.info/go/nekoq-api/rpc"
)

const (
	REDIS_SERVICE_NAME = "net.moetang.client.redis"
	REDIS_METHOD_INFO  = "INFO"
	REDIS_METHOD_GET   = "GET"
	REDIS_METHOD_SET   = "SET"

	_REDIS_SERVER_ADDRESS = "redis.server.addr"
)

func init() {
	rpc.RegisterClientFactory(REDIS_SERVICE_NAME, new(clientFactoryImpl))
	rpc.RegisterMethodFactory(REDIS_SERVICE_NAME, REDIS_METHOD_INFO, reflect.TypeOf([]byte{}), reflect.TypeOf(new(BulkString)))
}

var _ rpc.ClientFactory = &clientFactoryImpl{}

type clientFactoryImpl struct {
}

func (this *clientFactoryImpl) CreateClient(config map[string]string) (rpc.FullClient, error) {
	addr, ok := config[_REDIS_SERVER_ADDRESS]
	if ok {
		// single address mode
		c, err := net.DialTimeout("tcp", addr, 10*time.Second)
		if err != nil {
			return nil, err
		}
		cl := new(clientSingleImpl)
		cl.c = c.(*net.TCPConn)
		cl.bufIo = bufio.NewReaderSize(c, 1024)

		return cl, nil
	} else {
		//TODO cluster mode with service naming
		return nil, errorutil.New("no cluster mode supported")
	}
}

func (this *clientFactoryImpl) PreRegisterMethod(methodName string, in reflect.Type, out reflect.Type) error {
	return nil
}

var _ rpc.Client = &clientSingleImpl{}

type clientSingleImpl struct {
	sync.Mutex
	c     *net.TCPConn
	bufIo *bufio.Reader
}

func (this *clientSingleImpl) AsyncCall(param rpc.Param, resultPtr interface{}) (future.Future, error) {
	return nil, errorutil.New("async method not supported")
}

func (this *clientSingleImpl) Call(param rpc.Param, resultPtr interface{}) (bool, error) {
	this.Lock()
	defer this.Unlock()

	var toSend []byte
	switch param.Method {
	case "INFO":
		toSend = Parser().MethodInfo(param.Request)
		_, err := this.c.Write(toSend)
		if err != nil {
			return false, err
		}
		err = ProtocolCommonReader().ParseBulkString(this.bufIo, resultPtr.(*BulkString))
		return true, err
	default:
		return false, errorutil.New("unknown method=" + param.Method)
	}

	return false, nil
}
