package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// 📂 concurrent-file-downloader

// Single File
func DownLoadFile(url string, destDir string)error{
	fileName:=filepath.Base(url)
	fimePath:=filepath.Join(destDir,fileName)
	outPut,err:=os.Create(fimePath) // ex - ./downloads/sample.txt
	if err != nil {
		log.Fatal("ERROR:",err)
	}
	defer outPut.Close()
	fmt.Println("Downloading",url)

	// track the time
	start:=time.Now()

	resp,err:=http.Get(url)
	if err != nil {
		log.Fatal("ERROR:",err)
	}

	defer resp.Body.Close()

	if resp.StatusCode!=http.StatusOK{
		return fmt.Errorf("bad_status: %s",resp.Status)
	}

	// copy content of the body into our file
	_,err = io.Copy(outPut,resp.Body)
	if err != nil {
		log.Fatal("ERROR:",err)
	}

	fmt.Printf("Downloading %s took %s ✅\n",fileName,time.Since(start))

	return nil
}

// Multiple Files
func SequentialDownloader(urls []string,destDir string)error{
	if err:=os.MkdirAll(destDir,0755);err!=nil{
		log.Fatal("ERROR:",err)
	}

	// track the time
	start:=time.Now()

	for _, url := range urls {
		if err:=DownLoadFile(url,destDir);err!=nil{
			log.Println("ERROR doanloading...",url,err)
			continue
		}
	}
	fmt.Printf("Downloading %s took %s ✅\n",urls,time.Since(start))

	return nil
}

func main() {
	url:="https://go.dev/blog/gopher/header.jpg"
	err:=DownLoadFile(url,"./downloads")
	if err != nil {
		log.Println("ERROR:",err)
	}

	urls:=[]string{"https://go.dev/blog/gopher/wfmu.jpg","https://go.dev/blog/gopher/portrait.jpg"}
	err=SequentialDownloader(urls,"./downloads")
	if err != nil {
		log.Println("ERROR:",err)
	}

	log.Println("\nDone ☑️")
}