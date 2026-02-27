package main

import (
	"fmt"
	"net/http"
)

func main(){
	fmt.Println("starting server on :8080")

	http.HandleFunc("/", func (w http.ResponseWriter, r *http.Request)  {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		panic(err)
	}
}