package main

import (
	"bufio"
	"fmt"
	"Distributed-Grep/grepRPC"
	"log"
	"net/rpc"
	"os"
	"strings"
	"sync"
)

// serverlist format
// 192.168.0.1:80 ./machine.1.log
//

type SeverInfo struct {
	addr     string
	filepath string
}

type Args struct {
	pattern  string
	filepath string
}

// usage:  ./progam serverlistpath pattern file
func main() {
	//read command argument
	if len(os.Args) != 3 {
		log.Fatal(">Invalid input: ./grepClient <serverlistpath> <pattern> <filepath>")
	}
	// pattern filepath
	grepArg := Args{os.Args[2], os.Args[3]}
	serverListPath := os.Args[1]
	//fetch server list
	serverList, err := readServer(serverListPath)
	if err != nil {
		log.Fatal(">Can not read server list: ", err)
	}

	//connect each server sync.WaitGroup
	var wg sync.WaitGroup

	for i := range serverList {
		wg.Add(1)
		go distributedGrep(serverList[i], grepArg, &wg)
	}
	wg.Wait()
	fmt.Println(">Done fetching from all machines")

}

func readServer(path string) ([]SeverInfo, error) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(">Invalid Server List path")
	}

	defer file.Close()

	var Info []string // information of server addr, port filepath

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		Info = append(Info, scanner.Text())
	}

	var servers []SeverInfo

	for i := range Info {
		s := strings.Fields(Info[i])
		server := SeverInfo{s[0], s[1]}
		servers = append(servers, server)
	}

	return servers, scanner.Err()
}

//A client wishing to use the service establishes a connection
//and then invokes NewClient on the connection.

func distributedGrep(server SeverInfo, grepArg Args, wg *sync.WaitGroup) {
	defer wg.Done()
	//log.Println(">Server address:", server.addr)
	client, err := rpc.DialHTTP("tcp", server.addr)
	if err != nil {
		log.Fatal(">dialing error: ", err)
	}

	args2 := &grepRPC.GrepArgs{grepArg.pattern, grepArg.filepath}
	var reply string
	err = client.Call("GrepRes.GetGrep", args2, &reply)
	if err != nil {
		log.Fatal("grep error:", err)
	}
	fmt.Printf(">Results from %s: %s*%s=%s", server.addr, args2.Pattern, args2.File, reply)
}
