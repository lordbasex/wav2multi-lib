# Example 3: Advanced Usage with io.Reader and io.Writer

This example demonstrates advanced usage of `wav2multi-lib` using Go's `io.Reader` and `io.Writer` interfaces. This is essential for streaming, HTTP APIs, in-memory processing, and network communication.

## 📋 Overview

This example shows:
- ✅ How to use `TranscodeFromReader()` to read from any `io.Reader`
- ✅ How to use `TranscodeToWriter()` to write to any `io.Writer`
- ✅ Real-world use cases for each method
- ✅ Streaming and in-memory processing patterns
- ✅ HTTP API integration examples

**Perfect for:** Web APIs, microservices, streaming applications, serverless functions

## 🎯 Why Use io.Reader/Writer?

### Traditional Approach (Files Only):
```go
// ❌ Limited to files on disk
config := wav2multi.TranscoderConfig{
    InputPath:  "input.wav",      // Must be a file
    OutputPath: "output.ulaw",     // Must write to disk
    Format:     wav2multi.FormatULaw,
}
```

### Advanced Approach (Any I/O):
```go
// ✅ Can be: HTTP request, network, memory, pipe, etc.
result, err := transcoder.TranscodeFromReader(
    reader,              // Any io.Reader
    "output.ulaw",
    wav2multi.FormatULaw,
)

// ✅ Can be: HTTP response, network, memory, pipe, etc.
result, err := transcoder.TranscodeToWriter(
    "input.wav",
    writer,              // Any io.Writer
    wav2multi.FormatSLIN,
)
```

## 🚀 Running the Example

### Without CGO (μ-law, A-law, SLIN):
```bash
cd example/example3
CGO_ENABLED=0 go run main.go
```

### With CGO (including G.729):
```bash
export CGO_ENABLED=1
export CGO_CFLAGS="-I/usr/local/include"
export CGO_LDFLAGS="-L/usr/local/lib -lbcg729"
cd example/example3
go run main.go
```

## 📝 Expected Output

```
🎵 wav2multi-lib - Example 3: Advanced Usage
=============================================

📖 Example 1: TranscodeFromReader
   Reading WAV from io.Reader and saving to file
   (Useful for HTTP uploads, memory buffers, etc.)
   ✅ Success!
      Output: output-from-reader.alaw (16104 bytes)
      Format: A-law @ 64.0 kbps
      Processing time: 12 ms

📝 Example 2: TranscodeToWriter
   Reading WAV from file and writing to io.Writer
   (Useful for HTTP responses, network streams, etc.)
   ✅ Success!
      Written to buffer: 32208 bytes
      Format: SLIN @ 128.0 kbps
      Processing time: 8 ms
      Saved buffer to: output-to-writer.slin

🔄 Example 3: Reader to Writer (Streaming)
   Reading from io.Reader and writing to io.Writer
   (Useful for streaming, proxying, etc.)
   ✅ Success!
      Streamed: 16104 bytes
      Format: μ-law @ 64.0 kbps
      Processing time: 10 ms

========================================
✅ All advanced examples completed!
========================================

💡 Use Cases:
   • FromReader: HTTP uploads, streaming input
   • ToWriter: HTTP responses, network streams
   • Combined: Real-time audio processing pipelines
```

## 💻 Method 1: TranscodeFromReader

### Signature

```go
func (t *transcoder) TranscodeFromReader(
    reader io.Reader,           // Input source
    outputPath string,          // Where to save result
    format AudioFormat,         // Target format
) (*TranscoderResult, error)
```

### Basic Example

```go
// Open file as io.Reader
inputFile, err := os.Open("input.wav")
if err != nil {
    log.Fatal(err)
}
defer inputFile.Close()

// Transcode from reader
result, err := transcoder.TranscodeFromReader(
    inputFile,                     // io.Reader
    "output-from-reader.alaw",     // Save to file
    wav2multi.FormatALaw,          // Target format
)

fmt.Printf("Converted: %d bytes\n", result.OutputFile.Size)
```

### Use Case 1: HTTP File Upload

