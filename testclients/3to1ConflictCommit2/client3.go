/*

A trivial client to illustrate how the kvservice library can be used
from an application in assignment 6 for UBC CS 416 2016 W2.

Usage:
go run client.go
*/

package main

// Expects kvservice.go to be in the ./kvservice/ dir, relative to
// this client.go file
import (
	"./../../kvservice"
	"time"
	"fmt"
	"bufio"
	"io"
	"os"
	"strings"
	"log"
)

func main() {
	args := os.Args[1:]

	// nodes = append(nodes, "127.0.0.1:9091")
	// Missing command line args.
	if len(args) != 1 {
		fmt.Println("Usage: go run x.go [nodesFile]")
		return
	}

	var nodes []string

	pF, err := os.Open(args[0])
	checkError("nodesFile argument cannot be opened", err, true)
	bufIo := bufio.NewReader(pF)
	for {
		nodeIP, err := bufIo.ReadString('\n')
		nodeIP = strings.TrimSpace(nodeIP)
		fmt.Println("nodeIP: ", nodeIP)
		if err == io.EOF {
			break
		}
		nodes = append(nodes, nodeIP)		
		fmt.Println("nodes: ", nodes)
	}

	c := kvservice.NewConnection(nodes)
	fmt.Printf("NewConnection returned: %v\n", c)
	
	t, err := c.NewTX()
	fmt.Printf("NewTX returned: %v, %v\n", t, err)

	success, err := t.Put("k2", "t3 world")
	fmt.Printf("Put returned: %v, %v\n", success, err)

	success, v, err := t.Get("k2")
	fmt.Printf("Get returned: %v, %v, %v\n", success, v, err)


	success, err = t.Put("k3", "t3 hello world")
	fmt.Printf("Put returned: %v, %v\n", success, err)

	success, v, err = t.Get("k3")
	fmt.Printf("Get returned: %v, %v, %v\n", success, v, err)
	

	// commit order tx1, tx2, tx3
	fmt.Println("delay for 2 sec")
	DelayForSec(2)
	
	success, txID, err := t.Commit(2)
	fmt.Printf("Commit returned: %v, %v, %v\n", success, txID, err)

	// check getChildren
	fmt.Println("delay for 2 sec for getChildren")
	DelayForSec(2)

	var getChildren func(node string, parentHash string);
	getChildren = func(node string, parentHash string){
		children := c.GetChildren(node, parentHash)
		fmt.Printf("ParentHash: %v \n", parentHash)
		fmt.Printf("Children: %v \n", children)
		fmt.Printf("=============================================\n\n")
		if len(children) > 0 {
			for c := range children{
				getChildren(node, children[c])
			}
		}
	};

	for _, nodeIP := range nodes {
		fmt.Printf("GetChildren for node [%v] \n", nodeIP)
		getChildren(nodeIP, "")
		fmt.Printf("\n")
	}


	c.Close()
}


func DelayForSec(i time.Duration) {
	time.Sleep(time.Millisecond * 1000 * i)
}

func checkError(msg string, err error, exit bool) {
	if err != nil {
		log.Println(msg, err)
		if exit {
			os.Exit(-1)
		}
	}
}