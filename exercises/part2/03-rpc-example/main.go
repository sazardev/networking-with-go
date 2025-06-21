package main

import (
	"fmt"
	"net"
	"net/rpc"
)

type Args struct {
	A, B int
}

type Arith int

func (t *Arith) Multiply(args *Args, reply *int) error {
	*reply = args.A * args.B
	return nil
}

func main() {
	arith := new(Arith)
	rpc.Register(arith)
	ln, err := net.Listen("tcp", ":1234")
	if err != nil {
		panic(err)
	}
	fmt.Println("RPC server listening on port 1234...")
	rpc.Accept(ln)
}
