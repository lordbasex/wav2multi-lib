package wav2multi

import (
	"errors"
	"io"
)

// AudioFormat represents supported output formats
type AudioFormat string

const (
	FormatG729 AudioFormat = "g729"
	FormatULaw AudioFormat = "ulaw"
	FormatALaw AudioFormat = "alaw"
	FormatSLIN AudioFormat = "slin"
)

// TranscoderConfig holds configuration for the transcoder
type TranscoderConfig struct {
	// Input file path
	InputPath string
	// Output file path
	OutputPath string
	// Target format
	Format AudioFormat
}

// TranscoderResult holds the result of a transcoding operation
type TranscoderResult struct {
	// Input file information
	InputFile FileInfo
	// Output file information
	OutputFile FileInfo
	// Processing statistics
	Stats ProcessingStats
	// Any errors that occurred
	Error error
}

// FileInfo holds information about an audio file
type FileInfo struct {
	// File path
	Path string
	// File type (WAVE, etc.)
	Type string
	// Bit depth in bits
	BitDepth int
	// Sample rate in Hz
	SampleRate int
	// Number of channels
	Channels int
	// Total number of samples
	TotalSamples int
	// Duration in seconds
	Duration float64
	// File size in bytes
	Size int64
}

// ProcessingStats holds processing statistics
type ProcessingStats struct {
	// Processing time in milliseconds
	ProcessingTimeMs int64
	// Compression ratio (0.0 to 1.0)
	CompressionRatio float64
	// Bitrate in kbps
	BitrateKbps float64
	// Number of frames processed
	FramesProcessed int
}

// Transcoder interface defines the main transcoding functionality
type Transcoder interface {
	// Transcode converts audio from one format to another
	Transcode(config TranscoderConfig) (*TranscoderResult, error)
	// TranscodeFromReader converts audio from an io.Reader
	TranscodeFromReader(reader io.Reader, outputPath string, format AudioFormat) (*TranscoderResult, error)
	// TranscodeToWriter converts audio to an io.Writer
	TranscodeToWriter(inputPath string, writer io.Writer, format AudioFormat) (*TranscoderResult, error)
	// ValidateInput validates an input file
	ValidateInput(inputPath string) (*FileInfo, error)
	// GetSupportedFormats returns list of supported formats
	GetSupportedFormats() []AudioFormat
}

// CodecEncoder interface defines codec-specific encoding
type CodecEncoder interface {
	// Encode processes audio samples and writes encoded data
	Encode(samples []int16, writer io.Writer) error
	// GetFormat returns the format this encoder handles
	GetFormat() AudioFormat
	// GetBitrate returns the bitrate in kbps
	GetBitrate() float64
}

// Validation errors
var (
	ErrInvalidFormat     = errors.New("invalid audio format")
	ErrUnsupportedFormat = errors.New("unsupported format")
	ErrInvalidInput      = errors.New("invalid input file")
	ErrInvalidOutput     = errors.New("invalid output path")
	ErrCodecNotAvailable = errors.New("codec not available")
)

// Format validation
func IsValidFormat(format AudioFormat) bool {
	switch format {
	case FormatG729, FormatULaw, FormatALaw, FormatSLIN:
		return true
	default:
		return false
	}
}

// Get all supported formats
func GetSupportedFormats() []AudioFormat {
	return []AudioFormat{
		FormatG729,
		FormatULaw,
		FormatALaw,
		FormatSLIN,
	}
}
