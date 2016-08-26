package udpserver

import (
	"errors"
	"net"
)

var (
	ServerArgsNotInit = errors.New("server args Not init")
	ServerArgsError   = errors.New("server args error")
	ServerNotExist    = errors.New("server port not exsit")
	ServerStartError  = errors.New("server start error")
	SendTimeout       = errors.New("send time out 3s due to full buff")
	NilPointer        = errors.New("nil pointer")
)

const (
	UdpSendTimeout = 3
)

var (
	// UDP buffer size to revice package B
	UDPBUFFERSIZE = 1 * 1024 * 1024
	// max package size using to put package  B less than MTU
	MaxPackageSize = 1500
)

type PackageHandler func(*UdpPackage)

type Server interface {
	//发送数据
	Send(*UdpPackage) error
	//初始化.传入,监听ip:port格式的数据
	Start() error
}

//data package  ,using for net package
type Data struct {
	//store real data
	RealData []byte
	//store len(data) B,max uint16
	Data_len uint16
}

//udp pacakge
type UdpPackage struct {
	Destination *net.UDPAddr
	Source      *net.UDPAddr
	Data
}
