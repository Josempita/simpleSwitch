package messageProcessor

import (
	"log"
	"time"

	"github.com/michaelbironneau/asbclient"
)

func PollMessages(messageCh chan string) {

	i := 0
	log.Printf("Send: %d", i)
	namespace := "[Azure Data Service Namespace Here]"
	keyname := "[Azure data service key Name]"
	keyvalue := "[Azure data service Key]"

	client := asbclient.New(asbclient.Topic, namespace, keyname, keyvalue)
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
