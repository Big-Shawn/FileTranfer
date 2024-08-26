package main

import (
	"fmt"
	"google.golang.org/protobuf/proto"
	"iproto/frame"
	"testing"
)

func Test_resolveConnection(t *testing.T) {
	f := frame.Frame{
		Type:     0,
		Length:   0,
		Body:     nil,
		Reserved: nil,
	}

	marshal, err := proto.Marshal(&f)
	if err != nil {
		panic(err)
	}
	fmt.Println(len(marshal), proto.Size(&f))

	nf := frame.Frame{}

	err = proto.Unmarshal(marshal, &nf)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%d, %s, %d", nf.Type, nf.Body, nf.Length)
}
