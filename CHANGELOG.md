# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## About This Project

This library is a refactored version of the original [wav2multi](https://github.com/lordbasex/wav2multi) CLI application. **Both projects were created by Federico Pereira**, who refactored his own CLI tool into this reusable Go library with professional features and extensive documentation.

## [1.0.0] - 2025-10-22

### Added
- Initial release of wav2multi-lib
- Support for G.729 codec (with CGO)
- Support for μ-law (ulaw) codec
- Support for A-law (alaw) codec
- Support for SLIN (raw PCM) format
- Transcoder interface with multiple methods
- Support for io.Reader and io.Writer
- Input validation for WAV files
- Processing statistics
- Comprehensive error handling
- Two example applications (basic and multi-format)
- Comprehensive unit tests (codecs_test.go)
- CI/CD with GitHub Actions
- Makefile for development commands
- Apache 2.0 License
- Professional documentation (README, CGO_SETUP, examples)

### Changed
- Removed redundant `Verbose` field from `TranscoderConfig`
- Verbose logging now controlled only by `NewTranscoder(verbose bool)`
- Reorganized examples into separate directories (example1 and example2)
- Updated README with badges and improved structure
- **IMPORTANT:** Corrected module path from `github.com/cnsoluciones/wav2multi-lib` to `github.com/lordbasex/wav2multi-lib`

### Fixed
- Updated G.729 codec to support bcg729 v1.1.1 API
- Added missing `bitStreamLength` parameter in encoder
- Added missing parameters in decoder

### Examples
- **example1/**: Basic usage - single format conversion
- **example2/**: Advanced usage - convert to all formats with loop

## [Unreleased]

### Planned
- Streaming support for large files
- Additional codec support
- Performance optimizations
- Extended documentation
- Benchmarks