```go
func handleUpload(w http.ResponseWriter, r *http.Request) {
    // Parse multipart form
    err := r.ParseMultipartForm(32 << 20) // 32MB max
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    
    // Get uploaded file
    file, header, err := r.FormFile("audio")
    if err != nil {
        http.Error(w, "No file uploaded", http.StatusBadRequest)
        return
    }
    defer file.Close()
    
    // Transcode directly from upload (no temp file needed!)
    transcoder := wav2multi.NewTranscoder(false)
    
    outputPath := fmt.Sprintf("uploads/%s.ulaw", 
        strings.TrimSuffix(header.Filename, ".wav"))
    
    result, err := transcoder.TranscodeFromReader(
        file,                      // Uploaded file is io.Reader
        outputPath,
        wav2multi.FormatULaw,
    )
    
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    // Return success
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]interface{}{
        "status":   "success",
        "filename": header.Filename,
        "output":   outputPath,
        "size":     result.OutputFile.Size,
        "format":   result.OutputFile.Type,
        "bitrate":  result.Stats.BitrateKbps,
    })
}
```

**Client usage:**
```bash
curl -X POST -F "audio=@recording.wav" http://localhost:8080/upload
```

### Use Case 2: Read from Memory

```go
// Audio data in memory (e.g., from database, cache, etc.)
wavData := []byte{
    // ... WAV file content in bytes ...
}

// Create reader from memory
reader := bytes.NewReader(wavData)

// Transcode without touching disk
result, err := transcoder.TranscodeFromReader(
    reader,
    "output.alaw",
    wav2multi.FormatALaw,
)

fmt.Printf("Processed from memory: %d bytes\n", result.OutputFile.Size)
```

### Use Case 3: Read from Network Stream

```go
func handleNetworkStream(conn net.Conn) {
    defer conn.Close()
    
    transcoder := wav2multi.NewTranscoder(false)
    
    // Read audio from network connection
    result, err := transcoder.TranscodeFromReader(
        conn,                      // net.Conn is an io.Reader
        "stream-output.g729",
        wav2multi.FormatG729,
    )
    
    if err != nil {
        log.Printf("Stream transcode failed: %v", err)
        return
    }
    
    log.Printf("Stream processed: %d bytes", result.OutputFile.Size)
}
```

### Use Case 4: Read from Compressed Archive

```go
import (
    "archive/zip"
    "github.com/lordbasex/wav2multi-lib"
)

func processZipFile(zipPath string) error {
    // Open ZIP archive
    r, err := zip.OpenReader(zipPath)
    if err != nil {
        return err
    }
    defer r.Close()
    
    transcoder := wav2multi.NewTranscoder(false)
    
    // Process each WAV file in the archive
    for _, file := range r.File {
        if !strings.HasSuffix(file.Name, ".wav") {
            continue
        }
        
        // Open file in archive
        rc, err := file.Open()
        if err != nil {
            log.Printf("Skip %s: %v", file.Name, err)
            continue
        }
        
        // Transcode directly from ZIP (no extraction!)
        outputName := strings.TrimSuffix(file.Name, ".wav") + ".ulaw"
        
        result, err := transcoder.TranscodeFromReader(
            rc,                     // Reader from ZIP file
            outputName,
            wav2multi.FormatULaw,
        )
        rc.Close()
        
        if err != nil {
            log.Printf("Failed %s: %v", file.Name, err)
            continue
        }
        
        log.Printf("✅ %s -> %s (%d bytes)", 
            file.Name, outputName, result.OutputFile.Size)
    }
    
    return nil
}
```

## 💻 Method 2: TranscodeToWriter

### Signature

```go
func (t *transcoder) TranscodeToWriter(
    inputPath string,           // Input file
    writer io.Writer,           // Output destination
    format AudioFormat,         // Target format
) (*TranscoderResult, error)
```

### Basic Example

```go
// Create a buffer as destination
var buffer bytes.Buffer

// Transcode to writer
result, err := transcoder.TranscodeToWriter(
    "input.wav",               // Input file
    &buffer,                   // io.Writer (buffer)
    wav2multi.FormatSLIN,      // Target format
)

// Audio is now in buffer.Bytes()
fmt.Printf("Buffer size: %d bytes\n", buffer.Len())

// Optionally save to disk
os.WriteFile("output.slin", buffer.Bytes(), 0644)
```

### Use Case 1: HTTP Download/Stream

