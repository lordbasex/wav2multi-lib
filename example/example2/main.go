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

// Example2: Convert WAV to all supported formats
//
// This example demonstrates how to convert a single WAV file to all
// supported audio formats (G.729, μ-law, A-law, and SLIN) using a loop.
//
// Usage:
//   go run main.go
//
// Note: G.729 requires CGO and libbcg729. Without it, G.729 will be skipped.

package main

import (
	"fmt"
	"time"

	"github.com/lordbasex/wav2multi-lib"
)

func main() {
	fmt.Println("🎵 wav2multi-lib - Example 2: Convert to All Formats")
	fmt.Println("=====================================================")

	// Create a new transcoder with verbose logging disabled
	// We'll display our own summary instead
	transcoder := wav2multi.NewTranscoder(false)

	// Input file (same for all conversions)
	inputFile := "input.wav"

	// Define all codecs we want to convert to
	codecs := []string{"ulaw", "alaw", "g729", "slin"}

	// Track successful conversions
	successCount := 0
	failCount := 0

	// Convert to each format
	for i, codec := range codecs {
		fmt.Printf("📝 [%d/%d] Converting to %s...\n", i+1, len(codecs), codec)

		// Build output filename
		outputFile := fmt.Sprintf("output.%s", codec)

		// Configure transcoding
		config := wav2multi.TranscoderConfig{
			InputPath:  inputFile,
			OutputPath: outputFile,
			Format:     wav2multi.AudioFormat(codec),
		}

		// Perform transcoding
		result, err := transcoder.Transcode(config)
		if err != nil {
			fmt.Printf("   ❌ Failed: %v\n\n", err)
			failCount++

			// Wait 1 second before next conversion
			if i < len(codecs)-1 {
				time.Sleep(1 * time.Second)
			}
			continue
		}

		// Display results for this format
		fmt.Printf("   ✅ Success!\n")
		fmt.Printf("      Output: %s (%d bytes)\n", result.OutputFile.Path, result.OutputFile.Size)
		fmt.Printf("      Bitrate: %.1f kbps\n", result.Stats.BitrateKbps)
		fmt.Printf("      Processing time: %d ms\n", result.Stats.ProcessingTimeMs)
		fmt.Printf("      Compression: %.2f%%\n\n", result.Stats.CompressionRatio*100)
		successCount++

		// Wait 1 second before next conversion
		if i < len(codecs)-1 {
			time.Sleep(1 * time.Second)
		}
	}

	// Display final summary
	fmt.Println("========================================")
	fmt.Printf("✅ Completed: %d successful, %d failed\n", successCount, failCount)
	fmt.Println("========================================")

	if failCount > 0 {
		fmt.Println("\n⚠️  Note: G.729 requires CGO and libbcg729.")
		fmt.Println("   See CGO_SETUP.md for installation instructions.")
	}
}
