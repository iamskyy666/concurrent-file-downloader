
# Concurrent File Downloader in Golang 📂

A simple yet powerful Golang project demonstrating:

* Sequential file downloading
* Concurrent file downloading with Goroutines
* Channels for communication
* WaitGroups for synchronization
* Concurrency limiting using buffered channels
* File handling
* HTTP requests
* Error handling
* Download performance tracking

This project is a great practical introduction to Go concurrency and systems programming concepts.

---

# Features

## Single File Download

Download a single file from a URL and save it locally.

---

## Sequential Downloader

Download multiple files one-by-one.

Useful for understanding:

* blocking execution
* sequential processing
* baseline performance

---

## Concurrent Downloader

Download multiple files concurrently using Goroutines.

Features:

* Concurrent downloads
* Concurrency limiting
* WaitGroup synchronization
* Result aggregation
* Error collection
* Download statistics

---

# Concepts Covered

This project demonstrates many important Golang concepts:

| Concept           | Usage                       |
| ----------------- | --------------------------- |
| Goroutines        | Concurrent downloads        |
| Channels          | Communicating results       |
| Buffered Channels | Concurrency limiter         |
| WaitGroups        | Goroutine synchronization   |
| HTTP Requests     | Downloading files           |
| File Handling     | Saving downloaded content   |
| io.Copy           | Streaming file content      |
| Error Handling    | Graceful failure management |
| Time Tracking     | Performance monitoring      |
| Structs           | Result aggregation          |

---

# Project Structure

```text
.
├── main.go
├── downloads/
└── README.md
```

---

# How It Works

## 1. Single File Downloader

```go
func DownLoadFile(url string, destDir string) error
```

This function:

1. Creates the output file
2. Sends an HTTP GET request
3. Streams response body into the file
4. Tracks download time
5. Returns errors if something fails

---

## 2. Sequential Downloader

```go
func SequentialDownloader(urls []string, destDir string) error
```

This downloads files one-by-one.

Flow:

```text
URL 1 -> Download Complete
URL 2 -> Download Complete
URL 3 -> Download Complete
```

Only one download happens at a time.

---

## 3. Concurrent Downloader

```go
func ConcurrentDownloader(urls []string, destDir string, maxConcurrent int) error
```

This is the core part of the project.

It launches multiple Goroutines to download files simultaneously.

---

# Concurrency Architecture

```text
                 ┌─────────────┐
                 │ URLs Slice  │
                 └──────┬──────┘
                        │
                        ▼
             ┌──────────────────┐
             │ Goroutines Spawn │
             └────────┬─────────┘
                      │
      ┌───────────────┼───────────────┐
      ▼               ▼               ▼
┌──────────┐   ┌──────────┐   ┌──────────┐
│ Worker 1 │   │ Worker 2 │   │ Worker 3 │
└────┬─────┘   └────┬─────┘   └────┬─────┘
     │              │              │
     └──────────────┼──────────────┘
                    ▼
             ┌────────────┐
             │ results ch │
             └─────┬──────┘
                   ▼
            Main Result Loop
```

---

# Concurrency Limiter

One of the most important parts:

```go
limiter := make(chan struct{}, maxConcurrent)
```

This acts like a semaphore.

---

## Why Use It?

Without limiting:

* too many Goroutines may spawn
* memory usage may increase
* network congestion may occur
* remote servers may throttle requests

---

## How It Works

Before starting download:

```go
limiter <- struct{}{}
```

This acquires a slot.

After download completes:

```go
<-limiter
```

This releases the slot.

---

# WaitGroup Usage

```go
var wg sync.WaitGroup
```

WaitGroup tracks active Goroutines.

---

## Add Worker

```go
wg.Add(1)
```

---

## Mark Done

```go
defer wg.Done()
```

---

## Wait For All

```go
wg.Wait()
```

Ensures all downloads finish before program exits.

---

# Result Struct

```go
type Result struct {
    URL string
    FileName string
    Size int64
    Duration time.Duration
    Error error
}
```

This struct stores:

* downloaded file name
* file size
* download duration
* errors

---

# Result Channel

```go
results := make(chan Result)
```

Each worker sends its result into this channel.

Main Goroutine consumes results using:

```go
for result := range results {
}
```

---

# Download Flow

## Sequential

```text
Download File A
Wait...
Download File B
Wait...
Download File C
```

---

## Concurrent

```text
Download File A ─┐
Download File B ─┼── simultaneously
Download File C ─┘
```

Much faster for IO-bound operations.

---

# Why Concurrency Helps Here

File downloading is IO-bound.

Most of the time is spent waiting for:

* network responses
* remote servers
* disk IO

Concurrency allows Go to:

```text
work on other downloads while waiting
```

This improves throughput significantly.

---

# Example Output

```text
Downloading https://go.dev/blog/gopher/header.jpg
Downloading header.jpg took 624.7907ms ✅

Downloading https://go.dev/blog/gopher/wfmu.jpg
Downloading wfmu.jpg took 345.1078ms ✅

Downloading https://go.dev/blog/gopher/portrait.jpg
Downloading portrait.jpg took 446.9852ms ✅

🟢 Downloaded example1.png (32473 bytes) in 216.4452ms
🟢 Downloaded example2.png (104842 bytes) in 1.7387201s

✅ All downloads completed in 1.7405035s, Total: 137315 bytes
```

---

# Installation

## Clone Repository

```bash
git clone <your-repo-url>
cd concurrent-file-downloader
```

---

# Run Project

```bash
go run main.go
```

---

# Requirements

* Go 1.20+

---
