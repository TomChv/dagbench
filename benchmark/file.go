package benchmark

import (
	"strings"

	"github.com/iancoleman/strcase"
)

type File struct {
	Path     string
	Metadata map[string]string
	Entries  []Entry
}

type Entry struct {
	Name        string
	DurationSec float64
}

func (b *File) Name() string {
	name, exists := b.Metadata["name"]
	if !exists {
		return strcase.ToCamel(strings.TrimSuffix(b.Path, ".txt"))
	}

	return strcase.ToCamel(name)
}

func (b *File) Version() string {
	version := b.FullVersion()

	// dagger v0.19.2-xxxx -> dagger v0.19.2
	version = strings.Split(version, "-")[0]

	// Remove dagger and space
	version = strings.TrimPrefix(version, "dagger ")

	return version
}

func (b *File) FullVersion() string {
	version, exist := b.Metadata["version"]
	if !exist {
		return ""
	}

	return version
}

func (b *File) Aggregate() map[string]float64 {
	agg := make(map[string][]float64)
	for _, r := range b.Entries {
		agg[r.Name] = append(agg[r.Name], r.DurationSec)
	}

	avg := make(map[string]float64)
	for name, times := range agg {
		sum := 0.0
		for _, t := range times {
			sum += t
		}
		avg[name] = sum / float64(len(times))
	}

	return avg
}
