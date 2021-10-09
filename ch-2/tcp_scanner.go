package main

import (
	"fmt"
	"io"
	"sort"
	"net"
	"log"
	"os"
)


func main() {
	target := "127.0.0.1"
	maxPort := 1024
	numWorkers := 100
	scanTCP(target, maxPort, numWorkers)
}

////
// TCP POXY
////

// runs the proxy server
func startProxy() {
	listener, err := net.Listen("tcp", "0.0.0.0:20080")
	if err != nil {
		log.Fatalln("Unable to bind to port")
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalln("Unable to accept connection")
		}
		go handle(conn)
	}
}

// forwards all traffic to port 80 to google.com
func handle(src net.Conn) {
	dst, err := net.Dial("tcp", "google.com:80")
	if err != nil {
		log.Fatalln("Unable to connect")
	}
	defer dst.Close()

	go func() {
		if _, err = io.Copy(dst, src); err != nil {
			log.Fatalln(err)
		}
	}()

	if _, err = io.Copy(src, dst); err != nil {
		log.Fatalln(err)
	}

}

////
// ECHO SERVER
////

func serve() {
	port := 20080
	address := fmt.Sprintf("0.0.0.0:%d", port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalln("Unable to bind to port")
	}
	log.Printf("Listening on port %d\n", port)

	for {
		conn, err := listener.Accept()
		log.Println("Received connection")
		if err != nil {
			log.Fatalln("Unable to accept connection")
		}
		go echo(conn)
	}
}

func echo(conn net.Conn) {
	defer conn.Close()

	b := make([]byte, 512)

	for {
		s, err := conn.Read(b[0:])
		log.Printf("Read %d bytes", s)

		if err == io.EOF {
			log.Println("Client disconnected")
		}

		if err != nil {
			log.Fatalln("Unexpected error")
		}
		log.Printf("Received %d bytes: %s", s, string(b))

		log.Printf("Write data")
		_, err = conn.Write(b[0:s])
		if err != nil {
			log.Fatalln("Unable to write data")
		}
	}
}

////
// TCP PROXY
////

type FooReader struct {}
type FooWriter struct {}

func (fooReader *FooReader) Read(b []byte) (int, error) {
	fmt.Print("in > ")
	return os.Stdin.Read(b)
}

func (fooWriter *FooWriter) Write(b []byte) (int, error) {
	fmt.Print("out > ")
	return os.Stdin.Write(b)
}

func writeIO() {
	var (
		reader FooReader
		writer FooWriter
	)

	input := make([]byte, 4096)
	l, err := reader.Read(input)
	if err != nil {
		log.Fatalln("Unable to read data.")
	}
	fmt.Printf("Read %d bytes from stdin\n", l)

	l, err = writer.Write(input)
	if err != nil {
		log.Fatalln("Unable to write data.")
	}
	fmt.Printf("Wrote %d bytes to stdout\n", l)
}


////
// TCP SCANNER
///
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


