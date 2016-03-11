package vmail

import (
	"fmt"
	"net"
	"strconv"
	"sync"
	"time"
	"./vmail_proto"
	"github.com/golang/protobuf/proto"
)

const VMAIL_PORT = 9989

type VMailServer struct {
	connectedClients 	int
	clientLock 			*sync.Mutex
}

func (this *VMailServer) connectionHander(conn net.Conn){
	for {
		buf := make([]byte, 1024)
		len, err := conn.Read(buf)
		message := &vmail_proto.VMailMessage{}
		proto.Unmarshal(buf[:len], message)
		if err != nil {
			time.Sleep(1000)
			continue
		}
		switch *message.Type {
		case vmail_proto.MessageType_AUTH_REQUEST:
			auth_request := &vmail_proto.AuthRequest{}
			proto.Unmarshal(message.MessageData, auth_request)
			this.authenticate(*auth_request, conn)
		default:
			response := &vmail_proto.Error{Text:proto.String("Message Unknown")}
			sendMessage(response, conn)
		}
	}
}

func sendMessage(message proto.Message, conn net.Conn){
	vmail_message := &vmail_proto.VMailMessage{}
	switch message.(type) {
	case *vmail_proto.AuthResponse:
		mtype := vmail_proto.MessageType_AUTH_RESPONSE
		vmail_message.Type = &mtype
	default:
		fmt.Println("Invalid message Type")
	}
	vmail_message.MessageData, _ = proto.Marshal(message)
	data, _ :=proto.Marshal(vmail_message)
	conn.Write(data)
}

func (this *VMailServer) authenticate(auth_request vmail_proto.AuthRequest, conn net.Conn){
	fmt.Println("Authenticating")
	username := *auth_request.Username
	password := *auth_request.Password
	response := &vmail_proto.AuthResponse{}
	if username == "" || password == ""{
		response.Success = proto.Bool(false)
		sendMessage(response, conn)
		return
	}
	// TODO do proper login here
	if username == "bahus.vel@gmail.com" && password == "password"{
		fmt.Println("Authentication success")
		response.Success = proto.Bool(true)
		sendMessage(response, conn)
		return
	}
}

func (this *VMailServer) connectionListener(ln net.Listener){
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Connection failed to accept")
		}
		fmt.Printf("Clinet %s connected\n", conn.RemoteAddr())
		this.clientLock.Lock()
		this.connectedClients++
		this.clientLock.Unlock()
		go this.connectionHander(conn)
	}
}

func (this *VMailServer) Init() error {
	fmt.Println("Initilizing the vMail Server module")
	this.clientLock = &sync.Mutex{}
	ln, err := net.Listen("tcp", ":" + strconv.Itoa(VMAIL_PORT))
	if err != nil {
		return err
	} else {
		go this.connectionListener(ln)
	}
	return nil
}