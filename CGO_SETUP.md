# 🔧 CGO Setup for G.729 Support

This library supports G.729 encoding/decoding through CGO integration with `libbcg729`.

## 📋 Prerequisites

### 1. Install libbcg729

#### Ubuntu/Debian:
```bash
sudo apt-get update
sudo apt-get install libbcg729-dev
```

#### CentOS/RHEL:
```bash
sudo yum install libbcg729-devel
```

#### macOS (with Homebrew):
```bash
brew install bcg729
```

#### From Source:
```bash
git clone https://github.com/BelledonneCommunications/bcg729.git
cd bcg729
mkdir build && cd build
cmake .. -DCMAKE_INSTALL_PREFIX=/usr/local
make && sudo make install
```

### 2. Set CGO Environment Variables

```bash
export CGO_ENABLED=1
export CGO_CFLAGS="-I/usr/local/include"
export CGO_LDFLAGS="-L/usr/local/lib -lbcg729"
```

## 🚀 Usage

### With CGO (G.729 support):
```go
package main

import (
    "fmt"
    "log"
    "github.com/lordbasex/wav2multi-lib"
)

func main() {
    // Create transcoder with verbose logging
    transcoder := wav2multi.NewTranscoder(true)
    
    // Configure G.729 conversion
    config := wav2multi.TranscoderConfig{
        InputPath:  "input.wav",
        OutputPath: "output.g729",
        Format:     wav2multi.FormatG729,
    }
    
    // Transcode
    result, err := transcoder.Transcode(config)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("✅ Converted to G.729: %d bytes (%.1f kbps)\n", 
        result.OutputFile.Size, result.Stats.BitrateKbps)
}
```

### Without CGO (no G.729 support):
```go
package main

import (
    "fmt"
    "log"
    "github.com/lordbasex/wav2multi-lib"
)

func main() {
    // Create transcoder with verbose logging
    transcoder := wav2multi.NewTranscoder(true)
    
    // Configure μ-law conversion
    // Only μ-law, A-law, and SLIN will work without CGO
    config := wav2multi.TranscoderConfig{
        InputPath:  "input.wav",
        OutputPath: "output.ulaw",
        Format:     wav2multi.FormatULaw,
    }
    
    // Transcode
    result, err := transcoder.Transcode(config)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("✅ Converted to μ-law: %d bytes (%.1f kbps)\n", 
        result.OutputFile.Size, result.Stats.BitrateKbps)
}
```

## 🔍 Build Tags

The library uses build tags to handle CGO availability:

- **With CGO**: `// +build cgo` - Full G.729 support
- **Without CGO**: `// +build !cgo` - μ-law, A-law, SLIN only

## 🐳 Docker Usage

### With G.729 support:
```dockerfile
FROM golang:1.23-alpine AS builder

# Install build dependencies
RUN apk add --no-cache \
    gcc \
    musl-dev \
    cmake \
    make \
    git

# Install libbcg729
RUN git clone https://github.com/BelledonneCommunications/bcg729.git /tmp/bcg729 && \
    cd /tmp/bcg729 && \
    mkdir build && cd build && \
    cmake .. -DCMAKE_INSTALL_PREFIX=/usr/local && \
    make && make install

# Build the application
WORKDIR /app
COPY . .
RUN CGO_ENABLED=1 go build -o app .

FROM alpine:latest
RUN apk add --no-cache libstdc++
COPY --from=builder /usr/local/lib/libbcg729.so* /usr/local/lib/
COPY --from=builder /app/app /app/
CMD ["/app/app"]
```

### Without G.729 support:
```dockerfile
FROM golang:1.23-alpine AS builder

WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 go build -o app .

FROM alpine:latest
COPY --from=builder /app/app /app/
CMD ["/app/app"]
```

## 🧪 Testing

### Test G.729 availability:
```go
package main

import (
    "fmt"
    "github.com/lordbasex/wav2multi-lib"
)

func main() {
    encoder, err := wav2multi.GetEncoder(wav2multi.FormatG729)
    if err != nil {
        fmt.Printf("❌ G.729 not available: %v\n", err)
        fmt.Println("   Install libbcg729 and enable CGO")
        fmt.Println("   See CGO_SETUP.md for instructions")
    } else {
        fmt.Println("✅ G.729 encoder available!")
        fmt.Printf("   Bitrate: %.1f kbps\n", encoder.GetBitrate())
    }
}
```

### Quick Example Directory Check:
```bash
# Test example1 (basic usage)
cd example/example1
go run main.go

# Test example2 (all formats)
cd example/example2
go run main.go
```

## 📚 API Reference

### G.729 Specific Functions:

```go
// Create G.729 encoder
encoder, err := wav2multi.NewG729Encoder()
if err != nil {
    // G.729 not available
}

// Encode samples
err = encoder.Encode(samples, writer)

// Close encoder
encoder.Close()

// Create G.729 decoder
decoder, err := wav2multi.NewG729Decoder()
if err != nil {
    // G.729 not available
}

// Decode G.729 data
err = decoder.Decode(reader, writer)

// Close decoder
decoder.Close()
```

## ⚠️ Troubleshooting

### Common Issues:

1. **"G.729 encoder not available"**
   - Install libbcg729
   - Set CGO_ENABLED=1
   - Check CGO_CFLAGS and CGO_LDFLAGS

2. **"undefined: C.initBcg729EncoderChannel"**
   - libbcg729 not found
   - Check include paths

3. **"ld: library not found"**
   - libbcg729 not in library path
   - Check CGO_LDFLAGS

### Verify Installation:
```bash
# Check if libbcg729 is installed
pkg-config --exists bcg729 && echo "libbcg729 found" || echo "libbcg729 not found"

# Check CGO
go env CGO_ENABLED
go env CGO_CFLAGS
go env CGO_LDFLAGS
```

## 🎯 Summary

- **With CGO + libbcg729**: Full support for G.729, μ-law, A-law, SLIN
- **Without CGO**: Support for μ-law, A-law, SLIN only (G.729 will fail gracefully)
- **Docker**: Use multi-stage builds for CGO support
- **Testing**: Check encoder availability before use
- **Examples**: See `example/example1/` and `example/example2/` for working code

## 📚 See Also

- [example/example1/](example/example1/) - Basic usage example
- [example/example2/](example/example2/) - Convert to all formats
- [README.md](README.md) - Main documentation
- [CHANGELOG.md](CHANGELOG.md) - Version history
