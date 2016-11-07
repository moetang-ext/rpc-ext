package client

import (
	"bufio"
	"bytes"
	"net"
	"reflect"
	"strconv"
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
	REDIS_METHOD_DEL   = "DEL"

	_REDIS_SERVER_ADDRESS = "redis.server.addr"
)

func init() {
	cf := new(clientFactoryImpl)
	cf.methodMap = make(map[string]bool)
	rpc.RegisterClientFactory(REDIS_SERVICE_NAME, cf)
	rpc.RegisterMethodFactory(REDIS_SERVICE_NAME, REDIS_METHOD_INFO, reflect.TypeOf([]byte{}), reflect.TypeOf(new(BulkString)))
	rpc.RegisterMethodFactory(REDIS_SERVICE_NAME, REDIS_METHOD_SET, reflect.TypeOf(SetReq{}), reflect.TypeOf(new(string)))
	rpc.RegisterMethodFactory(REDIS_SERVICE_NAME, REDIS_METHOD_GET, reflect.TypeOf([]byte{}), reflect.TypeOf(new(BulkString)))
}

var _ rpc.ClientFactory = &clientFactoryImpl{}

type clientFactoryImpl struct {
	methodMap map[string]bool
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
	var toSend []byte
	switch param.Method {
	case REDIS_METHOD_INFO:
		toSend = Parser().MethodInfo(param.Request)
		this.Lock()
		defer this.Unlock()
		_, err := this.c.Write(toSend)
		if err != nil {
			return false, err
		}
		err = ProtocolCommonReader().ParseBulkString(this.bufIo, resultPtr.(*BulkString))
		if err != nil {
			return false, err
		}
		return true, nil
	case REDIS_METHOD_SET:
		data, err := param.Request.(SetReq).ToBytes()
		if err != nil {
			return false, err
		}
		this.Lock()
		defer this.Unlock()
		_, err = this.c.Write(data)
		if err != nil {
			return false, err
		}
		err = ProtocolCommonReader().ParseSimpleString(this.bufIo, resultPtr.(*string))
		if err != nil {
			return false, err
		}
		return true, nil
	case REDIS_METHOD_GET:
		reqData := param.Request.([]byte)
		if len(reqData) == 0 {
			return false, errorutil.New("request key is empty")
		}
		buf := new(bytes.Buffer)
		buf.WriteString("*2\r\n")
		buf.WriteString("$3\r\nGET\r\n")
		buf.WriteString("$" + strconv.Itoa(len(reqData)) + "\r\n")
		buf.Write(reqData)
		buf.WriteString("\r\n")
		this.Lock()
		defer this.Unlock()
		_, err := this.c.Write(buf.Bytes())
		if err != nil {
			return false, err
		}
		err = ProtocolCommonReader().ParseBulkString(this.bufIo, resultPtr.(*BulkString))
		if err != nil {
			return false, err
		}
		return true, nil
	default:
		return false, errorutil.New("unknown method=" + param.Method)
	}

	return false, nil
}
