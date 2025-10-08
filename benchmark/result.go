package benchmark

import (
	"fmt"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/iancoleman/strcase"
)

type Result []map[string]time.Duration

type Metadata map[string]string

// GoBenchFormat returns the result in the gobench format.
func (r *Result) GoBenchFormat(name string, metadatas Metadata) string {
	var content strings.Builder

	for key, value := range metadatas {
		_, _ = fmt.Fprintf(&content, "# %s: %s\n", key, value)
	}

	_, _ = fmt.Fprintf(&content, "\n\n")

	tw := tabwriter.NewWriter(&content, 0, 0, 2, ' ', 0)
	for _, iterations := range *r {
		for key, value := range iterations {
			msTime := fmt.Sprintf("%.1f s/op", value.Seconds())
			benKey := fmt.Sprintf("%s/%s", strcase.ToCamel(name), key)
			iteration := 1

			_, _ = fmt.Fprintf(tw, "Benchmark%s\t%d\t%s\n", benKey, iteration, msTime)
		}
	}

	_ = tw.Flush()

	return content.String()
}
