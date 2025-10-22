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

// Example1: Basic usage of wav2multi-lib
//
// This example demonstrates the simplest way to convert a WAV file
// to a single audio format (μ-law in this case).
//
// Usage:
//   go run main.go

package main

import (
	"fmt"
	"log"

	"github.com/lordbasex/wav2multi-lib"
)

func main() {
	fmt.Println("🎵 wav2multi-lib - Example 1: Basic Usage")
	fmt.Println("==========================================")

	// Create a new transcoder with verbose logging enabled
	// verbose = true shows detailed processing information
	transcoder := wav2multi.NewTranscoder(true)

	// Configure the transcoding operation
	// This is where you specify:
	// - InputPath: path to your WAV file
	// - OutputPath: where to save the converted file (can include directory)
	// - Format: which codec to use (ulaw, alaw, slin, or g729)
	config := wav2multi.TranscoderConfig{
		InputPath:  "input.wav",
		OutputPath: "output.ulaw",
		Format:     wav2multi.FormatULaw,
	}

	// Perform the transcoding
	fmt.Println("📝 Converting to μ-law format...")
	result, err := transcoder.Transcode(config)
	if err != nil {
		log.Fatalf("❌ Error during transcoding: %v", err)
	}

	// Display the results
	fmt.Println("\n✅ Conversion completed successfully!")
	fmt.Printf("   Input:  %s\n", result.InputFile.Path)
	fmt.Printf("   Output: %s (%d bytes)\n", result.OutputFile.Path, result.OutputFile.Size)
	fmt.Printf("   Format: %s @ %.1f kbps\n", result.OutputFile.Type, result.Stats.BitrateKbps)
	fmt.Printf("   Processing time: %d ms\n", result.Stats.ProcessingTimeMs)
	fmt.Printf("   Compression ratio: %.2f%%\n", result.Stats.CompressionRatio*100)
}
