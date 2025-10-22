package wav2multi

import (
	"bytes"
	"testing"
)

func TestPcmToULaw(t *testing.T) {
	// Test that the function produces consistent results
	tests := []struct {
		name  string
		input int16
	}{
		{"Zero", 0},
		{"Positive small", 100},
		{"Negative small", -100},
		{"Max positive", 32767},
		{"Max negative", -32767},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := pcmToULaw(tt.input)
			// Run twice to verify consistency
			result2 := pcmToULaw(tt.input)
			if result != result2 {
				t.Errorf("pcmToULaw(%d) inconsistent: %d != %d", tt.input, result, result2)
			}
		})
	}

	// Test that opposite signs produce different results
	pos := pcmToULaw(1000)
	neg := pcmToULaw(-1000)
	if pos == neg {
		t.Error("pcmToULaw should produce different results for opposite signs")
	}
}

func TestPcmToALaw(t *testing.T) {
	// Test that the function produces consistent results
	tests := []struct {
		name  string
		input int16
	}{
		{"Zero", 0},
		{"Positive small", 100},
		{"Negative small", -100},
		{"Max positive", 32767},
		{"Max negative", -32767},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := pcmToALaw(tt.input)
			// Run twice to verify consistency
			result2 := pcmToALaw(tt.input)
			if result != result2 {
				t.Errorf("pcmToALaw(%d) inconsistent: %d != %d", tt.input, result, result2)
			}
		})
	}

	// Test that opposite signs produce different results
	pos := pcmToALaw(1000)
	neg := pcmToALaw(-1000)
	if pos == neg {
		t.Error("pcmToALaw should produce different results for opposite signs")
	}
}

func TestULawEncoder(t *testing.T) {
	encoder := &ULawEncoder{}

	// Test GetFormat
	if encoder.GetFormat() != FormatULaw {
		t.Errorf("GetFormat() = %v, want %v", encoder.GetFormat(), FormatULaw)
	}

	// Test GetBitrate
	if encoder.GetBitrate() != 64.0 {
		t.Errorf("GetBitrate() = %v, want 64.0", encoder.GetBitrate())
	}

	// Test Encode
	samples := []int16{0, 100, -100, 1000, -1000}
	var buf bytes.Buffer

	err := encoder.Encode(samples, &buf)
	if err != nil {
		t.Fatalf("Encode() error = %v", err)
	}

	if buf.Len() != len(samples) {
		t.Errorf("Encode() produced %d bytes, want %d", buf.Len(), len(samples))
	}
}

func TestALawEncoder(t *testing.T) {
	encoder := &ALawEncoder{}

	// Test GetFormat
	if encoder.GetFormat() != FormatALaw {
		t.Errorf("GetFormat() = %v, want %v", encoder.GetFormat(), FormatALaw)
	}

	// Test GetBitrate
	if encoder.GetBitrate() != 64.0 {
		t.Errorf("GetBitrate() = %v, want 64.0", encoder.GetBitrate())
	}

	// Test Encode
	samples := []int16{0, 100, -100, 1000, -1000}
	var buf bytes.Buffer

	err := encoder.Encode(samples, &buf)
	if err != nil {
		t.Fatalf("Encode() error = %v", err)
	}

	if buf.Len() != len(samples) {
		t.Errorf("Encode() produced %d bytes, want %d", buf.Len(), len(samples))
	}
}

func TestSLINEncoder(t *testing.T) {
	encoder := &SLINEncoder{}

	// Test GetFormat
	if encoder.GetFormat() != FormatSLIN {
		t.Errorf("GetFormat() = %v, want %v", encoder.GetFormat(), FormatSLIN)
	}

	// Test GetBitrate
	if encoder.GetBitrate() != 128.0 {
		t.Errorf("GetBitrate() = %v, want 128.0", encoder.GetBitrate())
	}

	// Test Encode
	samples := []int16{0, 100, -100, 1000, -1000}
	var buf bytes.Buffer

	err := encoder.Encode(samples, &buf)
	if err != nil {
		t.Fatalf("Encode() error = %v", err)
	}

	// SLIN uses 2 bytes per sample (16-bit)
	expectedSize := len(samples) * 2
	if buf.Len() != expectedSize {
		t.Errorf("Encode() produced %d bytes, want %d", buf.Len(), expectedSize)
	}

	// Verify little-endian encoding
	data := buf.Bytes()
	// First sample is 0: should be [0x00, 0x00]
	if data[0] != 0x00 || data[1] != 0x00 {
		t.Errorf("First sample encoding incorrect: got [%02x %02x]", data[0], data[1])
	}

	// Second sample is 100 (0x0064): should be [0x64, 0x00] (little-endian)
	if data[2] != 0x64 || data[3] != 0x00 {
		t.Errorf("Second sample encoding incorrect: got [%02x %02x]", data[2], data[3])
	}
}

func TestGetEncoder(t *testing.T) {
	tests := []struct {
		name    string
		format  AudioFormat
		wantErr bool
	}{
		{"ULaw", FormatULaw, false},
		{"ALaw", FormatALaw, false},
		{"SLIN", FormatSLIN, false},
		{"Invalid", "invalid", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encoder, err := GetEncoder(tt.format)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetEncoder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && encoder == nil {
				t.Error("GetEncoder() returned nil encoder")
			}
			if !tt.wantErr && encoder.GetFormat() != tt.format {
				t.Errorf("GetEncoder() format = %v, want %v", encoder.GetFormat(), tt.format)
			}
		})
	}
}

func TestIsValidFormat(t *testing.T) {
	tests := []struct {
		name   string
		format AudioFormat
		want   bool
	}{
		{"G729", FormatG729, true},
		{"ULaw", FormatULaw, true},
		{"ALaw", FormatALaw, true},
		{"SLIN", FormatSLIN, true},
		{"Invalid", "mp3", false},
		{"Empty", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidFormat(tt.format); got != tt.want {
				t.Errorf("IsValidFormat() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetSupportedFormats(t *testing.T) {
	formats := GetSupportedFormats()

	if len(formats) != 4 {
		t.Errorf("GetSupportedFormats() returned %d formats, want 4", len(formats))
	}

	// Verify all expected formats are present
	expectedFormats := map[AudioFormat]bool{
		FormatG729: false,
		FormatULaw: false,
		FormatALaw: false,
		FormatSLIN: false,
	}

	for _, format := range formats {
		if _, exists := expectedFormats[format]; exists {
			expectedFormats[format] = true
		}
	}

	for format, found := range expectedFormats {
		if !found {
			t.Errorf("GetSupportedFormats() missing format %v", format)
		}
	}
}
