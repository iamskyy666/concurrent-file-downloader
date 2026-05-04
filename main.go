package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// 📂 concurrent-file-downloader

// Single File
func DownLoadFile(url string, destDir string)error{
	fileName:=filepath.Base(url)
	filePath:=filepath.Join(destDir,fileName)
	outPut,err:=os.Create(filePath) // ex - ./downloads/sample.txt
	if err != nil {
		log.Fatal("ERROR:",err)
	}
	defer outPut.Close()
	fmt.Println("Downloading",url)

	// track the time
	start:=time.Now()

	resp,err:=http.Get(url)
	if err != nil {
		_= os.Remove(filePath)
		log.Fatal("ERROR:",err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK{
		_= os.Remove(filePath)
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

// Now, with concurrency
type Result struct {
	URL string
	FileName string
	Size int64
	Duration time.Duration
	Error error
}


func ConcurrentDownloader(urls []string, destDir string, maxConcurrent int)error{
	if err:=os.MkdirAll(destDir,0755);err!=nil{
		log.Fatal("ERROR:",err)
	}

	results:=make(chan Result)

	var wg sync.WaitGroup

	limiter:= make(chan struct{}, maxConcurrent)

	for _, url := range urls {
		wg.Add(1)
		go func(url string){
			defer wg.Done()

			limiter<-struct{}{}
			defer func(){<-limiter}()

	// track the time
	start:=time.Now()

			fileName:=filepath.Base(url)
	filePath:=filepath.Join(destDir,fileName)

	outPut,err:=os.Create(filePath) // ex - ./downloads/sample.txt
	if err != nil {
		results<-Result{URL:url, Error: err}
		log.Fatal("ERROR:",err)
	}
		defer outPut.Close()

	resp,err:=http.Get(url)
	if err != nil {
		results<-Result{URL:url, Error: err}
		log.Fatal("ERROR:",err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK{
		results<-Result{URL:url, Error: fmt.Errorf("bad_status: %s",resp.Status)}
		return 
	}

	size,err:=io.Copy(outPut, resp.Body)
	if err != nil {
		results<-Result{URL:url, Error: err}
		log.Fatal("ERROR:",err)
	}

	timeSince:=time.Since(start)
	results<-Result{URL:url, FileName: fileName, Size: size, Duration: timeSince, Error: nil}

		}(url)


	}

	go func() {
		wg.Wait()
		close(results)
	}()

	var totalSize int64
	var errors []error
	start:=time.Now()

	for result:=range results{
		if result.Error!=nil{
			fmt.Printf("Error downloading %s: %s\n",result.URL,result.Error.Error())
			errors = append(errors, result.Error)
		}else{
			totalSize+=result.Size
			fmt.Printf("🟢 Downloaded %s (%d bytes) in %s\n",result.FileName,result.Size,result.Duration)
		}
	}
	startedSince:=time.Since(start)
	fmt.Printf("✅ All downloads completed in %s, Total: %d bytes\n",startedSince,totalSize)
	
	if len(errors)>0{
		return fmt.Errorf("errors in downloading: %+v",errors)
	}

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

	concurr_urls:=[]string{"https://media2.dev.to/dynamic/image/width=1000,height=420,fit=cover,gravity=auto,format=auto/https%3A%2F%2Fdev-to-uploads.s3.amazonaws.com%2Fuploads%2Farticles%2Felrr8p507bxbsy1rdni5.png","https://canopas-blogs.s3.ap-south-1.amazonaws.com/Guide_to_Concurrency_in_Golang_Key_Terms_and_Examples_9d87a1ec44.png"}
	err=ConcurrentDownloader(concurr_urls,"./downloads",3)
	if err != nil {
		log.Println("ERROR:",err)
	}

	log.Println("\nDone ☑️")
}

//Output:
// $ go run main.go
// Downloading https://go.dev/blog/gopher/header.jpg
// Downloading header.jpg took 624.7907ms ✅
// Downloading https://go.dev/blog/gopher/wfmu.jpg
// Downloading wfmu.jpg took 345.1078ms ✅
// Downloading https://go.dev/blog/gopher/portrait.jpg
// Downloading portrait.jpg took 446.9852ms ✅
// Downloading [https://go.dev/blog/gopher/wfmu.jpg https://go.dev/blog/gopher/portrait.jpg] took 795.426ms ✅
// 🟢 Downloaded Guide_to_Concurrency_in_Golang_Key_Terms_and_Examples_9d87a1ec44.png (32473 bytes) in 216.4452ms
// 🟢 Downloaded https%3A%2F%2Fdev-to-uploads.s3.amazonaws.com%2Fuploads%2Farticles%2Felrr8p507bxbsy1rdni5.png (104842 bytes) in 1.7387201s
// ✅ All downloads completed in 1.7405035s, Total: 137315 bytes
// 2026/05/04 23:25:08 
// Done ☑️
