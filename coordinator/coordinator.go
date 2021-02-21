package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"sync"
	"time"
)

func main() {

	var api = newAPI()
	must(rpc.Register(api))
	rpc.HandleHTTP()
	go pollWorker(api)
	lis, err := net.Listen("tcp", ":4040")
	must(err)

	fmt.Println("Coordinator started on 4040")
	must(http.Serve(lis, nil))

	//TODO: where to start worker check
}

//API server struct
type API struct {
	done   bool
	worker map[string]string //map [addr] -> status
	Mut    *sync.Mutex
}

func newAPI() *API {
	return &API{
		Mut:    &sync.Mutex{},
		done:   false,
		worker: make(map[string]string),
	}
}

//RegisterWorker : rpc for worker to register
func (a *API) RegisterWorker(addr string, reply *string) error { //reply *ResgisterResponse

	a.Mut.Lock()
	a.worker[addr] = "idle" //idle, working, complete
	a.Mut.Unlock()

	fmt.Printf("Wokers : %+v\n", a.worker)
	*reply = "./plugintest.so"

	return nil
}

func pollWorker(a *API) {
	for {
		for k, _ := range a.worker {
			//poll rpc of worker
			_, err := rpc.DialHTTP("tcp", k) //pollClient //TODO: is commented part really needed ?
			// if err != nil {
			// 	delete(a.worker, k)
			// }
			// var res bool
			// err = pollClient.Call("", true, &res)
			if err != nil {
				fmt.Println("Error during poll : ", err)
				a.Mut.Lock()
				delete(a.worker, k)
				a.Mut.Unlock()
			}

		}
		fmt.Printf("Workers [Poll] : %+v\n", a.worker)
		time.Sleep(10 * time.Second)
	}
}

func must(err error) {
	if err != nil {
		log.Fatal("Err: ", err)
		os.Exit(1)
	}
}
