package main

import (
    "log"
    "net/http"
)

func main() {
    fs := http.FileServer(http.Dir("."))

    http.Handle("/", fs)

    log.Println("Website running at http://localhost:5500")

    log.Fatal(http.ListenAndServe("127.0.0.1:5500", nil))
}