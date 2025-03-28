package smithyotelmetrics

import (
	"github.com/Enflick/smithy-go/metrics"
	otelmetric "go.opentelemetry.io/otel/metric"
)

func toInstrumentOpts(opts ...metrics.InstrumentOption) (unit, desc string) {
	var o metrics.InstrumentOptions
	for _, opt := range opts {
		opt(&o)
	}
	return o.UnitLabel, o.Description
}

func withMetricProps(opts ...metrics.RecordMetricOption) otelmetric.MeasurementOption {
	var o metrics.RecordMetricOptions
	for _, opt := range opts {
		opt(&o)
	}
	return otelmetric.WithAttributes(toOTELKeyValues(o.Properties)...)

}