```go
func handleDownload(w http.ResponseWriter, r *http.Request) {
    filename := r.URL.Query().Get("file")
    format := r.URL.Query().Get("format")
    
    if filename == "" || format == "" {
        http.Error(w, "Missing parameters", http.StatusBadRequest)
        return
    }
    
    // Set response headers
    w.Header().Set("Content-Type", "audio/basic")
    w.Header().Set("Content-Disposition", 
        fmt.Sprintf("attachment; filename=%s.%s", filename, format))
    
    transcoder := wav2multi.NewTranscoder(false)
    
    inputPath := fmt.Sprintf("recordings/%s.wav", filename)
    
    // Stream directly to HTTP response (no temp file!)
    result, err := transcoder.TranscodeToWriter(
        inputPath,
        w,                             // HTTP response is io.Writer
        wav2multi.AudioFormat(format),
    )
    
    if err != nil {
        // Can't send error response (headers already sent)
        log.Printf("Stream error: %v", err)
        return
    }
    
    log.Printf("Streamed %s: %d bytes to client", 
        filename, result.OutputFile.Size)
}
```

**Client usage:**
```bash
curl "http://localhost:8080/download?file=recording&format=ulaw" -o output.ulaw
```

### Use Case 2: Send to Network

```go
func sendToServer(inputFile string, serverAddr string) error {
    // Connect to server
    conn, err := net.Dial("tcp", serverAddr)
    if err != nil {
        return err
    }
    defer conn.Close()
    
    transcoder := wav2multi.NewTranscoder(false)
    
    // Send converted audio directly over network
    result, err := transcoder.TranscodeToWriter(
        inputFile,
        conn,                      // net.Conn is io.Writer
        wav2multi.FormatG729,
    )
    
    if err != nil {
        return err
    }
    
    log.Printf("Sent %d bytes to %s", result.OutputFile.Size, serverAddr)
    return nil
}
```

### Use Case 3: S3/Cloud Upload

```go
import (
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func uploadToS3(inputFile string, bucket string, key string) error {
    // Create S3 session
    sess := session.Must(session.NewSession())
    uploader := s3manager.NewUploader(sess)
    
    transcoder := wav2multi.NewTranscoder(false)
    
    // Create pipe for streaming
    pr, pw := io.Pipe()
    
    // Upload in goroutine
    var uploadErr error
    var uploadResult *s3manager.UploadOutput
    go func() {
        uploadResult, uploadErr = uploader.Upload(&s3manager.UploadInput{
            Bucket: aws.String(bucket),
            Key:    aws.String(key),
            Body:   pr,
        })
        pr.Close()
    }()
    
    // Transcode to pipe (which uploads to S3)
    result, err := transcoder.TranscodeToWriter(
        inputFile,
        pw,                        // Pipe writer
        wav2multi.FormatULaw,
    )
    pw.Close()
    
    if err != nil {
        return err
    }
    
    if uploadErr != nil {
        return uploadErr
    }
    
    log.Printf("Uploaded %d bytes to S3: %s", 
        result.OutputFile.Size, uploadResult.Location)
    return nil
}
```

### Use Case 4: Real-Time WebSocket Streaming

```go
import "github.com/gorilla/websocket"

func handleWebSocket(ws *websocket.Conn) {
    defer ws.Close()
    
    transcoder := wav2multi.NewTranscoder(false)
    
    // Create custom writer that sends to WebSocket
    wsWriter := &WebSocketWriter{conn: ws}
    
    // Stream audio to WebSocket client
    result, err := transcoder.TranscodeToWriter(
        "input.wav",
        wsWriter,                  // Custom io.Writer
        wav2multi.FormatULaw,
    )
    
    if err != nil {
        log.Printf("WebSocket stream error: %v", err)
        return
    }
    
    log.Printf("Streamed %d bytes via WebSocket", result.OutputFile.Size)
}

// Custom io.Writer for WebSocket
type WebSocketWriter struct {
    conn *websocket.Conn
}

func (w *WebSocketWriter) Write(p []byte) (n int, err error) {
    err = w.conn.WriteMessage(websocket.BinaryMessage, p)
    if err != nil {
        return 0, err
    }
    return len(p), nil
}
```

## 💻 Method 3: Combined (Reader → Writer)

### Use Case: HTTP Proxy/Gateway

```go
func handleConvert(w http.ResponseWriter, r *http.Request) {
    // Get format from query
    format := r.URL.Query().Get("format")
    if format == "" {
        format = "ulaw"
    }
    
    // Set response headers
    w.Header().Set("Content-Type", "audio/basic")
    w.Header().Set("Content-Disposition", 
        fmt.Sprintf("attachment; filename=converted.%s", format))
    
    transcoder := wav2multi.NewTranscoder(false)
    
    // Read from request body, write to response
    // (both are io.Reader/Writer interfaces)
    
    // Note: Currently we need a temp file, but this could be optimized
    tempFile := "/tmp/temp-upload.wav"
    
    // Save upload to temp
    f, err := os.Create(tempFile)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    _, err = io.Copy(f, r.Body)
    f.Close()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer os.Remove(tempFile)
    
    // Convert and stream to response
    result, err := transcoder.TranscodeToWriter(
        tempFile,
        w,                             // Stream to HTTP response
        wav2multi.AudioFormat(format),
    )
    
    if err != nil {
        log.Printf("Conversion error: %v", err)
        return
    }
    
    log.Printf("Proxied conversion: %d bytes", result.OutputFile.Size)
}
```

