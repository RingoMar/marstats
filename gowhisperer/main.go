package main

import (
	"bufio"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/textproto"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/ringomar/gowhisperer/rqueue"
	"github.com/ringomar/gowhisperer/utils"
)

var giftOrigin string
var queue = &rqueue.Queue{}
var SubLogger *log.Logger
var upStreamURL string = "localhost"

const SubEvent = iota

type USERNOTICE struct {
	ID         string
	Channel    string
	Name       string
	SubMethod  string
	SubAmount  string
	GiftAmount string
	SubPlan    string
	Created    string
}

func init() {
	file, err := os.OpenFile("usernotice.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Failed to open usernotice file:", err)
	}

	SubLogger = log.New(file, "", 0)
}

func main() {
	SetupCloseHandler()

	if os.Getenv("GO_ENV") == "production" {
		upStreamURL = "mongoserver"
	}

	// Edit channels you want to join
	channels := []string{"kaicenat", "feelssunnyman"}
	timeout := 200 * time.Second
	chat := connect()

	for _, v := range channels {
		chat.join(v)
	}

	go rqueue.MonitorQueue(queue, timeout)
	for message := range chat.messages {
		if strings.HasPrefix(message, "PING") {
			chat.send("PONG :tmi.twitch.tv")
			continue
		}

		ircTranslate, _ := utils.ParseIRCMessage(message)

		if ircTranslate.Command == "USERNOTICE" {
			// Loggin system
			stype := ircTranslate.Tags["msg-id"]
			if !(stype == "resub" || stype == "sub" || stype == "subgift" || stype == "submysterygift" || stype == "giftpaidupgrade" || stype == "anongiftpaidupgrade") {
				continue
			}

			SubLogger.Println(message)

			b, err := json.Marshal(ircTranslate)
			if err != nil {
				log.Println(err)
				return
			}
			if !(ircTranslate.Tags["msg-id"] == "subgift") {
				log.Println("\n\n<====>", string(b))

			}
		}

		if ircTranslate.Command == "USERNOTICE" {
			stype := ircTranslate.Tags["msg-id"]
			if !(stype == "resub" || stype == "sub" || stype == "subgift" || stype == "submysterygift" || stype == "giftpaidupgrade" || stype == "anongiftpaidupgrade") {
				continue
			}

			if stype == "submysterygift" {
				giftOrigin = ircTranslate.Tags["msg-param-community-gift-id"]
			}

			if stype == "subgift" {
				if ircTranslate.Tags["msg-param-community-gift-id"] == giftOrigin {
					log.Println("????? 			IGNORING SUBGIFT PART OF SUB MYSTERY")
					continue
				}
			}

			if stype == "sub" {
				data := url.Values{}

				data.Add("ID", ircTranslate.Tags["id"])
				data.Add("Channel", ircTranslate.Params[0])
				data.Add("Name", ircTranslate.Tags["display-name"])
				data.Add("SubMethod", ircTranslate.Tags["msg-id"])
				data.Add("SubAmount", ircTranslate.Tags["sub-amount"])
				data.Add("giftAmount", ircTranslate.Tags["msg-param-mass-gift-count"])
				data.Add("subPlan", ircTranslate.Tags["msg-param-sub-plan"])
				data.Add("Created", ircTranslate.Tags["tmi-sent-ts"])

				uploadOne(data)
			} else {

				queue.AddQueue(ircTranslate.Tags["id"],
					ircTranslate.Params[0],
					ircTranslate.Tags["display-name"],
					ircTranslate.Tags["msg-id"],
					ircTranslate.Tags["sub-amount"],
					ircTranslate.Tags["msg-param-mass-gift-count"],
					ircTranslate.Tags["msg-param-sub-plan"],
					ircTranslate.Tags["tmi-sent-ts"])
			}
		}

	}
}

func SetupCloseHandler() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("\n\n\r- Ctrl+C pressed in Terminal")
		os.Exit(0)
	}()
}

type irc struct {
	conn     net.Conn
	messages chan string
}

func connect() *irc {
	fmt.Printf("[%s][RIN-Connect] Dialing TCP server.\n", string(time.Now().Format(time.Stamp)))

	conn, err := tls.Dial("tcp", "irc.chat.twitch.tv:443", nil)
	if err != nil {
		log.Fatal("cannot connect to twitch irc server", err)
	}
	i := &irc{
		conn:     conn,
		messages: make(chan string, 10),
	}
	i.send("PASS oauth:1231231")
	i.send("NICK justinfan123")
	i.send("CAP REQ twitch.tv/membership")
	i.send("CAP REQ twitch.tv/commands")
	i.send("CAP REQ twitch.tv/tags")
	go i.read()

	fmt.Printf("[%s][RIN-Connect] connected to Twitch irc server\n", string(time.Now().Format(time.Stamp)))
	return i
}

func (i *irc) join(channel string) {
	i.send("JOIN #" + strings.ToLower(channel))
	fmt.Printf("[%s][RIN-Connect] Joined channel: %s\n", string(time.Now().Format(time.Stamp)), channel)

}

func (i *irc) send(msg string) {
	_, err := i.conn.Write([]byte(msg + "\r\n"))
	if err != nil {
		log.Fatal("Disconnected from twitch irc server", err)
	}
}

func (i *irc) read() {
	reader := bufio.NewReader(i.conn)
	tp := textproto.NewReader(reader)
	for {
		message, err := tp.ReadLine()
		if err != nil {
			log.Fatal("Disconnected from twitch irc server", err)
		}
		i.messages <- message
	}
}

func uploadOne(data url.Values) {
	url := fmt.Sprintf("http://%s:5284/api/addone", upStreamURL)
	resp, err := http.PostForm(url, data)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Error: received status code %d", resp.StatusCode)
	}
}
