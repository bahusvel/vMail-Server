package main

import (
	"./vmail"
	"./vmail/vproto"
)

func main(){
	mongo := &vmail.MongoStore{}
	mongo.Init("bahus.com", "192.168.99.100")
	msgChannels := vmail.TransportChannels{VServer:make(chan vproto.VMessage)}
	vserver := &vmail.VMailServer{}
	go vserver.Init(msgChannels.VServer, mongo)
	vmail.MessagePlane(msgChannels, mongo)
}
