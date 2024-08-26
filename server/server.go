package server

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)

func StartServer() {
	fmt.Println("Server started on port 52323")
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

func inLimitDir(target string) bool {
	limitDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current directory:", err)
		return false
	}
	return filepath.HasPrefix(target, limitDir)
}

func executeLocalShell(w http.ResponseWriter, r *http.Request) {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current directory:", err)
		w.Write([]byte("Error getting current directory"))
		return
	}
	queryParams := r.URL.Query()
	runPath := queryParams.Get("run")
	param := queryParams.Get("param")
	file := filepath.Join(dir, runPath)
	if !inLimitDir(file) {
		w.Write([]byte("Error: File not in current directory"))
		fmt.Println("Error: File not in current directory")
		return
	}
	fmt.Println("Executing file: ", file)
	cmd := exec.Command(file, param)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		w.Write([]byte("Error executing command" + err.Error()))
		fmt.Println("Error executing command: ", err)
	}
	w.Write([]byte("Command executed"))
}