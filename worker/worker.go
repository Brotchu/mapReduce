package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"plugin"
	"strings"
)

// type API int

// func (a *API) Ack()

func main() {
	//os args 1- local port to use ; 2- addr of coordinator
	if len(os.Args) < 3 {
		log.Fatal("Err: Specify (1)localport to use and (2)address of coordinator")
		os.Exit(1)
	}
	localPort := os.Args[1]
	coordinatorAddr := os.Args[2]

	//connect to coordinator
	client, err := rpc.DialHTTP("tcp", coordinatorAddr)
	must(err)

	var reply string //get function name
	err = client.Call("API.RegisterWorker", getLocalAddr()+":"+localPort, &reply)
	must(err)
	fmt.Println("Worker registered")

	workerFunc := loadPlugin(reply)
	workerFunc()
	//start go routine to reply to coordinator checks
	//TODO:
	go listenForPoll(localPort)

	//Testing : waiting indefinitely
	<-make(chan int)
}

func getLocalAddr() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	must(err)
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	addrString := strings.Split(localAddr.String(), ":")
	return addrString[0]
}

func loadPlugin(pName string) func() {
	p, err := plugin.Open(pName)
	must(err)

	f, err := p.Lookup("TestDisplay")
	funcTest := f.(func())
	return funcTest
}

type WorkerAPI struct{}

func (w *WorkerAPI) Poll(check bool, reply *bool) error {
	*reply = true
	return nil
}

func listenForPoll(port string) {
	var wAPI = &WorkerAPI{}
	must(rpc.Register(wAPI))
	rpc.HandleHTTP()

	lis, err := net.Listen("tcp", ":"+port)
	must(err)

	must(http.Serve(lis, nil))

}

func must(err error) {
	if err != nil {
		log.Fatal("Err: ", err)
		os.Exit(1)
	}
}
