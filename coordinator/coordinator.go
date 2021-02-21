package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
)

func main() {

	var api = newAPI()
	must(rpc.Register(api))
	rpc.HandleHTTP()

	lis, err := net.Listen("tcp", ":4040")
	must(err)

	fmt.Println("Coordinator started on 4040")
	must(http.Serve(lis, nil))
}

// type worker struct {
// 	addr string
// 	status string
// }

//API server struct
type API struct {
	worker map[string]string //map [addr] -> status
}

func newAPI() *API {
	return &API{
		worker: make(map[string]string),
	}
}

//RegisterWorker : rpc for worker to register
func (a *API) RegisterWorker(addr string, reply *bool) error {
	a.worker[addr] = "idle" //idle, working, complete
	fmt.Printf("Wokers : %+v\n", a.worker)
	*reply = true
	return nil
}

func must(err error) {
	if err != nil {
		log.Fatal("Err: ", err)
		os.Exit(1)
	}
}
