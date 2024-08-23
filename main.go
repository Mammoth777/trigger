package main

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)


func main() {
	fmt.Println("Server started on port 8080")
	http.HandleFunc("/", healthCheck)
	http.Handle("/execute-local-shell", http.HandlerFunc(executeLocalShell))

	err := http.ListenAndServe(":52323", nil)
	if err != nil {
		fmt.Println("Error starting server: ", err)
		return
	}
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("server is running"))
}

func executeLocalShell(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	runPath := queryParams.Get("run")
	file := filepath.Join(".", runPath)
	fmt.Println("Executing file: ", file)
	cmd := exec.Command(file)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		w.Write([]byte("Error executing command"))
		fmt.Println("Error executing command: ", err)
	}
	w.Write([]byte("Command executed"))
}