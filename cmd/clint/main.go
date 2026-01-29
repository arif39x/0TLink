package main

import (
	"github.com/hashicorp/yamux"
	"io"
	"net"
)

func main() {
	conn, _ = net.Dial("tcp", "localhost:7000") // TCP connection initiatate outbond connections are maximum allowed by firewalls so it will enable reverse tunnneling
	session, _ := yamux.Clint(conn, nil)  //Wrapper of tcp connection in a yamux Clint (i am thinking ogf one tcp socket and unlimited vvirtual coonnections)
	
	for {
		stream. _ := sesion.Accept     //Block until reliable server open a stream (one public user connection and one logi al tcp session)
		
		localApp, _ := net.Dial("tcp", "localhost:3000")
		
		go func() {    //A transparrent TCP bridge
			defer stream.Closer()
			defer localApp.Closer()
			go io.Copy(stream, localApp) //Local app → Public user
			io.Copy(localApp, stream) //Public user → Local app
		}()
	}
}