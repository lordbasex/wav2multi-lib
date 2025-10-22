//go:build cgo
// +build cgo

package wav2multi

/*
#cgo CFLAGS: -I/usr/local/include
#cgo LDFLAGS: -L/usr/local/lib -lbcg729
#include <bcg729/encoder.h>
#include <bcg729/decoder.h>
#include <stdlib.h>
*/
import "C"
import (
	"fmt"
	"io"
	"unsafe"
)

// G729Encoder implements G.729 encoding using libbcg729
type G729Encoder struct {
	encoder *C.bcg729EncoderChannelContextStruct
}

// NewG729Encoder creates a new G.729 encoder
func NewG729Encoder() (G729EncoderInterface, error) {
	encoder := C.initBcg729EncoderChannel(0) // 0 = disable VAD
	if encoder == nil {
		return nil, fmt.Errorf("failed to initialize G.729 encoder")
	}

	return &G729Encoder{
		encoder: encoder,
	}, nil
}

// Encode processes audio samples and writes G.729 encoded data
func (e *G729Encoder) Encode(samples []int16, writer io.Writer) error {
	if e.encoder == nil {
		return fmt.Errorf("encoder not initialized")
	}

	// Process samples in 80-sample frames (10ms at 8kHz)
	frameSize := 80
	for i := 0; i < len(samples); i += frameSize {
		// Get frame (pad with zeros if needed)
		frame := make([]int16, frameSize)
		copy(frame, samples[i:])

		// Convert to C array
		cFrame := (*C.int16_t)(unsafe.Pointer(&frame[0]))

		// Encode frame
		var output [10]C.uint8_t // G.729 produces 10 bytes per frame
		var bitStreamLength C.uint8_t
		C.bcg729Encoder(e.encoder, cFrame, (*C.uint8_t)(unsafe.Pointer(&output[0])), &bitStreamLength)

		// Write encoded data (use bitStreamLength for actual bytes written)
		if bitStreamLength > 0 {
			encodedData := (*[10]byte)(unsafe.Pointer(&output[0]))[:bitStreamLength]
			if _, err := writer.Write(encodedData); err != nil {
				return fmt.Errorf("failed to write G.729 data: %w", err)
			}
		}
	}

	return nil
}

// GetFormat returns the format this encoder handles
func (e *G729Encoder) GetFormat() AudioFormat {
	return FormatG729
}

// GetBitrate returns the bitrate in kbps
func (e *G729Encoder) GetBitrate() float64 {
	return 8.0 // 8 kbps
}

// Close releases the encoder resources
func (e *G729Encoder) Close() {
	if e.encoder != nil {
		C.closeBcg729EncoderChannel(e.encoder)
		e.encoder = nil
	}
}

// G729Decoder implements G.729 decoding using libbcg729
type G729Decoder struct {
	decoder *C.bcg729DecoderChannelContextStruct
}

// NewG729Decoder creates a new G.729 decoder
func NewG729Decoder() (*G729Decoder, error) {
	decoder := C.initBcg729DecoderChannel()
	if decoder == nil {
		return nil, fmt.Errorf("failed to initialize G.729 decoder")
	}

	return &G729Decoder{
		decoder: decoder,
	}, nil
}

// Decode processes G.729 encoded data and writes PCM samples
func (d *G729Decoder) Decode(reader io.Reader, writer io.Writer) error {
	if d.decoder == nil {
		return fmt.Errorf("decoder not initialized")
	}

	// Read and decode 10-byte frames
	buffer := make([]byte, 10)
	for {
		n, err := reader.Read(buffer)
		if err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("failed to read G.729 data: %w", err)
		}

		if n != 10 {
			return fmt.Errorf("incomplete G.729 frame: expected 10 bytes, got %d", n)
		}

		// Convert to C array
		cInput := (*C.uint8_t)(unsafe.Pointer(&buffer[0]))

		// Decode frame
		var output [80]C.int16_t // G.729 produces 80 samples per frame
		C.bcg729Decoder(d.decoder, cInput, 10, 0, 0, 0, (*C.int16_t)(unsafe.Pointer(&output[0])))

		// Write decoded PCM data
		decodedData := (*[160]byte)(unsafe.Pointer(&output[0]))[:] // 80 samples * 2 bytes
		if _, err := writer.Write(decodedData); err != nil {
			return fmt.Errorf("failed to write PCM data: %w", err)
		}
	}

	return nil
}

// Close releases the decoder resources
func (d *G729Decoder) Close() {
	if d.decoder != nil {
		C.closeBcg729DecoderChannel(d.decoder)
		d.decoder = nil
	}
}
