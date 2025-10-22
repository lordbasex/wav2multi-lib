# Example 1: Basic Usage

This example demonstrates the simplest way to use wav2multi-lib: converting a single WAV file to one audio format.

## 📋 Overview

This example shows:
- ✅ How to create a transcoder
- ✅ How to configure a simple conversion
- ✅ How to handle errors
- ✅ How to read conversion results
- ✅ Basic transcoding workflow

**Perfect for:** Beginners learning the library, simple CLI tools, single-format conversions

## 🚀 Running the Example

### Without CGO (μ-law, A-law, SLIN):
```bash
cd example/example1
CGO_ENABLED=0 go run main.go
```

### With CGO (all formats including G.729):
```bash
cd example/example1
export CGO_ENABLED=1
export CGO_CFLAGS="-I/usr/local/include"
export CGO_LDFLAGS="-L/usr/local/lib -lbcg729"
go run main.go
```

## 📁 Files

- `main.go` - Example source code
- `input.wav` - Input file (8kHz, mono, 16-bit PCM)
- `output.ulaw` - Generated output file (created when running)
- `go.mod` - Module dependencies

## 📝 Expected Output

```
🎵 wav2multi-lib - Example 1: Basic Usage
==========================================

📝 Converting to μ-law format...
=== TRANSCODING RESULT ===
Input:   (32252 bytes, 2.01 seconds)
Output: output.ulaw (16104 bytes)
Format: ulaw (64.0 kbps)
Processing: 50 ms
Compression: 50.00%
Samples: 16104
========================

✅ Conversion completed successfully!
   Input:  input.wav
   Output: output.ulaw (16104 bytes)
   Format: ulaw @ 64.0 kbps
   Processing time: 50 ms
   Compression ratio: 50.00%
```

## 💻 Code Walkthrough

### Step 1: Create a Transcoder

```go
transcoder := wav2multi.NewTranscoder(true) // verbose = true for detailed logs
```

**What it does:**
- Creates a new transcoder instance
- `verbose = true` enables detailed processing logs
- Returns a `Transcoder` interface

### Step 2: Configure the Conversion

```go
config := wav2multi.TranscoderConfig{
    InputPath:  "input.wav",      // Source WAV file
    OutputPath: "output.ulaw",     // Destination file
    Format:     wav2multi.FormatULaw,  // Target codec
}
```

**Configuration fields:**
- `InputPath`: Path to your input WAV file
- `OutputPath`: Where to save the converted audio (can include directories)
- `Format`: Which codec to use (`ulaw`, `alaw`, `slin`, `g729`)

### Step 3: Execute Transcoding

```go
result, err := transcoder.Transcode(config)
if err != nil {
    log.Fatalf("❌ Error during transcoding: %v", err)
}
```

**What happens:**
1. Validates input file (format, sample rate, channels)
2. Reads WAV data
3. Encodes to target format
4. Writes output file
5. Returns statistics

### Step 4: Read Results

```go
fmt.Printf("Output: %s (%d bytes)\n", 
    result.OutputFile.Path,      // Output file path
    result.OutputFile.Size)      // Output file size

fmt.Printf("Format: %s @ %.1f kbps\n", 
    result.OutputFile.Type,      // Format name
    result.Stats.BitrateKbps)    // Bitrate

fmt.Printf("Processing time: %d ms\n", 
    result.Stats.ProcessingTimeMs)  // Time taken

fmt.Printf("Compression: %.2f%%\n", 
    result.Stats.CompressionRatio*100)  // Compression ratio
```

## 🎓 Key Concepts

### 1. Verbose Logging

```go
// With verbose logging (shows all details)
transcoder := wav2multi.NewTranscoder(true)

// Without verbose logging (quiet mode)
transcoder := wav2multi.NewTranscoder(false)
```

**When to use verbose:**
- ✅ Debugging
- ✅ Learning the library
- ✅ Troubleshooting issues

**When to use quiet:**
- ✅ Production environments
- ✅ Batch processing
- ✅ When you handle your own logging

### 2. Format Selection

```go
// μ-law (64 kbps) - US telephony standard
Format: wav2multi.FormatULaw,

// A-law (64 kbps) - European telephony standard
Format: wav2multi.FormatALaw,

// SLIN (128 kbps) - Raw PCM, no compression
Format: wav2multi.FormatSLIN,

// G.729 (8 kbps) - High compression (requires CGO)
Format: wav2multi.FormatG729,
```

### 3. Error Handling

```go
result, err := transcoder.Transcode(config)
if err != nil {
    // Common errors:
    // - Invalid input file format
    // - Wrong sample rate (must be 8000 Hz)
    // - Not mono (must be 1 channel)
    // - Codec not available (G.729 without CGO)
    log.Fatalf("Transcoding failed: %v", err)
}
```

## 📚 Practical Use Cases

### Use Case 1: Simple CLI Tool

```go
package main

import (
    "flag"
    "log"
    "github.com/lordbasex/wav2multi-lib"
)

func main() {
    input := flag.String("input", "", "Input WAV file")
    output := flag.String("output", "", "Output file")
    format := flag.String("format", "ulaw", "Format (ulaw/alaw/slin/g729)")
    flag.Parse()

    transcoder := wav2multi.NewTranscoder(false)
    
    config := wav2multi.TranscoderConfig{
        InputPath:  *input,
        OutputPath: *output,
        Format:     wav2multi.AudioFormat(*format),
    }
    
    result, err := transcoder.Transcode(config)
    if err != nil {
        log.Fatal(err)
    }
    
    log.Printf("✅ Converted to %s (%d bytes)", 
        result.OutputFile.Type, 
        result.OutputFile.Size)
}
```

