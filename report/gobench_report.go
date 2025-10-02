package report

import (
	"fmt"
	"quartz/dagbench.io/config"
	"time"

	"strings"

	"text/tabwriter"

	"github.com/iancoleman/strcase"
)

func FormatGoBenchReport(conf *config.Config, reports []*Report) string {
	var content strings.Builder

	metadatas := getBenchMetadatas(conf)
	_, _ = fmt.Fprintf(&content, "%s\n\n", strings.Join(metadatas, "\n"))

	tw := tabwriter.NewWriter(&content, 0, 0, 2, ' ', 0)
	for _, report := range reports {
		for key, value := range report.values {
			msTime := fmt.Sprintf("%.1f s/op", value.Seconds())
			benKey := fmt.Sprintf("%s/%s", strcase.ToCamel(report.name), strcase.ToLowerCamel(key))
			iteration := 1

			_, _ = fmt.Fprintf(tw, "Benchmark%s\t%d\t%s\n", benKey, iteration, msTime)
		}
	}

	_ = tw.Flush()

	return content.String()
}

func getBenchMetadatas(conf *config.Config) []string {
	entries := map[string]string{
		"name":    conf.Name,
		"version": conf.Version,
		"date":    time.Now().Format("2006-01-02"),
		"sdk":     conf.SDK,
		"cloud":   fmt.Sprintf("%t", conf.Cloud),
	}

	var metadatas []string

	for key, value := range entries {
		metadatas = append(metadatas, fmt.Sprintf("# %s: %s", key, value))
	}

	return metadatas
}
