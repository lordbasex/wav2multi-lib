# Examples / Ejemplos

This directory contains practical examples demonstrating how to use wav2multi-lib in different scenarios.

## 📚 Available Examples

### 🎯 [Example 1: Basic Usage](example1/)
**Single-format conversion** - Demonstrates the simplest way to convert a WAV file to one audio format.

- **Ideal for:** Beginners, simple conversion tasks, learning the library
- **Concepts:** Creating transcoder, basic configuration, error handling
- **Output:** One output file
- **Difficulty:** ⭐ Beginner

```bash
cd example1
go run main.go
```

**What you'll learn:**
- How to create a transcoder with verbose logging
- How to configure a basic conversion
- How to read and display results
- Understanding the TranscoderConfig struct

---

### 🔄 [Example 2: Convert to All Formats](example2/)
**Multi-format batch conversion** - Demonstrates how to convert to multiple formats using a loop.

- **Ideal for:** Batch processing, testing, production systems
- **Concepts:** Iteration, per-format error handling, progress tracking
- **Output:** 4 files (G.729, μ-law, A-law, SLIN)
- **Difficulty:** ⭐⭐ Intermediate

```bash
cd example2
go run main.go
```

**What you'll learn:**
- How to iterate over multiple formats efficiently
- Graceful error handling (continues on failure)
- Dynamic filename generation
- Success/failure tracking and reporting
- Adding delays between conversions

---

### 🚀 [Example 3: Advanced Usage - io.Reader/Writer](example3/)
**Streaming and I/O interfaces** - Demonstrates advanced usage with Go's `io.Reader` and `io.Writer` interfaces.

- **Ideal for:** Web APIs, microservices, streaming, network applications
- **Concepts:** `io.Reader`, `io.Writer`, streaming, HTTP integration
- **Methods:** `TranscodeFromReader()`, `TranscodeToWriter()`
- **Difficulty:** ⭐⭐⭐ Advanced

```bash
cd example3
go run main.go
```

**What you'll learn:**
- How to use `TranscodeFromReader()` for flexible input
- How to use `TranscodeToWriter()` for flexible output
- HTTP file upload/download patterns
- In-memory processing without disk I/O
- Network streaming scenarios

---

## 🚀 Quick Start

### Option 1: Without CGO (μ-law, A-law, SLIN)
```bash
# Example 1
cd example1
CGO_ENABLED=0 go run main.go

# Example 2
cd example2
CGO_ENABLED=0 go run main.go

# Example 3
cd example3
CGO_ENABLED=0 go run main.go
```

### Option 2: With CGO (All formats including G.729)
```bash
# Configure CGO
export CGO_ENABLED=1
export CGO_CFLAGS="-I/usr/local/include"
export CGO_LDFLAGS="-L/usr/local/lib -lbcg729"

# Example 1
cd example1
go run main.go

# Example 2
cd example2
go run main.go

# Example 3
cd example3
go run main.go
```

## 📝 Structure of Each Example

```
example1/
├── main.go      # Example source code
├── input.wav    # Input audio file (included)
├── go.mod       # Module dependencies
└── README.md    # Detailed documentation

example2/
├── main.go      # Example source code
├── input.wav    # Input audio file (included)
├── go.mod       # Module dependencies
└── README.md    # Detailed documentation

example3/
├── main.go      # Example source code
├── input.wav    # Input audio file (included)
├── go.mod       # Module dependencies
└── README.md    # Detailed documentation
```

## 🎓 Learning Path

### For Beginners:
1. **Start with Example 1** - Learn the basics
2. **Experiment with formats** - Try changing `FormatULaw` to `FormatALaw`, etc.
3. **Add error handling** - Improve the error messages
4. **Try Example 2** - See multi-format conversion

### For Intermediate Users:
1. **Master Example 2** - Understand loop-based processing
2. **Modify for batch** - Process multiple input files
3. **Add progress tracking** - Implement counters and timers
4. **Try Example 3** - Learn advanced I/O patterns

### For Advanced Users:
1. **Master Example 3** - Understand io.Reader/Writer
2. **Build HTTP API** - Create a web service for conversion
3. **Implement streaming** - Real-time audio processing
4. **Optimize performance** - Parallel processing, caching

## 🎯 What You'll Learn

### In Example 1:
- ✅ Create a transcoder instance
- ✅ Configure a simple conversion
- ✅ Handle basic errors
- ✅ Read conversion results
- ✅ Display processing statistics

