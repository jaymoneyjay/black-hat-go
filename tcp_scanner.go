package main

import (
	"fmt"
	"sort"
	"net"
)

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

func main() {
	TAR := "127.0.0.1"
	MAXPORT := 1024
	ports := make(chan int, 100)
	results := make(chan int)
	var openports []int

	for i:=0; i<=cap(ports); i++ {
		go worker(ports, results, TAR)
	}

	go func() {
		for i:=1; i<=MAXPORT; i++ {
			ports <- i
		}
	}()

	for i:=0; i<MAXPORT; i++ {
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


