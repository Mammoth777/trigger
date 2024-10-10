package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strconv"
	"syscall"
)

func recordPid(pid int) {
	file, err := os.OpenFile("server.pid", os.O_CREATE|os.O_WRONLY, 0644);
	if err != nil {
		log.Println("Error creating pid file:", err)
		return
	}
	defer file.Close()
	_, err = file.WriteString(fmt.Sprintf("%d", pid))
	if err != nil {
		log.Println("Error writing pid file:", err)
	}
	log.Printf("Pid file created: %d\n", pid)
}


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
	
	recordPid(pid)

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
		log.Println("Error reading pid file:", err)
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
		log.Println("Error getting current directory:", err)
		return false
	}
	return filepath.HasPrefix(target, limitDir)
}

func executeLocalShell(w http.ResponseWriter, r *http.Request) {
	dir, err := os.Getwd()
	if err != nil {
		log.Println("Error getting current directory:", err)
		w.Write([]byte("Error getting current directory\n"))
		return
	}
	queryParams := r.URL.Query()
	runPath := queryParams.Get("run")
	param := queryParams.Get("param")
	file := filepath.Join(dir, runPath)
	if runPath == "" {
		w.Write([]byte("Error: No file specified\n"))
		log.Println("Error: No file specified")
		return
	}
	if !inLimitDir(file) {
		w.Write([]byte("Error: File not in current directory\n"))
		log.Println("Error: File not in current directory")
		return
	}
	log.Printf("Executing file: %s\n", file)
	cmd := exec.Command(file, param)
	// 确保目录存在
	logDir := "logs"
	if err := os.MkdirAll(logDir, os.ModePerm); err != nil {
			log.Println("Error creating directory:", err)
			return
	}
	outputFile, err := os.OpenFile("logs/exec-output.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		w.Write([]byte("Error opening output file: " + err.Error() + "\n"))
		log.Println("Error opening output file: ", err)
		return
	}
	defer outputFile.Close()
	cmd.Stdout = outputFile
	cmd.Stderr = outputFile
	log.Printf("Executing: %s %s \n", file, param)
	err = cmd.Run()
	if err != nil {
		w.Write([]byte("Error executing command: " + err.Error() + "\n"))
		log.Println("Error executing command: ", err)
		return
	}
	w.Write([]byte("Command executed\n"))
}