### In Example 2:
- ✅ Iterate over multiple formats
- ✅ Generate dynamic filenames
- ✅ Handle errors per format (graceful degradation)
- ✅ Track success/failure counts
- ✅ Display comprehensive summaries
- ✅ Implement delays between operations

### In Example 3:
- ✅ Use `TranscodeFromReader()` with `io.Reader`
- ✅ Use `TranscodeToWriter()` with `io.Writer`
- ✅ Process audio in streaming mode
- ✅ Integrate with HTTP APIs
- ✅ Implement in-memory processing
- ✅ Build network-based audio services

## 📊 Method Comparison

| Method | Input | Output | Example | Use Case |
|--------|-------|--------|---------|----------|
| `Transcode()` | File path | File path | 1, 2 | Simple file conversion |
| `TranscodeFromReader()` | `io.Reader` | File path | 3 | HTTP uploads, streams → file |
| `TranscodeToWriter()` | File path | `io.Writer` | 3 | File → HTTP downloads, streams |

## 🔧 Requirements

### Minimum Requirements (Without G.729):
- Go 1.21 or higher
- wav2multi-lib

### For G.729 Support:
- CGO enabled
- libbcg729 installed
- See [CGO_SETUP.md](../CGO_SETUP.md) for installation instructions

## 📊 Supported Formats

| Format | Bitrate | Compression | CGO Required | Best For |
|--------|---------|-------------|--------------|----------|
| G.729  | 8 kbps  | 93.7%       | ✅ Yes       | VoIP, maximum compression |
| μ-law  | 64 kbps | 50.0%       | ❌ No        | US telephony |
| A-law  | 64 kbps | 50.0%       | ❌ No        | European telephony |
| SLIN   | 128 kbps| 0.1%        | ❌ No        | Raw PCM, debugging |

## 🐛 Troubleshooting

### Error: "G.729 encoder not available"
**Solution:** G.729 requires CGO and libbcg729. See [CGO_SETUP.md](../CGO_SETUP.md).

```bash
export CGO_ENABLED=1
export CGO_CFLAGS="-I/usr/local/include"
export CGO_LDFLAGS="-L/usr/local/lib -lbcg729"
```

### Error: "invalid input file"
**Solution:** Verify your WAV file meets these requirements:
- Format: PCM
- Channels: Mono (1)
- Sample Rate: 8000 Hz
- Bit Depth: 16-bit

```bash
# Check file format
file input.wav

# Convert if needed
ffmpeg -i input.mp3 -ar 8000 -ac 1 -sample_fmt s16 output.wav
```

### Error: "missing go.sum entry"
**Solution:** Run `go mod tidy` in the example directory.

```bash
cd example1
go mod tidy
```

## 📚 More Information

- [Main README](../README.md) - Complete library documentation
- [CGO_SETUP.md](../CGO_SETUP.md) - G.729 installation guide
- [CHANGELOG.md](../CHANGELOG.md) - Version history

## 💡 Next Steps

1. ✅ Try Example 1 to understand the basics
2. ✅ Try Example 2 for multi-format conversion
3. ✅ Try Example 3 for advanced I/O patterns (io.Reader/Writer)
4. ✅ Modify the examples for your specific needs
5. ✅ Read the [complete documentation](../README.md)
6. ✅ Integrate into your project

## 🎯 Real-World Applications

### Example 1 is perfect for:
- Simple CLI tools
- Single-format conversion scripts
- Learning and experimentation
- Quick file conversions

### Example 2 is perfect for:
- Batch processing systems
- Testing environments
- Multi-format output requirements
- Production transcoding pipelines

### Example 3 is perfect for:
- HTTP/REST APIs
- Microservices
- WebSocket streaming
- Cloud storage integration
- Serverless functions
- Real-time audio processing

## 🤝 Contributing

Have a useful example? Share it!
- Open an issue in the repository
- Submit a pull request
- Contact the author

## 📖 Code Quality

All examples follow Go best practices:
- ✅ Clean, readable code
- ✅ Comprehensive comments
- ✅ Error handling
- ✅ Proper resource cleanup (`defer`)
- ✅ Standard Go project structure

## 🔗 Related Resources

### Documentation:
- [API Reference](../README.md#-api-reference)
- [Input Requirements](../README.md#-input-requirements)
- [Error Handling](../README.md#-error-handling)

### Setup:
- [Installation](../README.md#-installation)
- [CGO Setup](../CGO_SETUP.md)
- [Testing](../README.md#-testing)

---

**Author:** Federico Pereira <fpereira@cnsoluciones.com>  
**Library:** [wav2multi-lib](https://github.com/lordbasex/wav2multi-lib)

