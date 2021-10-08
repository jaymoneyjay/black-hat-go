package main

import (
	"fmt"
	"log"
	"os"
	"github.com/jaymoneyjay/black-hat-go/ch-3/msf_client/rpc"
	"github.com/joho/godotenv"
	)

func goDotEnvVariable(key string) string {

	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}

func main() {
	host := goDotEnvVariable("MSFHOST")
	pass := goDotEnvVariable("MSFPASS")
	user := "msf"

	if host == "" || pass == "" {
		log.Fatalln("Missing environment variable MSFHOST or MSFPASS.")
	}

	msf, err := rpc.New(host, user, pass)
	if err != nil {
		log.Fatalln(err)
	}

	defer msf.Logout()

	sessions, err := msf.SessionList()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("Sessions:")
	for _, s := range sessions {
		fmt.Printf("%5d: %s\n", s.ID, s.Info)
	}
}