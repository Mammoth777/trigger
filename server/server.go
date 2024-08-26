package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
)

var stop chan os.Signal
var pid int

func StartServer() {
	server := &http.Server{
		Addr: ":52323",
	}
	http.HandleFunc("/", healthCheck)
	http.Handle("/execute-local-shell", http.HandlerFunc(executeLocalShell))
	
	go func() {
		fmt.Println("Server started on port 52323")
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			fmt.Println("Error starting server: ", err)
			return
		}
	}()
	// 创建一个通道来监听系统中断信号
	stop = make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
  pid = os.Getpid()
	fmt.Println("pid", pid)

	// 阻塞主线程，直到接收到系统中断信号
	<- stop
	fmt.Println("Shutting down server...")
	
	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		fmt.Println("Error shutting down server: ", err)
	}
	fmt.Println("Server stopped")
}

func StopServer() {
	fmt.Println("Server stopped")
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