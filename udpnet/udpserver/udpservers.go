package udpserver

import (
	"dana-tech.com/wbw/logs"
	"strconv"
	"strings"
)

//servers using to listening multiport at the same time
type udpServers struct {
	udpservers map[int]*udpServer
}

//ipports must more than 1 buffsize motethan 1
func NewUdpServers(ipports []string, buffsize uint16, handle PackageHandler) (*udpServers, error) {

	if len(ipports) <= 1 || buffsize < 0 {
		return nil, ServerArgsError
	}

	servers := new(udpServers)
	servers.udpservers = make(map[int]*udpServer)
	for _, ipport := range ipports {
		ipportslice := strings.Split(ipport, ":")
		if len(ipportslice) != 2 {
			return nil, ServerArgsError
		}
		port, err := strconv.Atoi(ipportslice[1])
		if err != nil || port <= 0 || port > 65535 {
			return nil, ServerArgsError
		}
		servers.udpservers[port], _ = Newudpserver(ipport, buffsize, handle)
	}
	return servers, nil
}

func (servers *udpServers) Start() error {

	if len(servers.udpservers) == 0 {
		return ServerArgsNotInit
	}
	result := make(chan bool, 1)
	var ok bool
	for _, server := range servers.udpservers {
		go func() {
			err := server.Start()
			if err != nil {
				logs.Logger.Errorf("server %s  start   error %s", server.ipport, err.Error())
				result <- false
				return
			} else {
				result <- true

			}
			logs.Logger.Tracef("server %s  start  ", server.ipport)
		}()
		select {
		case ok = <-result:
			if ok == false {
				break
			}
		}
	}
	if !ok {
		for _, server := range servers.udpservers {
			if server.conn != nil {
				server.conn.Close()
			}
		}
		return ServerStartError
	}
	return nil
}

func (servers *udpServers) Send(data *UdpPackage) error {
	if data == nil {
		return NilPointer
	}
	server, has := servers.udpservers[data.Source.Port]
	if !has {
		return ServerNotExist
	}

	return server.Send(data)

}
