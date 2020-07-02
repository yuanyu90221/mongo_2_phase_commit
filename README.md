# mongo_2_phase_commit

## introduction

This is the topic about how the mongo 2 phase commit

cause problem while concurrent problem happen

## solution with goroutine channel

想法：

建立一個 in channel接收 request
建立一個out channel 印出結果

```golang
var in chan string
var out chan result
type result struct {
    Account string
    Result float64
}

func pay(w http.ResponseWriter, r *http.Request) {
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		in <- account // send account to goroutine handle store logic
		for {
			select {
			case data := <-out:
				fmt.Printf("%+v\n", data)
				wg.Done()
				return
			}
		}
	}(&wg)
	wg.Wait()
	io.WriteString(w, "ok")
}
func main() {
    in = make(chan string)
    out = make(chan result)
    // logic to access request
    go func(in *chan string) {
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

				out <- result{
					Account: account,
					Result:  entry.Amount,
				}
			}
		}
    }(&in)
    // rest logic
}
```

## solution with multpile queue
想法：

建立多個 in channel接收 request
建立多個out channel 透過 select 印出結果

```golang
var in []chan string
var out []chan result
type result struct {
    Account string
    Result float64
}

var maxUser = 100
var maxThread = 10
func pay(w http.ResponseWriter, r *http.Request) {
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		// random choose channel
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
func main() {
	// initial space
	in = make([]chan string, maxThread)
	out = make([]chan result, maxThread)
	// initial channel
	for i := range in {
		in[i] = make(chan string)
		out[i] = make(chan result)
	}
	// inital accounts
	for i := 0; i < maxUser; i++ {
		account = "user" + strconv.Itoa(i+1)
		user := currency{Account: account, Amount: 1000.00, Code: "USD"}
		if err := globalDB.C("bank").Insert(&user); err != nil {
			panic("insert error")
		}
	}

	// setup multiple select
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
}
```
