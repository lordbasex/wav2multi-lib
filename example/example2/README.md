# Example 2: Convert to All Formats

This example demonstrates how to convert a single WAV file to **all supported audio formats** using a loop, with progress tracking and error handling per format.

## 📋 Overview

This example shows:
- ✅ How to iterate over multiple formats
- ✅ How to handle errors per format (graceful degradation)
- ✅ How to generate dynamic output filenames
- ✅ How to display a summary of results
- ✅ Progress tracking with counters

**Perfect for:** Batch processing, testing, production systems that need multiple formats

## 🚀 Running the Example

### Without CGO (μ-law, A-law, SLIN - G.729 will fail gracefully):
```bash
cd example/example2
CGO_ENABLED=0 go run main.go
```

### With CGO (all formats including G.729):
```bash
cd example/example2
export CGO_ENABLED=1
export CGO_CFLAGS="-I/usr/local/include"
export CGO_LDFLAGS="-L/usr/local/lib -lbcg729"
go run main.go
```

## 📁 Generated Files

After running, you'll get:
- `output.ulaw` - μ-law format (64 kbps)
- `output.alaw` - A-law format (64 kbps)
- `output.slin` - SLIN format (128 kbps)
- `output.g729` - G.729 format (8 kbps) - only with CGO

## 📝 Expected Output

### Without CGO (3 of 4 succeed):
```
🎵 wav2multi-lib - Example 2: Convert to All Formats
=====================================================

📝 [1/4] Converting to ulaw...
   ✅ Success!
      Output: output.ulaw (16104 bytes)
      Bitrate: 64.0 kbps
      Processing time: 50 ms
      Compression: 50.00%

📝 [2/4] Converting to alaw...
   ✅ Success!
      Output: output.alaw (16104 bytes)
      Bitrate: 64.0 kbps
      Processing time: 45 ms
      Compression: 50.00%

📝 [3/4] Converting to g729...
   ❌ Failed: G.729 encoder not available: codec not available

📝 [4/4] Converting to slin...
   ✅ Success!
      Output: output.slin (32208 bytes)
      Bitrate: 128.0 kbps
      Processing time: 42 ms
      Compression: 99.88%

========================================
✅ Completed: 3 successful, 1 failed
========================================

⚠️  Note: G.729 requires CGO and libbcg729.
   See CGO_SETUP.md for installation instructions.
```

### With CGO (all 4 succeed):
```
🎵 wav2multi-lib - Example 2: Convert to All Formats
=====================================================

📝 [1/4] Converting to ulaw...
   ✅ Success!
      Output: output.ulaw (16104 bytes)
      Bitrate: 64.0 kbps
      Processing time: 50 ms
      Compression: 50.00%

📝 [2/4] Converting to alaw...
   ✅ Success!
      Output: output.alaw (16104 bytes)
      Bitrate: 64.0 kbps
      Processing time: 45 ms
      Compression: 50.00%

📝 [3/4] Converting to g729...
   ✅ Success!
      Output: output.g729 (2020 bytes)
      Bitrate: 8.0 kbps
      Processing time: 15 ms
      Compression: 6.26%

📝 [4/4] Converting to slin...
   ✅ Success!
      Output: output.slin (32208 bytes)
      Bitrate: 128.0 kbps
      Processing time: 42 ms
      Compression: 99.88%

========================================
✅ Completed: 4 successful, 0 failed
========================================
```

## 💻 Code Walkthrough

### Step 1: Define Formats to Convert

```go
// Simple list of codec names
codecs := []string{"ulaw", "alaw", "g729", "slin"}
```

**Why this approach:**
- ✅ Simple and easy to understand
- ✅ Easy to add/remove formats
- ✅ Works directly with `AudioFormat` type

### Step 2: Track Success/Failure

```go
successCount := 0
failCount := 0
```

**Purpose:**
- Count successful conversions
- Count failures (like G.729 without CGO)
- Display final summary

### Step 3: Loop Through Formats

