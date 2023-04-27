package main

import (
	"fmt"
	"net"
	"net/rpc"
	"ppgti/remotelist/pkg"
)

func main() {
	list := remotelist.NewPersistentRemoteList()
	rpcs := rpc.NewServer()
	rpcs.Register(list)

	l, err := net.Listen("tcp", "[localhost]:5000")
	defer l.Close()

	if err != nil {
		fmt.Println("Something went wrong while listing to port: %w\n")
		return 
	}

	for {
		conn, err := l.Accept()
		if err == nil {
			go rpcs.ServeConn(conn)
		} else {
			break
		}
	}
}
