package main

import (
	"fmt"
	"github.com/mlctrez/godom-example/example"
	"github.com/mlctrez/godom/app/ctx"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func main() {
	deadline := time.Now().Add(time.Second * 30)
	go func() { ctx.Run(example.New()) }()

	success := false
	for time.Now().Before(deadline) {
		if err := copyFiles(); err != nil {
			time.Sleep(1 * time.Second)
			continue
		}
		success = true
		break
	}
	if !success {
		log.Fatal("error copying files")
	}
}

func copyFiles() (err error) {
	files := []string{"index.html", "app.js"}
	for _, f := range files {
		if err = copyUrlToFile(f); err != nil {
			return err
		}
		fmt.Println("copied", f)
	}
	return nil
}

func copyUrlToFile(url string) (err error) {
	var outFile *os.File
	var out = url
	if url == "index.html" {
		url = ""
	}
	if outFile, err = os.Create(filepath.Join("build", out)); err != nil {
		return err
	}
	defer func() { _ = outFile.Close() }()
	var resp *http.Response
	if resp, err = http.Get(fmt.Sprintf("http://localhost:8080/%s", url)); err != nil {
		return err
	}
	_, err = io.Copy(outFile, resp.Body)
	return err
}
