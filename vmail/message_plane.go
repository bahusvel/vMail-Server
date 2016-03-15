package vmail
import (
	"fmt"
	"./vmail_proto"
	"strings"
	"gopkg.in/mgo.v2/bson"
)

type TransportChannels struct {
	VServer chan vmail_proto.VMessage
}

func getDomain(address string) string{
	return address[strings.Index(address, "@") + 1:]
}

func allRecepients(message *vmail_proto.VMessage) []string {
	return append(message.GetReceivers(), message.GetHiddenReceivers()...)
}

func newMail(mongo *MongoStore, username string) []vmail_proto.VMessage {
	return mongo.msgsFromDB(bson.M{}, username+"_rx")
}

func MessagePlane(channels TransportChannels, mongo *MongoStore) {
	for {
		var msg vmail_proto.VMessage
		select {
		case msg = <- channels.VServer:
			fmt.Println("vServer Message")
		}
		mongo.toDB(msg)
	}
}