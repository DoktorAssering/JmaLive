package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Тест прошёл успешно!"))
	})

	log.Fatal(http.ListenAndServeTLS(":443", "cert.pem", "key.pem", nil))

	log.Println("Сервер успешно запустился")
}
