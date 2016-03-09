package proto
import (
	"fmt"
	"net"
	"strconv"
	"sync"
)

const VMAIL_PORT = 9989

type VMailServer struct {
	connectedClients 	int
	clientLock 			*sync.Mutex
}

func (this *VMailServer) connectionHander(conn *net.TCPConn){

}

func (this *VMailServer) connectionListener(ln *net.TCPListener){
	for {
		conn, err := ln.AcceptTCP()
		if err != nil {
			fmt.Println("Connection failed to accept")
		}
		this.clientLock.Lock()
		this.connectedClients++
		this.clientLock.Unlock()
		go this.connectionHander(conn)
	}
}

func (this *VMailServer) Init() error {
	fmt.Println("Initilizing the vMail Server module")
	this.clientLock = &sync.Mutex{}
	tcpAddr, _ := net.ResolveTCPAddr("ip4", ":" + strconv.Itoa(VMAIL_PORT))
	ln, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		return err
	} else {
		go this.connectionListener(ln )
	}
	return nil
}