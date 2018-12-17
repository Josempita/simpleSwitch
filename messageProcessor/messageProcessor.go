package messageProcessor

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/michaelbironneau/asbclient"
)

type Config struct {
	NamespaceArg string `json:"namespaceArg"`
	KeynameArg   string `json:"keynameArg"`
	KeyvalueArg  string `json:"keyvalueArg"`
}

func getConfiguration() Config {
	conf := LoadConfiguration()
	//conf := Config{namespaceArg: os.Args[1], keynameArg: os.Args[2], keyvalueArg: os.Args[3]}
	return conf
}

func LoadConfiguration() Config {
	log.Printf("loading config from config.json")
	var config Config
	configFile, err := os.Open("config.json")
	defer configFile.Close()
	if err != nil {
		fmt.Println(err.Error())
	}
	jsonParser := json.NewDecoder(configFile)
	err = jsonParser.Decode(&config)
	if err != nil {
		fmt.Println(err.Error())
	}
	log.Printf("Namespace: " + config.NamespaceArg)
	return config
}

func PollMessages(messageCh chan string) {
	conf := getConfiguration()
	i := 1
	log.Printf("Send: %d", i)

	client := asbclient.New(asbclient.Topic, conf.NamespaceArg, conf.KeynameArg, conf.KeyvalueArg)
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
