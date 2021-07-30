package main

import (
	"file-server/server"
	"fmt"
	"net/http"
)

func main() {
	go http.HandleFunc("/upload", server.Upload)
	go http.HandleFunc("/download/",server.Download)
	http.Handle("/static/", http.StripPrefix("/static/", server.ServeHandle))
	err := http.ListenAndServe(":8888", nil)
	if err != nil {
		fmt.Println(err)
	}
}

//curl -i --range 2-10 localhost:8888/download/test1
//curl -i localhost:8888/download/test1
//curl -F "filename=@test4" localhost:8888/upload
