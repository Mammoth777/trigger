package server

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strconv"
	"syscall"
)


func StartServer() {
	server := &http.Server{
		Addr: ":52323",
	}
	http.HandleFunc("/", healthCheck)
	http.Handle("/execute-local-shell", http.HandlerFunc(executeLocalShell))

	pid := os.Getpid()
	go func ()  {
		fmt.Printf("(%d)Server started on port 52323\n", pid)
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			fmt.Println("Error starting server: ", err)
			return
		}
	}()
	err := os.WriteFile("server.pid", []byte(fmt.Sprintf("%d", pid)), 0644)
	if err != nil {
		fmt.Println("Error writing pid file:", err)
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	StopServer()
}

func StartServerDeamon() {
	cmd := exec.Command(os.Args[0], "serve")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Start()
	if err != nil {
		fmt.Println("Error starting server as deamon:", err)
		return
	}
}

func StopServer() {
	pidFile, err := os.ReadFile("server.pid")
	if err != nil {
		fmt.Println("Error reading pid file:", err)
		return
	}
	pid, err := strconv.Atoi(string(pidFile))
	if err != nil {
		fmt.Println("Error converting pid to int:", err)
		return
	}
	pcs, err := os.FindProcess(pid)
	if err != nil {
		fmt.Println("Error finding process:", err)
		return
	}
	err = pcs.Signal(syscall.SIGTERM)
	if err != nil {
		fmt.Println("Error stopping process:", err)
		return
	}
	err = os.Remove("server.pid")
	if err != nil {
		fmt.Println("Error removing pid file:", err)
		return
	}
	fmt.Println("Process stopped")
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
