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
	connectedClients 	map[string]*net.Conn
	clientLock 			*sync.Mutex
	msgChan				chan vmail_proto.VMessage
	mongo				*MongoStore
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
		this.messageIn(message, conn)
	}
}

func (this *VMailServer) messageIn(message *vmail_proto.VMailMessage, conn net.Conn){
	switch *message.Mtype {
	case vmail_proto.MessageType_AUTH_REQUEST:
		auth_request := &vmail_proto.AuthRequest{}
		proto.Unmarshal(message.MessageData, auth_request)
		this.authenticate(*auth_request, conn)
	case vmail_proto.MessageType_VMESSAGE:
		vmessage := &vmail_proto.VMessage{}
		proto.Unmarshal(message.MessageData, vmessage)
		this.msgChan <- *vmessage
	default:
		response := &vmail_proto.Error{Text:proto.String("Message Unknown")}
		sendMessage(response, conn)
	}
}

func sendMessage(message proto.Message, conn net.Conn){
	vmail_message := &vmail_proto.VMailMessage{}
	var mtype vmail_proto.MessageType
	switch message.(type) {
	case *vmail_proto.AuthResponse:
		mtype = vmail_proto.MessageType_AUTH_RESPONSE
	case *vmail_proto.VMessage:
		mtype = vmail_proto.MessageType_VMESSAGE
	default:
		fmt.Println("Invalid message Type")
		return
	}
	vmail_message.Mtype = &mtype
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
	fmt.Printf("The mail %v\n", mail)
	this.deliverMail(mail)
}

func (this *VMailServer) deliverMail(messeges []vmail_proto.VMessage){
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

func (this *VMailServer) Init(msgChan chan vmail_proto.VMessage, mongo *MongoStore) error {
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