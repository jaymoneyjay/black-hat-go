package main

import (
	"github.com/jaymoneyjay/black-hat-go/ch-3/shodan_client/shodan"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
)

// use godot package to load/read the .env file and
// return the value of the key
func goDotEnvVariable(key string) string {

	// load .env file
	err := godotenv.Load()

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}

func main() {
	apiKey := goDotEnvVariable("SHODAN_API_KEY")
	client := shodan.New(apiKey)

	// call api info
	retAPI, err := client.APIInfo()
	if err != nil {
		log.Panicln(err)
	}
	fmt.Println(retAPI)

	// query host
	retHost, err := client.HostSearch(os.Args[1])
	if err != nil {
		log.Panicln(err)
	}
	for _, host := range retHost.Matches {
		fmt.Printf("%18s%8d\n", host.IPString, host.Port)
	}
}

