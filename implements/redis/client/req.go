package client

import (
	"bytes"
	"strconv"

	"import.moetang.info/go/nekoq-api/errorutil"
)

type SET_MODE int

const (
	SET_NONE SET_MODE = 0
	SET_NX   SET_MODE = 1
	SET_XX   SET_MODE = 2
)

type SET_EXPIRE_MODE int

const (
	SET_EXPIRE_NONE          = 0
	SET_EXPIRE_SECOND        = 1
	SET_EXPIRE_MILLIS_SECOND = 2
)

type SetReq struct {
	Key        []byte
	Value      []byte
	ExpireCnt  int
	ExpireMode SET_EXPIRE_MODE
	SetMode    SET_MODE
}

func (this SetReq) ToBytes() ([]byte, error) {
	expire := false
	if len(this.Key) == 0 {
		return nil, errorutil.New("key length is 0")
	}
	paramCnt := 3
	if this.ExpireCnt > 0 && this.ExpireMode > 0 {
		paramCnt += 2
		expire = true
	}
	if this.SetMode > 0 {
		paramCnt += 1
	}
	buf := new(bytes.Buffer)
	buf.WriteString("*" + strconv.Itoa(paramCnt) + "\r\n")

	buf.WriteString("$3\r\nSET\r\n")

	buf.WriteString("$" + strconv.Itoa(len(this.Key)) + "\r\n")
	buf.Write(this.Key)
	buf.WriteString("\r\n")

	buf.WriteString("$" + strconv.Itoa(len(this.Value)) + "\r\n")
	if len(this.Value) > 0 {
		buf.Write(this.Value)
	}
	buf.WriteString("\r\n")

	if expire {
		switch this.ExpireMode {
		case SET_EXPIRE_SECOND:
			buf.WriteString("$2\r\nEX\r\n")
			ss := strconv.Itoa(this.ExpireCnt)
			buf.WriteString("$" + strconv.Itoa(len(ss)) + "\r\n")
			buf.WriteString(ss + "\r\n")
		case SET_EXPIRE_MILLIS_SECOND:
			buf.WriteString("$2\r\nPX\r\n")
			ss := strconv.Itoa(this.ExpireCnt)
			buf.WriteString("$" + strconv.Itoa(len(ss)) + "\r\n")
			buf.WriteString(ss + "\r\n")
		}
	}

	switch this.SetMode {
	case SET_NX:
		buf.WriteString("$2\r\nNX\r\n")
	case SET_XX:
		buf.WriteString("$2\r\nXX\r\n")
	}

	return buf.Bytes(), nil
}
