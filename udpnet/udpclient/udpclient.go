package main

import (
	"bytes"
	"log"
	"math/rand"
	"net"
	"strconv"
	"time"
)

var udpconn *net.UDPConn

var prepared = false

func prepare() {
	if prepared == true {
		return
	}
	prepared = true
	addr, err := net.ResolveUDPAddr("udp", "127.0.0.1:2345")
	if err != nil {
		log.Println(err)
		return
	}
	udpconn, err = net.DialUDP("udp", nil, addr)
	if err != nil {
		log.Println(err)
		return
	}
	udpconn.SetWriteBuffer(1024 * 1024 * 10)

}

var fails = 0

func Send(i int) {
	tmp := strconv.FormatInt(int64(i), 10)

	log.Println(tmp)
	_, err := udpconn.Write([]byte(tmp))
	if err != nil {
		fails++
		log.Println(err)
	}
	var recv []byte
	recv = make([]byte, 1024)
	l, err := udpconn.Read(recv)
	if err != nil {
		fails++
		log.Println(err)
	}
	result := bytes.Compare(recv[:l], []byte(tmp))
	if result != 0 {
		log.Println(err)
	}

}

func main() {
	prepare()
	log.Println("over")
	times := time.Now()
	i := rand.Int()
	for i = 0; i < 100000; i++ {
		go Send(i)
	}

	log.Println(time.Now().Sub(times))
	log.Println(i)
	log.Println(fails)
	select {}
}
