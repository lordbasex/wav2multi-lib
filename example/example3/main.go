// Copyright 2025 Federico Pereira <fpereira@cnsoluciones.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Example3: Advanced usage with io.Reader and io.Writer
//
// This example demonstrates advanced usage patterns:
// 1. TranscodeFromReader - Read from io.Reader (e.g., HTTP, memory, pipe)
// 2. TranscodeToWriter - Write to io.Writer (e.g., network, buffer, pipe)
//
// Usage:
//   go run main.go

package main

import (
	"bytes"
	"fmt"
	"log"
	"os"

	"github.com/lordbasex/wav2multi-lib"
)

func main() {
	fmt.Println("🎵 wav2multi-lib - Example 3: Advanced Usage")
	fmt.Println("=============================================")

	// Create a new transcoder with verbose logging disabled
	transcoder := wav2multi.NewTranscoder(false)

	// Example 1: TranscodeFromReader
	// Read from io.Reader (simulating HTTP, memory, etc.)
	fmt.Println("\n📖 Example 1: TranscodeFromReader")
	fmt.Println("   Reading WAV from io.Reader and saving to file")
	fmt.Println("   (Useful for HTTP uploads, memory buffers, etc.)")

	// Open input file as io.Reader
	inputFile, err := os.Open("input.wav")
	if err != nil {
		log.Fatalf("❌ Error opening input file: %v", err)
	}
	defer inputFile.Close()

	// Transcode from reader to file
	result1, err := transcoder.TranscodeFromReader(inputFile, "output-from-reader.alaw", wav2multi.FormatALaw)
	if err != nil {
		log.Fatalf("❌ Error transcoding from reader: %v", err)
	}

	fmt.Printf("   ✅ Success!\n")
	fmt.Printf("      Output: %s (%d bytes)\n", result1.OutputFile.Path, result1.OutputFile.Size)
	fmt.Printf("      Format: A-law @ %.1f kbps\n", result1.Stats.BitrateKbps)
	fmt.Printf("      Processing time: %d ms\n", result1.Stats.ProcessingTimeMs)

	// Example 2: TranscodeToWriter
	// Write to io.Writer (simulating network, buffer, etc.)
	fmt.Println("\n📝 Example 2: TranscodeToWriter")
	fmt.Println("   Reading WAV from file and writing to io.Writer")
	fmt.Println("   (Useful for HTTP responses, network streams, etc.)")

	// Create a buffer as our io.Writer (could be network connection, etc.)
	var buffer bytes.Buffer

	// Transcode from file to writer
	result2, err := transcoder.TranscodeToWriter("input.wav", &buffer, wav2multi.FormatSLIN)
	if err != nil {
		log.Fatalf("❌ Error transcoding to writer: %v", err)
	}

	fmt.Printf("   ✅ Success!\n")
	fmt.Printf("      Written to buffer: %d bytes\n", buffer.Len())
	fmt.Printf("      Format: SLIN @ %.1f kbps\n", result2.Stats.BitrateKbps)
	fmt.Printf("      Processing time: %d ms\n", result2.Stats.ProcessingTimeMs)

	// Optionally save the buffer to a file for verification
	err = os.WriteFile("output-to-writer.slin", buffer.Bytes(), 0644)
	if err != nil {
		log.Fatalf("❌ Error saving buffer to file: %v", err)
	}
	fmt.Printf("      Saved buffer to: output-to-writer.slin\n")

	// Example 3: Combining both (Reader to Writer)
	// This is useful for streaming scenarios
	fmt.Println("\n🔄 Example 3: Reader to Writer (Streaming)")
	fmt.Println("   Reading from io.Reader and writing to io.Writer")
	fmt.Println("   (Useful for streaming, proxying, etc.)")

	// Open input as reader
	inputFile2, err := os.Open("input.wav")
	if err != nil {
		log.Fatalf("❌ Error opening input file: %v", err)
	}
	defer inputFile2.Close()

	// Create output writer
	var streamBuffer bytes.Buffer

	// First, we need to read from reader to a temporary file
	// (since TranscodeToWriter expects a file path currently)
	tempData, err := os.ReadFile("input.wav")
	if err != nil {
		log.Fatalf("❌ Error reading input: %v", err)
	}

	// Write to temporary file
	tempFile := "temp-input.wav"
	err = os.WriteFile(tempFile, tempData, 0644)
	if err != nil {
		log.Fatalf("❌ Error writing temp file: %v", err)
	}
	defer os.Remove(tempFile)

	// Transcode to writer
	result3, err := transcoder.TranscodeToWriter(tempFile, &streamBuffer, wav2multi.FormatULaw)
	if err != nil {
		log.Fatalf("❌ Error in streaming transcode: %v", err)
	}

	fmt.Printf("   ✅ Success!\n")
	fmt.Printf("      Streamed: %d bytes\n", streamBuffer.Len())
	fmt.Printf("      Format: μ-law @ %.1f kbps\n", result3.Stats.BitrateKbps)
	fmt.Printf("      Processing time: %d ms\n", result3.Stats.ProcessingTimeMs)

	// Final summary
	fmt.Println("\n========================================")
	fmt.Println("✅ All advanced examples completed!")
	fmt.Println("========================================")
	fmt.Println("\n💡 Use Cases:")
	fmt.Println("   • FromReader: HTTP uploads, streaming input")
	fmt.Println("   • ToWriter: HTTP responses, network streams")
	fmt.Println("   • Combined: Real-time audio processing pipelines")
}
