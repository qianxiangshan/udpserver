package main

import (
	"dana-tech.com/wbw/logs"
	"log"
	"sync/atomic"
	"time"
	"tool/udpnet/udpserver"
)

var i uint64

func packagehandle(data *udpserver.UdpPackage) {

	defer func() {
		tmp := recover()
		if tmp != nil {
			logs.Logger.Errorf("panic in handle data %v", tmp)
		}
	}()

	if data == nil {
		log.Println("error")
	}
	atomic.AddUint64(&i, 1)
	log.Println("recv data : ", string(data.Data.RealData))
	log.Println("recv", data.Destination, data.Source)
	sorce := data.Destination
	dest := data.Source
	data.Destination = dest
	data.Source = sorce
	log.Println("send ", data.Destination, data.Source)
	err := server.Send(data)
	if err != nil {
		log.Println(err)
	}
}

var server udpserver.Server

func main() {
	var err error
	udpserver.UDPBUFFERSIZE = 50 * 1024 * 1024
	ipports := make([]string, 0)
	ipports = append(ipports, "0.0.0.0:1234")
	ipports = append(ipports, "127.0.0.1:2345")
	server, err = udpserver.NewUdpServers(ipports, 0, packagehandle)
	if err != nil {
		log.Fatal(err)
	}
	go func(i *uint64) {
		for {
			log.Println(*i)
			time.Sleep(time.Second)
		}

	}(&i)

	err = server.Start()
	if err != nil {
		log.Fatal(err)
	}

	select {}

}
