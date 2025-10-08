package hook

import (
	"errors"
	"fmt"
)

var ErrExtractDurationFromTrace = errors.New("failed to extract time from trace line")

func ExtractDurationFromTraceError(line string) error {
	return fmt.Errorf("%w: %s", ErrExtractDurationFromTrace, line)
}
