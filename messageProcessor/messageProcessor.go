package messageProcessor

import (
	"log"
	"os"
	"time"

	"github.com/michaelbironneau/asbclient"
)

type config struct {
	namespaceArg, keynameArg, keyvalueArg string
}

func getConfiguration() config {
	conf := config{namespaceArg: os.Args[1], keynameArg: os.Args[2], keyvalueArg: os.Args[3]}
	return conf
}

func PollMessages(messageCh chan string) {

	conf := getConfiguration()
	i := 0
	log.Printf("Send: %d", i)

	client := asbclient.New(asbclient.Topic, conf.namespaceArg, conf.keynameArg, conf.keyvalueArg)
	client.SetSubscription("stateSubscription")

	for {
		msg, err := client.PeekLockMessage("state", 30)

		if err != nil {
			log.Printf("Peek error: %s", err)
		} else {
			if msg != nil {
				message := string(msg.Body)
				messageCh <- message
				log.Printf("Peeked message: '%s'", message)

				err = client.DeleteMessage(msg)
				if err != nil {
					log.Printf("Delete error: %s", err)

				}
			}
		}
		time.Sleep(time.Millisecond * 200)
	}
}
