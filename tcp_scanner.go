package main

import (
	"fmt"
	"net"
)

func scanPort(port int) {
	address := fmt.Sprintf("scanme.nmap.org:%d", port)
	con, err := net.Dial("tcp", address)
	if err != nil {
		return
	}
	statement := fmt.Sprintf("Port %d open.", port)
	fmt.Println(statement)
	con.Close()
}

func main() {
	for i:=1; i<=1024; i++ {
		scanPort(i)
	}
}


