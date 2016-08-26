package udpserver

import (
	"dana-tech.com/wbw/logs"
	"net"
	"time"
)

type udpServer struct {
	conn        *net.UDPConn
	handle      PackageHandler
	ipport      string
	packagebuff chan *UdpPackage
}

func Newudpserver(ipport string, buffsize uint16, handle PackageHandler) (*udpServer, error) {
	udpserver := new(udpServer)
	udpserver.ipport = ipport
	udpserver.handle = handle
	if buffsize == 0 {
		udpserver.packagebuff = make(chan *UdpPackage, 1)
		return udpserver, nil
	} else {
		udpserver.packagebuff = make(chan *UdpPackage, buffsize)
		return udpserver, nil
	}
}

func (udpserver *udpServer) recv() {
	var packagebuffer []byte
	packagebuffer = make([]byte, MaxPackageSize)
	for {
		packagesize, raddr, err := udpserver.conn.ReadFromUDP(packagebuffer)
		if err != nil {
			logs.Logger.Errorf("udp server %s  readfromudp error %s", udpserver.ipport, err.Error())
			continue
		}
		//
		var data = new(UdpPackage)
		//这里如果导致gc 过高,可以使用pool来减少,不同大小等级的pool,依据协议数据包的大小
		data.RealData = make([]byte, packagesize)
		data.Data_len = uint16(packagesize)
		copy(data.RealData, packagebuffer[:packagesize])
		data.Source = raddr
		data.Destination, err = net.ResolveUDPAddr("udp", udpserver.ipport)
		if err != nil {
			logs.Logger.Errorf("udp server %s  readfromudp error %s", udpserver.ipport, err.Error())
			continue
		}
		go udpserver.handle(data)
	}

}

//阻塞函数,
func (udpserver *udpServer) Start() error {
	var err error
	go udpserver.send()
	if len(udpserver.ipport) == 0 {
		return ServerArgsNotInit
	}
	udpaddr, err := net.ResolveUDPAddr("udp", udpserver.ipport)
	if err != nil {
		logs.Logger.Errorf("ResolveUDPAddr %s  error %s", udpserver.ipport, err.Error())
		return nil
	}
	udpserver.conn, err = net.ListenUDP("udp", udpaddr)
	if err != nil {
		logs.Logger.Errorf("ListenUDP %s  error %s", udpserver.ipport, err.Error())
		return err
	}
	//缓冲区大小1M,支持1w个心跳包
	udpserver.conn.SetReadBuffer(UDPBUFFERSIZE)
	//单个包最大1500B,MTU的值

	go udpserver.recv()
	return err
}

//单独启动的用于发送处理的函数,内部使用,不对外公开
func (udpserver *udpServer) send() {
	var data *UdpPackage
	for {
		select {
		case data = <-udpserver.packagebuff:
			_, err := udpserver.conn.WriteToUDP(data.RealData, data.Destination)
			if err != nil {
				logs.Logger.Errorf("send to %v error %v", data.Destination, err)
			}
		}
	}
}

//用户级别的发送函数,放入缓冲区,放入缓冲区即认为发送成功
func (udpserver *udpServer) Send(data *UdpPackage) error {
	if data == nil {
		return NilPointer
	}

	select {
	case <-time.After(time.Second * UdpSendTimeout):
		return SendTimeout
	case udpserver.packagebuff <- data:
		//未初始化的channle,是阻塞的
		return nil
	}
}

func (udpserver *udpServer) GetInfo() string {
	return "listen" + udpserver.ipport
}