```go
for i, codec := range codecs {
    fmt.Printf("📝 [%d/%d] Converting to %s...\n", i+1, len(codecs), codec)
    
    // Build output filename dynamically
    outputFile := fmt.Sprintf("output.%s", codec)
    
    // Configure transcoding
    config := wav2multi.TranscoderConfig{
        InputPath:  inputFile,
        OutputPath: outputFile,
        Format:     wav2multi.AudioFormat(codec),
    }
    
    // Perform transcoding with error handling
    result, err := transcoder.Transcode(config)
    if err != nil {
        fmt.Printf("   ❌ Failed: %v\n\n", err)
        failCount++
        
        // Wait before next conversion
        if i < len(codecs)-1 {
            time.Sleep(1 * time.Second)
        }
        continue  // Don't stop, continue with next format
    }
    
    // Display success
    fmt.Printf("   ✅ Success!\n")
    fmt.Printf("      Output: %s (%d bytes)\n", result.OutputFile.Path, result.OutputFile.Size)
    successCount++
    
    // Wait before next conversion
    if i < len(codecs)-1 {
        time.Sleep(1 * time.Second)
    }
}
```

### Step 4: Display Summary

```go
fmt.Println("========================================")
fmt.Printf("✅ Completed: %d successful, %d failed\n", successCount, failCount)
fmt.Println("========================================")

if failCount > 0 {
    fmt.Println("\n⚠️  Note: G.729 requires CGO and libbcg729.")
}
```

## 🎓 Key Concepts

### 1. Graceful Error Handling

The example **doesn't stop** if one format fails:

```go
result, err := transcoder.Transcode(config)
if err != nil {
    fmt.Printf("❌ Failed: %v\n", err)
    failCount++
    continue  // Continue with next format
}
```

**Benefits:**
- ✅ G.729 fails without CGO, but others continue
- ✅ Partial success is better than total failure
- ✅ Production-ready behavior

### 2. Progress Tracking

```go
fmt.Printf("📝 [%d/%d] Converting to %s...\n", i+1, len(codecs), codec)
```

**Shows:**
- Current progress: `[1/4]`, `[2/4]`, etc.
- Which format is being processed
- User knows how much work remains

### 3. Dynamic Filename Generation

```go
outputFile := fmt.Sprintf("output.%s", codec)
```

**Pattern:**
- Input: `"ulaw"` → Output: `"output.ulaw"`
- Input: `"g729"` → Output: `"output.g729"`

### 4. Sleep Between Conversions

```go
time.Sleep(1 * time.Second)
```

**Why:**
- ✅ Prevents overwhelming the system
- ✅ Makes output easier to read
- ✅ Useful for rate-limiting in production

## 📚 Practical Use Cases

### Use Case 1: Process Multiple Input Files

```go
package main

import (
    "fmt"
    "path/filepath"
    "github.com/lordbasex/wav2multi-lib"
)

func main() {
    transcoder := wav2multi.NewTranscoder(false)
    
    // List of input files
    inputFiles := []string{
        "recording-001.wav",
        "recording-002.wav",
        "recording-003.wav",
    }
    
    codecs := []string{"ulaw", "alaw", "g729", "slin"}
    
    for _, inputFile := range inputFiles {
        fmt.Printf("\n📂 Processing: %s\n", inputFile)
        
        // Get base name without extension
        baseName := filepath.Base(inputFile)
        baseName = baseName[:len(baseName)-len(filepath.Ext(baseName))]
        
        for _, codec := range codecs {
            // Dynamic output: recording-001.ulaw, recording-001.alaw, etc.
            outputFile := fmt.Sprintf("%s.%s", baseName, codec)
            
            config := wav2multi.TranscoderConfig{
                InputPath:  inputFile,
                OutputPath: outputFile,
                Format:     wav2multi.AudioFormat(codec),
            }
            
            result, err := transcoder.Transcode(config)
            if err != nil {
                fmt.Printf("  ❌ %s: %v\n", codec, err)
                continue
            }
            
            fmt.Printf("  ✅ %s: %d bytes\n", codec, result.OutputFile.Size)
        }
    }
}
```