**Client usage:**
```bash
# Upload WAV, get converted audio back immediately
curl -X POST --data-binary "@input.wav" \
     "http://localhost:8080/convert?format=ulaw" \
     -o output.ulaw
```

## 📊 Comparison: Methods Overview

| Method | Input | Output | Use Case |
|--------|-------|--------|----------|
| `Transcode()` | File | File | Simple file conversion |
| `TranscodeFromReader()` | `io.Reader` | File | HTTP uploads, streams → file |
| `TranscodeToWriter()` | File | `io.Writer` | File → HTTP downloads, streams |
| Combined | `io.Reader` | `io.Writer` | Full streaming pipeline |

## 🔧 Advanced Patterns

### Pattern 1: Concurrent Processing

```go
func processMultipleStreams(inputs []io.Reader) {
    var wg sync.WaitGroup
    transcoder := wav2multi.NewTranscoder(false)
    
    for i, reader := range inputs {
        wg.Add(1)
        go func(idx int, r io.Reader) {
            defer wg.Done()
            
            outputPath := fmt.Sprintf("stream-%d.ulaw", idx)
            result, err := transcoder.TranscodeFromReader(
                r,
                outputPath,
                wav2multi.FormatULaw,
            )
            
            if err != nil {
                log.Printf("Stream %d failed: %v", idx, err)
                return
            }
            
            log.Printf("Stream %d: %d bytes", idx, result.OutputFile.Size)
        }(i, reader)
    }
    
    wg.Wait()
}
```

### Pattern 2: Progress Tracking

```go
type ProgressReader struct {
    reader io.Reader
    total  int64
    read   int64
}

func (pr *ProgressReader) Read(p []byte) (int, error) {
    n, err := pr.reader.Read(p)
    pr.read += int64(n)
    
    if pr.total > 0 {
        percent := float64(pr.read) / float64(pr.total) * 100
        fmt.Printf("\rProgress: %.1f%%", percent)
    }
    
    return n, err
}

// Usage
file, _ := os.Open("large-file.wav")
stat, _ := file.Stat()

progressReader := &ProgressReader{
    reader: file,
    total:  stat.Size(),
}

result, err := transcoder.TranscodeFromReader(
    progressReader,
    "output.ulaw",
    wav2multi.FormatULaw,
)
```

## 🐛 Common Issues

### Issue 1: "invalid input file" with io.Reader

**Cause:** Input doesn't meet WAV requirements

**Solution:** Validate before transcoding
```go
// Save to temp file for validation
tempFile, _ := os.CreateTemp("", "validate-*.wav")
io.Copy(tempFile, reader)
tempFile.Close()

// Validate
inputInfo, err := transcoder.ValidateInput(tempFile.Name())
if err != nil {
    return fmt.Errorf("invalid input: %w", err)
}

// Now transcode
result, err := transcoder.TranscodeFromReader(...)
```

### Issue 2: Memory issues with large files

**Solution:** Use streaming with buffers
```go
// Use buffered reader
bufferedReader := bufio.NewReaderSize(reader, 32*1024) // 32KB buffer

result, err := transcoder.TranscodeFromReader(
    bufferedReader,
    "output.ulaw",
    wav2multi.FormatULaw,
)
```

## 📚 Next Steps

1. ✅ Integrate into your HTTP API
2. ✅ Build a WebSocket streaming service
3. ✅ Create a cloud storage integration
4. ✅ Implement a microservice architecture
5. ✅ Add real-time audio processing

## 🔗 See Also

- [Example 1](../example1/) - Basic file-to-file conversion
- [Example 2](../example2/) - Multi-format batch processing
- [Main README](../../README.md) - Complete library documentation
- [CGO_SETUP.md](../../CGO_SETUP.md) - G.729 installation guide

---

**Author:** Federico Pereira <fpereira@cnsoluciones.com>  
**Library:** [wav2multi-lib](https://github.com/lordbasex/wav2multi-lib)
