package wav2multi

import (
	"fmt"
	"io"
	"os"

	youpywav "github.com/youpy/go-wav"
)

// G729Encoder interface for G.729 encoding
type G729EncoderInterface interface {
	CodecEncoder
	Close()
}

// ULawEncoder implements μ-law encoding
type ULawEncoder struct{}

func (e *ULawEncoder) Encode(samples []int16, writer io.Writer) error {
	for _, sample := range samples {
		ulawByte := pcmToULaw(sample)
		if _, err := writer.Write([]byte{ulawByte}); err != nil {
			return err
		}
	}
	return nil
}

func (e *ULawEncoder) GetFormat() AudioFormat {
	return FormatULaw
}

func (e *ULawEncoder) GetBitrate() float64 {
	return 64.0 // 64 kbps
}

// ALawEncoder implements A-law encoding
type ALawEncoder struct{}

func (e *ALawEncoder) Encode(samples []int16, writer io.Writer) error {
	for _, sample := range samples {
		alawByte := pcmToALaw(sample)
		if _, err := writer.Write([]byte{alawByte}); err != nil {
			return err
		}
	}
	return nil
}

func (e *ALawEncoder) GetFormat() AudioFormat {
	return FormatALaw
}

func (e *ALawEncoder) GetBitrate() float64 {
	return 64.0 // 64 kbps
}

// SLINEncoder implements SLIN (PCM 16-bit) encoding
type SLINEncoder struct{}

func (e *SLINEncoder) Encode(samples []int16, writer io.Writer) error {
	for _, sample := range samples {
		// Write 16-bit PCM in little-endian format
		bytes := []byte{
			byte(sample & 0xFF),        // Low byte
			byte((sample >> 8) & 0xFF), // High byte
		}
		if _, err := writer.Write(bytes); err != nil {
			return err
		}
	}
	return nil
}

func (e *SLINEncoder) GetFormat() AudioFormat {
	return FormatSLIN
}

func (e *SLINEncoder) GetBitrate() float64 {
	return 128.0 // 128 kbps
}

// pcmToULaw converts 16-bit PCM to μ-law
func pcmToULaw(pcm int16) byte {
	// Get sign and magnitude
	sign := pcm < 0
	if sign {
		pcm = -pcm
	}

	// Note: pcm is already constrained to int16 range (-32768 to 32767)
	// After taking absolute value, it's in range (0 to 32767)

	// Add bias
	pcm += 33

	// Find segment
	segment := 0
	temp := pcm >> 7
	for temp > 0 {
		segment++
		temp >>= 1
	}
	if segment > 7 {
		segment = 7
	}

	// Quantize
	quantization := (pcm >> (segment + 3)) & 0x0F

	// Build μ-law byte
	ulaw := byte(segment<<4) | byte(quantization)
	if sign {
		ulaw |= 0x80
	}

	return ^ulaw // Invert all bits
}

// pcmToALaw converts 16-bit PCM to A-law
func pcmToALaw(pcm int16) byte {
	// Get sign and magnitude
	sign := pcm < 0
	if sign {
		pcm = -pcm
	}

	// Note: pcm is already constrained to int16 range (-32768 to 32767)
	// After taking absolute value, it's in range (0 to 32767)

	// Add bias
	pcm += 33

	// Find segment
	segment := 0
	temp := pcm >> 7
	for temp > 0 {
		segment++
		temp >>= 1
	}
	if segment > 7 {
		segment = 7
	}

	// Quantize
	quantization := (pcm >> (segment + 3)) & 0x0F

	// Build A-law byte
	alaw := byte(segment<<4) | byte(quantization)
	if sign {
		alaw |= 0x80
	}

	// A-law uses even bits for segment
	alaw ^= 0x55 // XOR with 0x55 to get even bits

	return alaw
}

// GetEncoder returns the appropriate encoder for the given format
func GetEncoder(format AudioFormat) (CodecEncoder, error) {
	switch format {
	case FormatG729:
		encoder, err := NewG729Encoder()
		if err != nil {
			return nil, fmt.Errorf("G.729 encoder not available: %w", err)
		}
		return encoder, nil
	case FormatULaw:
		return &ULawEncoder{}, nil
	case FormatALaw:
		return &ALawEncoder{}, nil
	case FormatSLIN:
		return &SLINEncoder{}, nil
	default:
		return nil, ErrUnsupportedFormat
	}
}

// ReadWAVSamples reads samples from a WAV file using youpy/go-wav
func ReadWAVSamples(reader io.Reader) ([]int16, *FileInfo, error) {
	// Convert io.Reader to a file-like interface
	// For now, we'll use a simplified approach
	file, ok := reader.(*os.File)
	if !ok {
		return nil, nil, fmt.Errorf("reader must be *os.File for youpy/go-wav")
	}

	wavReader := youpywav.NewReader(file)

	// Get format information
	format, err := wavReader.Format()
	if err != nil {
		return nil, nil, err
	}

	// Validate format
	if format.AudioFormat != 1 {
		return nil, nil, ErrInvalidFormat
	}
	if format.NumChannels != 1 {
		return nil, nil, ErrInvalidFormat
	}
	if format.SampleRate != 8000 {
		return nil, nil, ErrInvalidFormat
	}
	if format.BitsPerSample != 16 {
		return nil, nil, ErrInvalidFormat
	}

	// Read all samples
	var samples []int16
	for {
		sampleBatch, err := wavReader.ReadSamples(1024)
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, nil, err
		}

		for _, s := range sampleBatch {
			samples = append(samples, int16(s.Values[0]))
		}
	}

	// Create file info
	fileInfo := &FileInfo{
		Type:         "WAVE",
		BitDepth:     int(format.BitsPerSample),
		SampleRate:   int(format.SampleRate),
		Channels:     int(format.NumChannels),
		TotalSamples: len(samples),
		Duration:     float64(len(samples)) / float64(format.SampleRate),
	}

	return samples, fileInfo, nil
}

// AnalyzeWAVFile analyzes a WAV file and returns detailed information
func AnalyzeWAVFile(inputPath string) (*FileInfo, error) {
	// This would use go-audio/wav for detailed analysis
	// For now, we'll return a basic implementation
	return &FileInfo{
		Path: inputPath,
		Type: "WAVE",
		// Other fields would be populated by actual analysis
	}, nil
}
