package main

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"sync"
	"time"

	_ "github.com/joho/godotenv/autoload"
	config "github.com/yuanyu90221/mongo_2_phase_commit/config"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var globalDB *mgo.Database
var account = "json"
var mu = &sync.Mutex{}

type currency struct {
	ID      bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Amount  float64       `bson:"amount"`
	Account string        `bson:"account"`
	Code    string        `bson:"code"`
}

func pay(w http.ResponseWriter, r *http.Request) {
	entry := currency{}
	mu.Lock()
	defer mu.Unlock()
	// step 1: get current amount
	err := globalDB.C("bank").Find(bson.M{"account": account}).One(&entry)

	if err != nil {
		panic(err)
	}

	wait := Random(1, 100)
	time.Sleep(time.Duration(wait) * time.Millisecond)

	// step 3: subtract current balance and update back to the db
	entry.Amount = entry.Amount + 50.000
	err = globalDB.C("bank").UpdateId(entry.ID, &entry)

	if err != nil {
		panic("update error")
	}

	fmt.Printf("%+v\n", entry)
	io.WriteString(w, "ok")
}
func Random(min, max int) int {
	rand.Seed(time.Now().UTC().UnixNano())
	return rand.Intn(max-min+1) + min
}
func main() {
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8000"
	}
	connectString := config.GetConnectString()
	session, _ := mgo.Dial(connectString)
	globalDB = session.DB("test")
	globalDB.C("bank").DropCollection()

	user := currency{Account: account, Amount: 1000.00, Code: "USD"}
	err := globalDB.C("bank").Insert(&user)
	if err != nil {
		panic("Insert error")
	}
	log.Println("Listen server on " + port + " port")
	http.HandleFunc("/", pay)

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
