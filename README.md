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