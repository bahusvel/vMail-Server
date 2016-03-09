package main

import (
	"./proto"
	"time"
)

func main(){
	var vmail = &proto.VMailServer{}
	vmail.Init()
	for {
		time.Sleep(10)
	}
}
