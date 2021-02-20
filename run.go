package main

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
)

func main() {
	http.HandleFunc("/", scriptHandler)

	// Determine port for HTTP service.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	// Start HTTP server.
	log.Printf("Listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

func fileExists(name string) bool {
	_, err := os.Stat(name)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func scriptFilename(url string) string {
	filename := url + ".sh"
	if fileExists(filename) {
		return filename
	}
	return "script.sh"
}

func scriptHandler(w http.ResponseWriter, r *http.Request) {
	filename := scriptFilename(r.URL.Path[1:])

	// From exec.CommandContext():
	// The provided context is used to kill the process (by calling
	// os.Process.Kill) if the context becomes done before the command
	// completes on its own.

	// So if script takes too long it is eventually killed.
	// So, using exec.Command() instead.
	cmd := exec.Command("/bin/bash", filename)

	// https://blog.kowalczyk.info/article/wOYk/advanced-command-execution-in-go-with-osexec.html
	var stdoutBuf, stderrBuf bytes.Buffer
	cmd.Stdout = io.MultiWriter(os.Stdout, &stdoutBuf)
	cmd.Stderr = io.MultiWriter(os.Stderr, &stderrBuf)

	err := cmd.Run()
	if err != nil {
		w.WriteHeader(500)
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}
	outStr, errStr := string(stdoutBuf.Bytes()), string(stderrBuf.Bytes())
	log.Printf("\nout:\n%s\nerr:\n%s\n", outStr, errStr)
	_, _ = w.Write([]byte("done"))
}
