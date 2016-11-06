package client

import (
	"bufio"
	"io"
	"strconv"
)

type BulkString struct {
	null bool
	data []byte
}

func (this *BulkString) NullBulkString() bool {
	return this.null
}

func (this *BulkString) Data() []byte {
	return this.data
}

func ProtocolCommonReader() commonReader {
	return commonReaderImpl
}

var commonReaderImpl = commonReader{}

type commonReader struct {
}

//TODO return value may be: BulkString(NullString), SimpleString, Error

func (this commonReader) ParseBulkString(r *bufio.Reader, resp *BulkString) error {
	var length int
	data, err := r.ReadBytes('\n')
	if err != nil {
		return err
	}
	s, err := strconv.Atoi(string(data[1 : len(data)-2]))
	if err != nil {
		return err
	}
	length = s
	if length == -1 {
		resp.null = true
		return nil
	}
	resp.null = false
	if length == 0 {
		r.Discard(2)
		return nil
	} else {
		result := make([]byte, length)
		_, err := io.ReadFull(r, result)
		if err != nil {
			return err
		} else {
			r.Discard(2)
			resp.data = result
			return nil
		}
	}
}

func (this commonReader) ParseSimpleString(r *bufio.Reader, resp *string) error {
	data, err := r.ReadBytes('\n')
	if err != nil {
		return err
	}
	*resp = string(data[1 : len(data)-2])
	return nil
}

//================

type parser struct {
}

var parserImpl = parser{}

func Parser() parser {
	return parserImpl
}

func (this parser) MethodInfo(req interface{}) []byte {
	return []byte("INFO\r\n")
}
