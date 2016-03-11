package vmail
import (
	"fmt"
	"./vmail_proto"
)

func vmailIn(vmail *vmail_proto.VMessage){
	fmt.Println(vmail)
}