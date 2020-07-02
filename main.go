package main

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	_ "github.com/joho/godotenv/autoload"
	config "github.com/yuanyu90221/mongo_2_phase_commit/config"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var globalDB *mgo.Database
var account = "json"
var in []chan string
var out []chan result
var maxUser = 100
var maxThread = 10

type result struct {
	Account string
	Result  float64
}

type currency struct {
	ID      bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Amount  float64       `bson:"amount"`
	Account string        `bson:"account"`
	Code    string        `bson:"code"`
}

func pay(w http.ResponseWriter, r *http.Request) {
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		number := Random(1, maxUser)
		channelNumber := number % maxThread
		account := "user" + strconv.Itoa(number)

		in[channelNumber] <- account
		for {
			select {
			case data := <-out[channelNumber]:
				fmt.Printf("%+v\n", data)
				wg.Done()
				return
			}
		}
	}(&wg)
	wg.Wait()
	io.WriteString(w, "ok")
}
func Random(min, max int) int {
	rand.Seed(time.Now().UTC().UnixNano())
	return rand.Intn(max-min+1) + min
}
func main() {
	in = make([]chan string, maxThread)
	out = make([]chan result, maxThread)

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8000"
	}
	connectString := config.GetConnectString()
	session, _ := mgo.Dial(connectString)
	globalDB = session.DB("test")
	globalDB.C("bank").DropCollection()

	for i := range in {
		in[i] = make(chan string)
		out[i] = make(chan result)
	}

	for i := 0; i < maxUser; i++ {
		account = "user" + strconv.Itoa(i+1)
		user := currency{Account: account, Amount: 1000.00, Code: "USD"}
		if err := globalDB.C("bank").Insert(&user); err != nil {
			panic("insert error")
		}
	}
	for i := range in {
		go func(in *chan string, i int) {
			for {
				select {
				case account := <-*in:
					entry := currency{}
					// step 1: get current amount
					err := globalDB.C("bank").Find(bson.M{"account": account}).One(&entry)

					if err != nil {
						panic(err)
					}
					// step 3: subtract current balance and update back to the db
					entry.Amount = entry.Amount + 50.000
					err = globalDB.C("bank").UpdateId(entry.ID, &entry)

					if err != nil {
						panic("update error")
					}

					out[i] <- result{
						Account: account,
						Result:  entry.Amount,
					}
				}
			}
		}(&in[i], i)
	}

	log.Println("Listen server on " + port + " port")
	http.HandleFunc("/", pay)

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