**Output structure:**
```
recording-001.ulaw
recording-001.alaw
recording-001.g729
recording-001.slin
recording-002.ulaw
recording-002.alaw
...
```

### Use Case 2: Organize by Directory

```go
import "os"

codecs := []string{"ulaw", "alaw", "g729", "slin"}

for _, codec := range codecs {
    // Create directory for each format
    dir := fmt.Sprintf("converted/%s", codec)
    os.MkdirAll(dir, 0755)
    
    // Save in format-specific directory
    outputFile := fmt.Sprintf("%s/output.%s", dir, codec)
    
    config := wav2multi.TranscoderConfig{
        InputPath:  "input.wav",
        OutputPath: outputFile,
        Format:     wav2multi.AudioFormat(codec),
    }
    
    result, err := transcoder.Transcode(config)
    // ...
}
```

**Output structure:**
```
converted/
  ulaw/
    output.ulaw
  alaw/
    output.alaw
  g729/
    output.g729
  slin/
    output.slin
```

### Use Case 3: Conditional Format Selection

```go
transcoder := wav2multi.NewTranscoder(false)

// Get available formats based on CGO
availableFormats := transcoder.GetSupportedFormats()

fmt.Printf("Converting to %d available formats\n", len(availableFormats))

for _, format := range availableFormats {
    outputFile := fmt.Sprintf("output.%s", format)
    
    config := wav2multi.TranscoderConfig{
        InputPath:  "input.wav",
        OutputPath: outputFile,
        Format:     format,
    }
    
    result, err := transcoder.Transcode(config)
    if err != nil {
        fmt.Printf("❌ %s failed: %v\n", format, err)
        continue
    }
    
    fmt.Printf("✅ %s: %d bytes\n", format, result.OutputFile.Size)
}
```

### Use Case 4: Production Batch Processing

```go
type ConversionJob struct {
    InputFile string
    OutputDir string
    Formats   []string
}

func ProcessBatch(jobs []ConversionJob) {
    transcoder := wav2multi.NewTranscoder(false)
    
    totalJobs := 0
    successJobs := 0
    
    for _, job := range jobs {
        for _, codec := range job.Formats {
            totalJobs++
            
            outputPath := filepath.Join(
                job.OutputDir, 
                fmt.Sprintf("%s.%s", 
                    filepath.Base(job.InputFile),
                    codec,
                ),
            )
            
            config := wav2multi.TranscoderConfig{
                InputPath:  job.InputFile,
                OutputPath: outputPath,
                Format:     wav2multi.AudioFormat(codec),
            }
            
            _, err := transcoder.Transcode(config)
            if err != nil {
                log.Printf("ERROR: %s -> %s: %v", 
                    job.InputFile, codec, err)
                continue
            }
            
            successJobs++
        }
    }
    
    log.Printf("Batch complete: %d/%d successful", successJobs, totalJobs)
}
```

## 🔧 Customization Examples

### Add Timestamp to Filenames

```go
import "time"

timestamp := time.Now().Format("20060102-150405")

for _, codec := range codecs {
    outputFile := fmt.Sprintf("output-%s.%s", timestamp, codec)
    // Results in: output-20250122-143052.ulaw
}
```

### Filter Formats by Size Requirement

```go
// Only use formats under 20KB
codecs := []string{"ulaw", "alaw", "g729", "slin"}
maxSize := 20000 // bytes

for _, codec := range codecs {
    config := wav2multi.TranscoderConfig{
        InputPath:  "input.wav",
        OutputPath: fmt.Sprintf("output.%s", codec),
        Format:     wav2multi.AudioFormat(codec),
    }
    
    result, err := transcoder.Transcode(config)
    if err != nil {
        continue
    }
    
    if result.OutputFile.Size > int64(maxSize) {
        fmt.Printf("⚠️  %s too large (%d bytes), skipping\n", 
            codec, result.OutputFile.Size)
        os.Remove(result.OutputFile.Path)  // Delete if too large
        continue
    }
    
    fmt.Printf("✅ %s: %d bytes (within limit)\n", 
        codec, result.OutputFile.Size)
}
```

