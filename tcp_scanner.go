package main

import (
	"fmt"
	"net"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	for i:=0; i<=65535;i++ {
		wg.Add(1)
		go func(port int) {
			defer wg.Done()
			address := fmt.Sprintf("127.0.0.1:%d", port)
			con, err := net.Dial("tcp", address)
			if err != nil {
				return
			}
			statement := fmt.Sprintf("Port %d open.", port)
			fmt.Println(statement)
			con.Close()

		}(i)
	}
	wg.Wait()
}


