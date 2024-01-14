package main

import (
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
)

func main() {

	tmplargs := []string{"generate", "--watch", "--proxy=http://localhost:8080", "--cmd=go run ."}
	templWatchCmd := exec.Command("templ", tmplargs...)
	go processCommand(templWatchCmd)

	// Common pattern for managing shutdowns of long-running processes in go.
	// SIGINT: Ctrl+C
	// SIGTERM: termination signal, typically sent by system-level tools
	// to request a graceful shutdown.
	// Since we're listening for the signal and blocking until it's received,
	// this allows us to decide what and how we want to cleanup before exiting.
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGINT, syscall.SIGTERM)
	<-sigc

}

func processCommand(cmd *exec.Cmd) {
	//	templOut, err := cmd.StdoutPipe()
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//
	//	if err := cmd.Start(); err != nil {
	//		log.Fatal(err)
	//	}
	//	tmplR := io.TeeReader(templOut, os.Stdout)
	//
	//	if _, err := io.ReadAll(tmplR); err != nil {
	//		log.Fatal(err)
	//	}
	//
	//	if err := cmd.Wait(); err != nil {
	//		log.Fatal(err)
	//	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		log.Printf("Error starting command: %v\n", err)
		println("exiting")
		return // instead of log.Fatal to avoid exiting the entire program
	}

	if err := cmd.Wait(); err != nil {
		log.Printf("Command finished with error: %v\n", err)
	}

}
