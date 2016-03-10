package main

import (
	"./vmail"
	"time"
)

func main(){
	var vmail = &vmail.VMailServer{}
	vmail.Init()
	for {
		time.Sleep(1000000)
	}
}
