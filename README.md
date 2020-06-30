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