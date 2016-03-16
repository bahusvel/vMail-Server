package vmail

import (
	"fmt"
	"net"
	"strconv"
	"sync"
	"time"
	"./vproto"
	"github.com/golang/protobuf/proto"
	"encoding/binary"
)

const VMAIL_PORT = 9989

type VMailServer struct {
	connectedClients 	map[string]*net.Conn
	clientLock 			*sync.Mutex
	msgChan				chan vproto.VMessage
	mongo				*MongoStore
}

func (this *VMailServer) connectionHander(conn net.Conn){
	for {
		buf := make([]byte, 1024)
		lenbuf := make([]byte, 4)
		var msgData []byte = make([]byte, 0)
		lnread, err := conn.Read(lenbuf)
		if lnread == 4 && err == nil{
			msgLength := int(binary.LittleEndian.Uint32(lenbuf))
			for len(msgData) < msgLength {
				if msgLength - len(msgData) >= len(buf) {
					lnread, err = conn.Read(buf)
				} else {
					buf = make([]byte, msgLength - len(msgData))
					lnread, err = conn.Read(buf)
				}
				msgData = append(msgData, buf[:lnread]...)
			}
			message := &vproto.VMailMessage{}
			proto.Unmarshal(msgData, message)
			if err != nil {
				time.Sleep(1000)
				continue
			}
			this.messageIn(message, conn)
		}

	}
}

func (this *VMailServer) messageIn(message *vproto.VMailMessage, conn net.Conn){
	switch *message.Mtype {
	case vproto.MessageType_AUTH_REQUEST:
		auth_request := &vproto.AuthRequest{}
		proto.Unmarshal(message.MessageData, auth_request)
		this.authenticate(*auth_request, conn)
	case vproto.MessageType_VMESSAGE:
		vmessage := &vproto.VMessage{}
		proto.Unmarshal(message.MessageData, vmessage)
		this.msgChan <- *vmessage
	default:
		response := &vproto.Error{Text:proto.String("Message Unknown")}
		sendMessage(response, conn)
	}
}

func sendMessage(message proto.Message, conn net.Conn){
	vmail_message := &vproto.VMailMessage{}
	var mtype vproto.MessageType
	switch message.(type) {
	case *vproto.AuthResponse:
		mtype = vproto.MessageType_AUTH_RESPONSE
	case *vproto.VMessage:
		mtype = vproto.MessageType_VMESSAGE
	default:
		fmt.Println("Invalid message Type")
		return
	}
	vmail_message.Mtype = &mtype
	vmail_message.MessageData, _ = proto.Marshal(message)
	data, _ :=proto.Marshal(vmail_message)
	length := make([]byte, 4)
	binary.LittleEndian.PutUint32(length, uint32(len(data)))
	conn.Write(length)
	conn.Write(data)
}

func (this *VMailServer) authenticate(auth_request vproto.AuthRequest, conn net.Conn){
	fmt.Println("Authenticating")
	username := *auth_request.Username
	password := *auth_request.Password
	response := &vproto.AuthResponse{}
	if username == "" || password == ""{
		response.Success = proto.Bool(false)
		sendMessage(response, conn)
		return
	}
	// TODO do proper login here
	if username == "bahus.vel@bahus.com" && password == "password"{
		fmt.Println("Authentication success")
		response.Success = proto.Bool(true)
		this.clientLock.Lock()
		this.connectedClients[username] = &conn
		this.clientLock.Unlock()
		sendMessage(response, conn)
		// login hook here
		go this.loginHook(username)
	} else {
		response.Success = proto.Bool(false)
		sendMessage(response, conn)
	}
}

func (this *VMailServer) loginHook(username string){
	mail := newMail(this.mongo, username)
	this.deliverMail(mail)
}

func (this *VMailServer) deliverMail(messeges []vproto.VMessage){
	for _, message := range messeges {
		recipients := allRecepients(&message)
		for _, recipient := range recipients{
			if conn, ok := this.connectedClients[recipient]; ok {
				sendMessage(&message, *conn)
			}
		}
	}
}

func (this *VMailServer) connectionListener(ln net.Listener){
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Connection failed to accept")
		}
		fmt.Printf("Clinet %s connected\n", conn.RemoteAddr())
		go this.connectionHander(conn)
	}
}

func (this *VMailServer) Init(msgChan chan vproto.VMessage, mongo *MongoStore) error {
	fmt.Println("Initilizing the vMail Server module")
	this.msgChan = msgChan
	this.mongo = mongo
	this.clientLock = &sync.Mutex{}
	this.connectedClients = make(map[string]*net.Conn)
	ln, err := net.Listen("tcp", ":" + strconv.Itoa(VMAIL_PORT))
	if err != nil {
		return err
	} else {
		go this.connectionListener(ln)
	}
	return nil
}