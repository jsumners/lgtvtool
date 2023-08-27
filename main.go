package main

import (
	"crypto/tls"
	"embed"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"strings"

	"github.com/gorilla/websocket"
)

//go:embed resources/*
var resources embed.FS

var log *slog.Logger
var messages chan *ReceivedMessage
var interrupt chan os.Signal

func main() {
	log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	slog.SetDefault(log)

	var tvIpAddr string
	flag.Parse()
	args := flag.Args()
	if len(args) < 1 {
		log.Error("Must supply television's IP address as a parameter.")
		os.Exit(1)
	}
	tvIpAddr = args[0]

	interrupt = make(chan os.Signal)
	messages = make(chan *ReceivedMessage)
	defer close(interrupt)
	defer close(messages)

	signal.Notify(interrupt, os.Interrupt)

	// We need a custom dialer because the TLS certificate presented by the
	// television is not a valid certificate.
	dialer := websocket.Dialer{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	log.Debug("dialing television...")
	tvUri := fmt.Sprintf("wss://%s:3001/", tvIpAddr)
	conn, _, err := dialer.Dial(tvUri, nil)
	if err != nil {
		log.Error("failed to dial television", "error", err)
		panic(err)
	}
	defer conn.Close()
	go messageReceivedHandler(conn)

	// First we need to send a handshake message. This handshake takes one of two
	// forms:
	// 1. We have a locally stored client key. In this case, we send the register
	// message with our client key embedded, and our app should then be
	// authorized.
	// 2. We do not have a locally stored client key. In this case, we need to
	// send the register message without a key, get a response back _with_ the
	// key, send a new register using the acquired key, and then our app should
	// be authorized.
	registerJSON, err := resources.ReadFile("resources/lg_register.json")
	if err != nil {
		log.Error("Could not load handshake payload.", "error", err)
		os.Exit(1)
	}

	var clientKey *string
	clientKey = readClientKey()
	if clientKey != nil {
		registerJSON = []byte(strings.Replace(string(registerJSON), "CLIENTKEYGOESHERE", *clientKey, 1))
		sendTextMessage(conn, registerJSON)

		// TV will tell us that it received the message and we do not care.
		_ = <-messages
	} else {
		sendTextMessage(conn, registerJSON)
		// First response is ignorable.
		response := <-messages
		// Second response contains our key.
		response = <-messages

		payload := &PayloadRegisterKey{}
		err = json.Unmarshal(response.Payload, payload)
		if err != nil {
			log.Error("Could not parse client key response.", "error", err)
			os.Exit(1)
		}

		saveClientKey(&payload.ClientKey)
	}
	log.Debug("...television dialed")

	log.Debug("sending request for service menu")
	//msg := `{"id": "something_0", "type": "request", "uri":"ssap://com.webos.applicationManager/launch", "payload": { "id": "com.webos.app.factorywin", "params": { "id": "executeFactory", "irKey": "inStart" }}}`
	cmd, err := json.Marshal(NewServiceMenuCommand(0))
	if err != nil {
		log.Error("Failed to send service menu command.", "error", err)
		os.Exit(1)
	}
	sendTextMessage(conn, cmd)

	log.Debug("entering main application loop")
mainLoop:
	for {
		select {
		case receivedMessage := <-messages:
			if receivedMessage == nil {
				// Will never happen. This is just a stub in case we want to do
				// something extra with the message in the future (e.g. inspect the type
				// and perform other actions).
			}
			continue

		case <-interrupt:
			conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			break mainLoop
		}
	}
}

func messageReceivedHandler(connection *websocket.Conn) {
	for {
		_, msg, err := connection.ReadMessage()
		if err != nil {
			if errors.Is(err, net.ErrClosed) != true {
				log.Error("Error receiving message.", "error", err)
			}
			return
		}

		log.Debug("received message", "message", msg)
		rm := &ReceivedMessage{}
		json.Unmarshal(msg, rm)
		messages <- rm
	}
}

func readClientKey() *string {
	file, err := os.ReadFile("./client-key.txt")
	if err != nil {
		log.Debug("could not read client-key.txt", "error", err)
		return nil
	}
	key := string(file)
	return &key
}

func saveClientKey(key *string) {
	err := os.WriteFile("./client-key.txt", []byte(*key), 0666)
	if err != nil {
		log.Error("Could not write client-key.txt", "error", err)
	}
}

func sendTextMessage(conn *websocket.Conn, msg []byte) {
	//log.Debug("sending message", "message", string(msg))
	conn.WriteMessage(websocket.TextMessage, msg)
}
