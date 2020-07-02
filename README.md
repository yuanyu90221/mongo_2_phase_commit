# mongo_2_phase_commit

## introduction

This is the topic about how the mongo 2 phase commit

cause problem while concurrent problem happen

## solution with mutex

```golang
var mu = &sync.Mutex{}
// add mutex Logic in race condition
// this is pay func

func pay(w http.ResponseWriter, r *http.Request) {
    entry := currency{}
    mu.Lock()
    defer mu.Unlock() // this will call before return
	// step 1: get current amount
	err := globalDB.C("bank").Find(bson.M{"account": account}).One(&entry)
//
}
```
## 多個application 同時執行
新增一個version 欄位

每次寫入把version+1

```golang
type currency struct {
	ID      bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Amount  float64       `bson:"amount"`
	Account string        `bson:"account"`
	Code    string        `bson:"code"`
	Version int           `bson:"version"`
}
// 每次寫入把version加一
func pay(w http.ResponseWriter, r *http.Request) {
	entry := currency{}
LOOP:
	// step 1: get current amount
	err := globalDB.C("bank").Find(bson.M{"account": account}).One(&entry)

	if err != nil {
		panic(err)
	}

	wait := Random(1, 100)
	time.Sleep(time.Duration(wait) * time.Millisecond)

	// step 3: subtract current balance and update back to the db
	entry.Amount = entry.Amount + 50.000
	err = globalDB.C("bank").Update(bson.M{
		"version": entry.Version,
		"_id":     entry.ID,
	}, bson.M{"$set": map[string]interface{}{
		"amount":  entry.Amount,
		"version": (entry.Version + 1),
	}})

	if err != nil {
		goto LOOP
	}

	fmt.Printf("%+v\n", entry)
	io.WriteString(w, "ok")
}
```