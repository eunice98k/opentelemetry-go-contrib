package cortex_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"go.opentelemetry.io/otel/sdk/export/metric"
	"go.opentelemetry.io/otel/sdk/export/metric/aggregation"
	"opentelemetry.io/contrib/exporters/metric/cortex"
)

func TestExportKindFor(t *testing.T) {
	exporter := cortex.Exporter{}
	got := exporter.ExportKindFor(nil, aggregation.Kind(0))
	want := metric.CumulativeExporter

	if got != want {
		t.Errorf("ExportKindFor() =  %q, want %q", got, want)
	}
}

// TestNewRawExporter tests whether NewRawExporter successfully creates an Exporter with the same
// Config struct as the one passed in.
func TestNewRawExporter(t *testing.T) {
	exporter, err := cortex.NewRawExporter(ValidConfig)
	if err != nil {
		t.Errorf("Failed to create exporter with error %v", err)
	}

	if !cmp.Equal(ValidConfig, exporter.Config) {
		t.Errorf("Got configuration %v, wanted %v", exporter.Config, ValidConfig)
	}
}

// TestNewExportPipeline tests whether a push Controller was successfully created with an Exporter
// from New RawExporter. Errors in this function will be from calls to push controller package and
// NewRawExport. Both have their own tests.
func TestNewExportPipeline(t *testing.T) {
	_, err := cortex.NewExportPipeline(ValidConfig)
	if err != nil {
		t.Errorf("Failed to create export pipeline with error %v", err)
	}
}
