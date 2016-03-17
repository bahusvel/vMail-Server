package main

import (
	"./vmail"
	"./vmail/vproto"
)

const STORAGE_DOMAIN = "bahus.com"

func main(){
	mongo := &vmail.MongoStore{}
	mongo.Init(STORAGE_DOMAIN, "192.168.99.100")
	defer mongo.Close()
	msgChannels := vmail.TransportChannels{VServer:make(chan vproto.VMessage)}
	vserver := &vmail.VMailServer{}
	go vserver.Init(msgChannels.VServer, mongo)
	mplane := vmail.MessagePlane{Channels: msgChannels, Mongo: mongo, StorageDomain: STORAGE_DOMAIN}
	mplane.Start()
}
