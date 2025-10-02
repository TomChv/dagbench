package hook

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

type SpanMarker struct {
	markers []string
	values  map[string]time.Duration
}

func NewSpanMarker(markers []string) *SpanMarker {
	return &SpanMarker{
		markers: markers,
		values:  make(map[string]time.Duration, len(markers)),
	}
}

func (s *SpanMarker) Values() map[string]time.Duration {
	return s.values
}

func (s *SpanMarker) Hook(line string) error {
	for _, marker := range s.markers {
		if strings.Contains(line, marker) && strings.Contains(line, "DONE") {
			duration, err := extractTimeFromTraceLine(line)
			if err != nil {
				return fmt.Errorf("failed to extract time from trace line: %w", err)
			}

			s.values[marker] = duration

			break
		}
	}

	return nil
}

// extractTimeFromTraceLine extracts the time from a trace line
// We assume the trace is with `--progress=plain`
func extractTimeFromTraceLine(line string) (time.Duration, error) {
	re := regexp.MustCompile(`\[(.*?)\]`)
	match := re.FindStringSubmatch(line)

	if len(match) > 1 {
		return time.ParseDuration(match[1])
	}

	return time.Duration(0), fmt.Errorf("failed to extract time from trace line: %s", line)
}
