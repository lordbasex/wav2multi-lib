//go:build !cgo
// +build !cgo

package wav2multi

import (
	"fmt"
	"io"
)

// G729EncoderNoCGO implements G.729 encoding (CGO disabled)
type G729EncoderNoCGO struct{}

// NewG729Encoder creates a new G.729 encoder (CGO disabled)
func NewG729Encoder() (G729EncoderInterface, error) {
	return nil, fmt.Errorf("G.729 encoding requires CGO and libbcg729 library")
}

// Encode processes audio samples and writes G.729 encoded data (CGO disabled)
func (e *G729EncoderNoCGO) Encode(samples []int16, writer io.Writer) error {
	return fmt.Errorf("G.729 encoding requires CGO and libbcg729 library")
}

// GetFormat returns the format this encoder handles
func (e *G729EncoderNoCGO) GetFormat() AudioFormat {
	return FormatG729
}

// GetBitrate returns the bitrate in kbps
func (e *G729EncoderNoCGO) GetBitrate() float64 {
	return 8.0 // 8 kbps
}

// Close releases the encoder resources
func (e *G729EncoderNoCGO) Close() {
	// No-op for non-CGO version
}

// G729Decoder implements G.729 decoding (CGO disabled)
type G729Decoder struct{}

// NewG729Decoder creates a new G.729 decoder (CGO disabled)
func NewG729Decoder() (*G729Decoder, error) {
	return nil, fmt.Errorf("G.729 decoding requires CGO and libbcg729 library")
}

// Decode processes G.729 encoded data and writes PCM samples (CGO disabled)
func (d *G729Decoder) Decode(reader io.Reader, writer io.Writer) error {
	return fmt.Errorf("G.729 decoding requires CGO and libbcg729 library")
}

// Close releases the decoder resources
func (d *G729Decoder) Close() {
	// No-op for non-CGO version
}
