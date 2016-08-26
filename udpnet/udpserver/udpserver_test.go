package udpserver

import (
	"log"
	"testing"
)

func packagehandle(data *UdpPackage) {
	if data == nil {
		log.Println("error")
	}
	log.Println(data)

}

func TestNewudpserver(t *testing.T) {

	udpserver, err := Newudpserver("0.0.0.0:1234", 1000, packagehandle)
	if err != nil {
		t.Fatal(err)
	}

	if _, ok := udpserver.(*udpServer); !ok {
		t.Log(ok)
		t.Fail()
	}

}

func TestStart(t *testing.T) {

	var udp udpServer
	var p *UdpPackage
	err := udp.Send(p)
	if err != nil {
		t.Fatal(err)
	}

}
