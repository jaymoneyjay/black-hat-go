package main

import (
	"fmt"
	"sort"
	"net"
)

func main() {
	target := "127.0.0.1"
	maxPort := 1024
	numWorkers := 100
	scanTCP(target, maxPort, numWorkers)
}

func worker(ports, results chan int, target string) {
	for p := range ports {
		address := fmt.Sprintf("%s:%d", target, p)
		conn, err := net.Dial("tcp", address)
		if err != nil {
			results <- 0
			continue
		}
		conn.Close()
		results <- p
	}
}

func scanTCP(target string, maxPort, numWorkers int) {
	ports := make(chan int, numWorkers)
	results := make(chan int)
	var openports []int

	for i:=0; i<=cap(ports); i++ {
		go worker(ports, results, target)
	}

	go func() {
		for i:=1; i<=maxPort; i++ {
			ports <- i
		}
	}()

	for i:=0; i<maxPort; i++ {
		port := <-results
		if port != 0 {
			openports = append(openports, port)
		}
	}

	close(ports)
	close(results)

	sort.Ints(openports)
	for _, port := range openports {
		fmt.Printf("%d open\n", port)
	}
}


