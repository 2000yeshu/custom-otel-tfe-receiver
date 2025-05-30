package awscloudwatchmetricsreceiver

import (
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/pmetric"
	conventions "go.opentelemetry.io/otel/semconv/v1.21.0"
)

func convertToOTLPMetrics(prefix string, m CloudWatchMetrics, r pcommon.Resource, timestamp pcommon.Timestamp) pmetric.Metrics {
	md := pmetric.NewMetrics()
	rm := md.ResourceMetrics().AppendEmpty()
	rm.SetSchemaUrl(conventions.SchemaURL)
	r.CopyTo(rm.Resource())

	ilms := rm.ScopeMetrics()

	appendDoubleGauge(prefix+attributeCPUUtilized, unitNone, m.CPUUtilized, timestamp, ilms.AppendEmpty())

	return md
}

func appendDoubleGauge(metricName, unit string, value float64, ts pcommon.Timestamp, ilm pmetric.ScopeMetrics) {
	metric := appendMetric(ilm, metricName, unit)
	gauge := metric.SetEmptyGauge()
	dp := gauge.DataPoints().AppendEmpty()
	dp.SetDoubleValue(value)
	dp.SetTimestamp(ts)
}

// Append a metric and set name and unit
func appendMetric(ilm pmetric.ScopeMetrics, name, unit string) pmetric.Metric {
	metric := ilm.Metrics().AppendEmpty()
	metric.SetName(name)
	metric.SetUnit(unit)

	return metric
}
