package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type USERNOTICE struct {
	ID         string `json:"id"`
	Channel    string `json:"channel"`
	Name       string `json:"name"`
	SubMethod  string `json:"subMethod"`
	SubAmount  string `json:"subAmount"`
	GiftAmount string `json:"giftAmount"`
	SubPlan    string `json:"subPlan"`
	Created    string `json:"created"`
}

func main() {
	for {

		if !IsOnline() {
			log.Println("Internet Connected. Welcome to the Party")
			break
		}
	}

	connect()
	introScreen := `
.--------------------------.		
|                      .   |
|                     / V\ |
|                   / .  / |
|                  <<   |  |
|                  /    |  |
|__________      /      |  |
||MAR      \   /        |  |
||DATABASE | /    \  \ /   |
||         |(      ) | |   |
||  ________|   _/_  | |   |
|<__________\______)\__)   |
.__________________________.`

	fmt.Printf("%s%s%s\n\n", "\033[34m", introScreen, "\033[0m")
	handleRequest()
}

func IsOnline() bool {
	_, err := http.Get("https://mongodb.com")
	return err != nil
}

var client *mongo.Client

// Connecting with the database (MongoDB)
func connect() {
	// Connect using the srv url)
	clientOptions := options.Client().ApplyURI(os.Getenv("srv_link"))
	client, _ = mongo.NewClient(clientOptions)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := client.Connect(ctx)

	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(context.Background(), readpref.Primary())

	if err != nil {
		log.Fatal("Couldn't connect to the database", err)
	} else {
		log.Println("Connected to MondoDB Server")
	}

}

// function for handling the request from the client.
func handleRequest() {

	http.HandleFunc("/", homePage)
	http.HandleFunc("/api", homePage)

	http.HandleFunc("/api/addone", addOne)
	http.HandleFunc("/api/dump", dumpSubs)

	err := http.ListenAndServe(":5284", nil)
	if err != nil {
		log.Fatal("ListenAndServe", err)
	}
}

func dumpSubs(response http.ResponseWriter, request *http.Request) {

	if request.Method != http.MethodPost {
		http.Error(response, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var items []USERNOTICE
	if err := json.NewDecoder(request.Body).Decode(&items); err != nil {
		http.Error(response, "Cant Parse body", http.StatusNoContent)
	}

	if len(items) == 0 {
		http.Error(response, "No items in queue", http.StatusNoContent)
		return
	}
	log.Println("DUMP -> ", items)

	collection := client.Database(os.Getenv("database")).Collection(os.Getenv("collection"))
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var docs []interface{}
	for _, item := range items {
		docs = append(docs, bson.M{
			"ID":         item.ID,
			"Channel":    item.Channel,
			"Name":       item.Name,
			"SubMethod":  item.SubMethod,
			"SubAmount":  item.SubAmount,
			"GiftAmount": item.GiftAmount,
			"SubPlan":    item.SubPlan,
			"Created":    item.Created,
		})
	}

	// Bulk insert documents to MongoDB
	_, err := collection.InsertMany(ctx, docs)
	if err != nil {
		http.Error(response, "Failed to insert items", http.StatusInternalServerError)
		log.Println("InsertMany error:", err)
		return
	}

	response.WriteHeader(http.StatusOK)
	fmt.Fprintln(response, "Queue items written to MongoDB successfully")

}

func addOne(response http.ResponseWriter, request *http.Request) {

	if request.Method != http.MethodPost {
		http.Error(response, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	err := request.ParseForm()
	if err != nil {
		http.Error(response, "Failed to parse form", http.StatusBadRequest)
		return
	}

	q := request.Form
	if len(q) == 0 {
		http.Error(response, "No form items", http.StatusNoContent)
		return
	}

	log.Println("SINGLE ->", q)

	response.Write([]byte("Muxing data..."))

	collection := client.Database(os.Getenv("database")).Collection(os.Getenv("collection"))
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Create a document to insert
	document := map[string]interface{}{
		"ID":         q.Get("ID"),
		"Channel":    q.Get("Channel"),
		"Name":       q.Get("Name"),
		"SubMethod":  q.Get("SubMethod"),
		"SubAmount":  q.Get("SubAmount"),
		"GiftAmount": q.Get("giftAmount"),
		"SubPlan":    q.Get("subPlan"),
		"Created":    q.Get("Created"),
	}

	_, err = collection.InsertOne(ctx, document)
	if err != nil {
		http.Error(response, "Failed to insert data into MongoDB", http.StatusInternalServerError)
		return
	}

	response.Write([]byte("Sent data to MongoDB...."))

}

func homePage(w http.ResponseWriter, r *http.Request) {
	log.Println("home ok")
}