### Parallel Processing

```go
import "sync"

var wg sync.WaitGroup
codecs := []string{"ulaw", "alaw", "slin"}  // Exclude G.729 for simplicity

for _, codec := range codecs {
    wg.Add(1)
    go func(c string) {
        defer wg.Done()
        
        transcoder := wav2multi.NewTranscoder(false)
        config := wav2multi.TranscoderConfig{
            InputPath:  "input.wav",
            OutputPath: fmt.Sprintf("output.%s", c),
            Format:     wav2multi.AudioFormat(c),
        }
        
        result, err := transcoder.Transcode(config)
        if err != nil {
            fmt.Printf("❌ %s: %v\n", c, err)
            return
        }
        
        fmt.Printf("✅ %s: %d bytes\n", c, result.OutputFile.Size)
    }(codec)
}

wg.Wait()
fmt.Println("All conversions complete!")
```

## 📊 Format Comparison

| Format | Bitrate | File Size* | Compression | Processing** | Best For |
|--------|---------|-----------|-------------|--------------|----------|
| **G.729** | 8 kbps | ~2 KB | 93.7% | 15 ms | VoIP, bandwidth-critical |
| **μ-law** | 64 kbps | ~16 KB | 50.0% | 50 ms | US telephony systems |
| **A-law** | 64 kbps | ~16 KB | 50.0% | 45 ms | European telephony |
| **SLIN** | 128 kbps | ~32 KB | 0.1% | 42 ms | Raw PCM, debugging |

*For a 2-second audio file  
**Approximate processing time

## 🐛 Common Issues

### Issue 1: G.729 fails without CGO

**Expected behavior:** G.729 fails gracefully, others continue

**Solution:** This is normal. Install libbcg729 for G.729 support:
```bash
# See CGO_SETUP.md
export CGO_ENABLED=1
export CGO_CFLAGS="-I/usr/local/include"
export CGO_LDFLAGS="-L/usr/local/lib -lbcg729"
```

### Issue 2: Output files already exist

**Behavior:** Existing files are overwritten

**Solution:** Add timestamp or check before writing:
```go
outputFile := fmt.Sprintf("output.%s", codec)

if _, err := os.Stat(outputFile); err == nil {
    fmt.Printf("⚠️  %s already exists, skipping\n", outputFile)
    continue
}
```

### Issue 3: Disk space issues with many formats

**Solution:** Check available space before converting:
```go
import "syscall"

func checkDiskSpace(path string) (uint64, error) {
    var stat syscall.Statfs_t
    err := syscall.Statfs(path, &stat)
    if err != nil {
        return 0, err
    }
    return stat.Bavail * uint64(stat.Bsize), nil
}

// Before conversion
freeSpace, _ := checkDiskSpace(".")
if freeSpace < 100*1024*1024 {  // Less than 100MB
    log.Fatal("Insufficient disk space")
}
```

## 📚 Next Steps

After mastering Example 2, try:

1. ✅ **[Example 3](../example3/)** - Advanced usage with `io.Reader`/`io.Writer`
2. ✅ Process multiple input files in a batch
3. ✅ Add parallel processing for faster conversion
4. ✅ Implement a queuing system for large batches
5. ✅ Add database logging for conversions

## 🔗 See Also

- [Example 1](../example1/) - Basic single-format conversion
- [Example 3](../example3/) - Advanced I/O with streaming
- [Main README](../../README.md) - Complete documentation
- [CGO_SETUP.md](../../CGO_SETUP.md) - G.729 setup guide

---

**Author:** Federico Pereira <fpereira@cnsoluciones.com>  
**Library:** [wav2multi-lib](https://github.com/lordbasex/wav2multi-lib)
