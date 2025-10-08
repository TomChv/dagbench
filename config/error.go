package config

import (
	"errors"
	"fmt"
)

var (
	ErrMissingSDKInInit     = errors.New("init.sdk is required")
	ErrNonPositiveIteration = errors.New("iteration must be greater than 0")
	ErrInitAndModuleOverlap = errors.New("cannot set both module and auto-init")
	ErrUnsuportedExtension  = errors.New("unsupported  file extension")
	ErrUnsupportedFormat    = errors.New("unsupported format")
)

func UnsupportedExtensionError(ext string) error {
	return fmt.Errorf("%w: %s", ErrUnsuportedExtension, ext)
}

func UnsupportedFormatError(format Format) error {
	return fmt.Errorf("%w: %s", ErrUnsupportedFormat, format)
}
