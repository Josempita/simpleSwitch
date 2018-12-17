package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"simpleSwitch/messageProcessor"
	"time"

	"github.com/stianeikeland/go-rpio"
)

var (
	// Use mcu pin 21 (gpio21) corresponds to physical pin 40 on the pi
	pin = rpio.Pin(21)
)

func main() {

	// Open and map memory to access gpio, check for errors
	if err := rpio.Open(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Unmap gpio memory when done
	defer rpio.Close()

	// Set pin to output mode
	pin.Output()

	//lets just send a signal to show that the service is up
	pin.High()
	time.Sleep(time.Second / 2)
	pin.Low()
	time.Sleep(time.Second / 2)
	pin.High()
	time.Sleep(time.Second / 2)
	pin.Low()

	http.HandleFunc("/switchoff", func(w http.ResponseWriter, r *http.Request) {

		handleSwitchOff(w)
	})

	http.HandleFunc("/switchon", func(w http.ResponseWriter, r *http.Request) {

		handleSwitchOn(w)
	})
	messageCh := make(chan string)
	go messageProcessor.PollMessages(messageCh)
	go messageListener(messageCh)
	log.Fatal(http.ListenAndServe(":8083", nil))
}

func handleSwitchOff(w http.ResponseWriter) {

	mapResponse, _ := json.Marshal("switchOff")
	fmt.Fprintf(w, "%q", string(mapResponse))
	pin.Low()
}

func handleSwitchOn(w http.ResponseWriter) {

	mapResponse, _ := json.Marshal("switchOn")
	time.Sleep(time.Second)
	fmt.Fprintf(w, "%q", string(mapResponse))
	pin.High()
}

func messageListener(messageCh chan string) {
	for {
		//Block untile message received from channel
		msg := <-messageCh
		if msg == "on" {
			fmt.Println("On State Received")
			pin.High()
		} else if msg == "off" {
			fmt.Println("Off State Received")
			pin.Low()
		}

	}
}