**Usage:**
```bash
./convert -input voice.wav -output voice.ulaw -format ulaw
```

### Use Case 2: Directory with Output Path

```go
config := wav2multi.TranscoderConfig{
    InputPath:  "recordings/call-001.wav",
    OutputPath: "converted/call-001.ulaw",  // Directory + filename
    Format:     wav2multi.FormatULaw,
}

// Note: Make sure "converted/" directory exists
os.MkdirAll("converted", 0755)

result, err := transcoder.Transcode(config)
```

### Use Case 3: Validate Before Converting

```go
transcoder := wav2multi.NewTranscoder(false)

// Validate input first
inputInfo, err := transcoder.ValidateInput("input.wav")
if err != nil {
    log.Fatalf("Invalid input: %v", err)
}

log.Printf("Input validated: %s (%d bytes, %.2f seconds)",
    inputInfo.Path,
    inputInfo.Size,
    inputInfo.DurationSeconds)

// Now convert
config := wav2multi.TranscoderConfig{
    InputPath:  "input.wav",
    OutputPath: "output.ulaw",
    Format:     wav2multi.FormatULaw,
}

result, err := transcoder.Transcode(config)
```

### Use Case 4: Check Supported Formats

```go
transcoder := wav2multi.NewTranscoder(false)

// Get list of supported formats
formats := transcoder.GetSupportedFormats()

log.Println("Available formats:")
for _, format := range formats {
    log.Printf("  - %s", format)
}

// Output (without CGO):
// Available formats:
//   - ulaw
//   - alaw
//   - slin

// Output (with CGO):
// Available formats:
//   - g729
//   - ulaw
//   - alaw
//   - slin
```

## 🔧 Customization Examples

### Change Output Format

```go
// Convert the same input to different formats
formats := []wav2multi.AudioFormat{
    wav2multi.FormatULaw,
    wav2multi.FormatALaw,
    wav2multi.FormatSLIN,
}

for _, format := range formats {
    config := wav2multi.TranscoderConfig{
        InputPath:  "input.wav",
        OutputPath: fmt.Sprintf("output.%s", format),
        Format:     format,
    }
    
    result, err := transcoder.Transcode(config)
    if err != nil {
        log.Printf("Failed %s: %v", format, err)
        continue
    }
    
    log.Printf("✅ %s: %d bytes", format, result.OutputFile.Size)
}
```

### Add Error Recovery

```go
config := wav2multi.TranscoderConfig{
    InputPath:  "input.wav",
    OutputPath: "output.ulaw",
    Format:     wav2multi.FormatULaw,
}

// Retry logic
maxRetries := 3
var result *wav2multi.TranscoderResult
var err error

for i := 0; i < maxRetries; i++ {
    result, err = transcoder.Transcode(config)
    if err == nil {
        break
    }
    log.Printf("Attempt %d failed: %v", i+1, err)
    time.Sleep(time.Second)
}

if err != nil {
    log.Fatalf("All retries failed: %v", err)
}
```

## 📊 Format Comparison

| Format | Bitrate | File Size* | Compression | Best For |
|--------|---------|-----------|-------------|----------|
| **G.729** | 8 kbps | ~2 KB | 93.7% | VoIP, bandwidth-limited |
| **μ-law** | 64 kbps | ~16 KB | 50.0% | US telephony |
| **A-law** | 64 kbps | ~16 KB | 50.0% | European telephony |
| **SLIN** | 128 kbps | ~32 KB | 0.1% | Raw PCM, debugging |

*For a 2-second audio file

## 🐛 Common Issues

### Issue 1: "invalid input file"

**Cause:** Input file doesn't meet requirements

**Solution:**
```bash
# Check file format
file input.wav

# Should output:
# WAVE audio, Microsoft PCM, 16 bit, mono 8000 Hz

# Convert if needed
ffmpeg -i input.mp3 -ar 8000 -ac 1 -sample_fmt s16 output.wav
```

### Issue 2: "G.729 encoder not available"

**Cause:** G.729 requires CGO and libbcg729

**Solution:**
```bash
# See CGO_SETUP.md for detailed instructions
export CGO_ENABLED=1
export CGO_CFLAGS="-I/usr/local/include"
export CGO_LDFLAGS="-L/usr/local/lib -lbcg729"
```

### Issue 3: Output directory doesn't exist

**Cause:** `OutputPath` includes a directory that doesn't exist

**Solution:**
```go
import "path/filepath"

outputDir := filepath.Dir(config.OutputPath)
os.MkdirAll(outputDir, 0755)  // Create directory first

result, err := transcoder.Transcode(config)
```

## 📚 Next Steps

After mastering Example 1, try:

1. ✅ **[Example 2](../example2/)** - Convert to multiple formats in a loop
2. ✅ **[Example 3](../example3/)** - Advanced usage with `io.Reader`/`io.Writer`
3. ✅ Modify the format to try different codecs
4. ✅ Process multiple input files
5. ✅ Build a CLI tool with flags

## 🔗 See Also

- [Example 2](../example2/) - Multi-format conversion
- [Example 3](../example3/) - Advanced I/O patterns
- [Main README](../../README.md) - Complete documentation
- [CGO_SETUP.md](../../CGO_SETUP.md) - G.729 setup guide

---

**Author:** Federico Pereira <fpereira@cnsoluciones.com>  
**Library:** [wav2multi-lib](https://github.com/lordbasex/wav2multi-lib)
