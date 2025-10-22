package wav2multi

import (
	"fmt"
	"io"
	"os"
	"time"
)

// DefaultTranscoder implements the Transcoder interface
type DefaultTranscoder struct {
	verbose bool
}

// NewTranscoder creates a new transcoder instance
func NewTranscoder(verbose bool) Transcoder {
	return &DefaultTranscoder{
		verbose: verbose,
	}
}

// Transcode converts audio from one format to another
func (t *DefaultTranscoder) Transcode(config TranscoderConfig) (*TranscoderResult, error) {
	startTime := time.Now()

	// Validate input
	if !IsValidFormat(config.Format) {
		return nil, ErrUnsupportedFormat
	}

	// Validate input file
	_, err := t.ValidateInput(config.InputPath)
	if err != nil {
		return nil, fmt.Errorf("input validation failed: %w", err)
	}

	// Create output file
	outputFile, err := os.Create(config.OutputPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create output file: %w", err)
	}
	defer func() { _ = outputFile.Close() }()

	// Get encoder for the target format
	encoder, err := GetEncoder(config.Format)
	if err != nil {
		return nil, fmt.Errorf("failed to get encoder: %w", err)
	}

	// Read input file
	inputFile, err := os.Open(config.InputPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open input file: %w", err)
	}
	defer func() { _ = inputFile.Close() }()

	// Read WAV samples
	samples, fileInfo, err := ReadWAVSamples(inputFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read WAV samples: %w", err)
	}

	// Encode samples
	if err := encoder.Encode(samples, outputFile); err != nil {
		return nil, fmt.Errorf("encoding failed: %w", err)
	}

	// Get output file info
	outputStat, err := os.Stat(config.OutputPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get output file info: %w", err)
	}

	// Calculate processing time
	processingTime := time.Since(startTime)

	// Calculate compression ratio
	compressionRatio := 0.0
	if fileInfo.Size > 0 {
		compressionRatio = float64(outputStat.Size()) / float64(fileInfo.Size)
	}

	// Create result
	result := &TranscoderResult{
		InputFile: *fileInfo,
		OutputFile: FileInfo{
			Path: config.OutputPath,
			Size: outputStat.Size(),
			Type: string(config.Format),
		},
		Stats: ProcessingStats{
			ProcessingTimeMs: processingTime.Milliseconds(),
			CompressionRatio: compressionRatio,
			BitrateKbps:      encoder.GetBitrate(),
			FramesProcessed:  len(samples),
		},
	}

	if t.verbose {
		t.logResult(result)
	}

	return result, nil
}

// TranscodeFromReader converts audio from an io.Reader
func (t *DefaultTranscoder) TranscodeFromReader(reader io.Reader, outputPath string, format AudioFormat) (*TranscoderResult, error) {
	startTime := time.Now()

	// Validate format
	if !IsValidFormat(format) {
		return nil, ErrUnsupportedFormat
	}

	// Get encoder
	encoder, err := GetEncoder(format)
	if err != nil {
		return nil, fmt.Errorf("failed to get encoder: %w", err)
	}

	// Read WAV samples from reader
	samples, fileInfo, err := ReadWAVSamples(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to read WAV samples: %w", err)
	}

	// Create output file
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create output file: %w", err)
	}
	defer func() { _ = outputFile.Close() }()

	// Encode samples
	if err := encoder.Encode(samples, outputFile); err != nil {
		return nil, fmt.Errorf("encoding failed: %w", err)
	}

	// Get output file info
	outputStat, err := os.Stat(outputPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get output file info: %w", err)
	}

	// Calculate processing time
	processingTime := time.Since(startTime)

	// Create result
	result := &TranscoderResult{
		InputFile: *fileInfo,
		OutputFile: FileInfo{
			Path: outputPath,
			Size: outputStat.Size(),
			Type: string(format),
		},
		Stats: ProcessingStats{
			ProcessingTimeMs: processingTime.Milliseconds(),
			BitrateKbps:      encoder.GetBitrate(),
			FramesProcessed:  len(samples),
		},
	}

	if t.verbose {
		t.logResult(result)
	}

	return result, nil
}

// TranscodeToWriter converts audio to an io.Writer
func (t *DefaultTranscoder) TranscodeToWriter(inputPath string, writer io.Writer, format AudioFormat) (*TranscoderResult, error) {
	startTime := time.Now()

	// Validate input
	if !IsValidFormat(format) {
		return nil, ErrUnsupportedFormat
	}

	// Validate input file
	_, err := t.ValidateInput(inputPath)
	if err != nil {
		return nil, fmt.Errorf("input validation failed: %w", err)
	}

	// Get encoder
	encoder, err := GetEncoder(format)
	if err != nil {
		return nil, fmt.Errorf("failed to get encoder: %w", err)
	}

	// Read input file
	inputFile, err := os.Open(inputPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open input file: %w", err)
	}
	defer func() { _ = inputFile.Close() }()

	// Read WAV samples
	samples, fileInfo, err := ReadWAVSamples(inputFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read WAV samples: %w", err)
	}

	// Encode samples to writer
	if err := encoder.Encode(samples, writer); err != nil {
		return nil, fmt.Errorf("encoding failed: %w", err)
	}

	// Calculate processing time
	processingTime := time.Since(startTime)

	// Create result
	result := &TranscoderResult{
		InputFile: *fileInfo,
		OutputFile: FileInfo{
			Type: string(format),
		},
		Stats: ProcessingStats{
			ProcessingTimeMs: processingTime.Milliseconds(),
			BitrateKbps:      encoder.GetBitrate(),
			FramesProcessed:  len(samples),
		},
	}

	if t.verbose {
		t.logResult(result)
	}

	return result, nil
}

// ValidateInput validates an input file
func (t *DefaultTranscoder) ValidateInput(inputPath string) (*FileInfo, error) {
	// Check if file exists
	stat, err := os.Stat(inputPath)
	if err != nil {
		return nil, fmt.Errorf("file not found: %w", err)
	}

	// Open file for analysis
	file, err := os.Open(inputPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer func() { _ = file.Close() }()

	// Read WAV samples to validate format
	_, fileInfo, err := ReadWAVSamples(file)
	if err != nil {
		return nil, fmt.Errorf("invalid WAV file: %w", err)
	}

	// Set file size
	fileInfo.Size = stat.Size()
	fileInfo.Path = inputPath

	return fileInfo, nil
}

// GetSupportedFormats returns list of supported formats
func (t *DefaultTranscoder) GetSupportedFormats() []AudioFormat {
	return GetSupportedFormats()
}

// logResult logs the transcoding result
func (t *DefaultTranscoder) logResult(result *TranscoderResult) {
	fmt.Printf("=== TRANSCODING RESULT ===\n")
	fmt.Printf("Input:  %s (%d bytes, %.2f seconds)\n",
		result.InputFile.Path, result.InputFile.Size, result.InputFile.Duration)
	fmt.Printf("Output: %s (%d bytes)\n",
		result.OutputFile.Path, result.OutputFile.Size)
	fmt.Printf("Format: %s (%.1f kbps)\n",
		result.OutputFile.Type, result.Stats.BitrateKbps)
	fmt.Printf("Processing: %d ms\n", result.Stats.ProcessingTimeMs)
	fmt.Printf("Compression: %.2f%%\n", result.Stats.CompressionRatio*100)
	fmt.Printf("Samples: %d\n", result.Stats.FramesProcessed)
	fmt.Printf("========================\n")
}
