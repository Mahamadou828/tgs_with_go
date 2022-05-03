package main

import (
	"fmt"
	"html"
	"net/http"
)

//The build represent the environment that the current program is running
//for this specific programm we have 3 stages: dev, staging, prod
var build = "dev"

func main() {
	//Simply print a start message with the build and block the programm on a syscall
	fmt.Println("Starting the program", build)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	})

	http.ListenAndServe(":3000", nil)
}
