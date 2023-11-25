package main

import (
	"github.com/fayca121/cardchecker/src/controllers"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/check", controllers.CheckCardNumber)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}

}
