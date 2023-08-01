package customkafkareceiver

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"

	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/pmetric"
)


type customMetricsUnmarshaler struct {
	pmetric.Unmarshaler
	encoding string
}


func (n customMetricsUnmarshaler) Encoding() string {
	return n.encoding
}

func (n customMetricsUnmarshaler) Unmarshal(data []byte) (pmetric.Metrics, error) {
	if len(data) == 0 {
		return pmetric.Metrics{}, fmt.Errorf("err: no %s data received", "custom metrics")
	}
	var format map[string]interface{}
	payload := bytes.NewReader(data)
	dec := json.NewDecoder(payload)
	err := dec.Decode(&format)
	if err != nil {
		return pmetric.Metrics{}, err
	}
	metrics := pmetric.NewMetrics()
	// OTLP metric doesn't identify entity type hence, Only one resource is created, hence resourceMetrics and
	// corresponding scopeMetrics is initialized outside loop
	resourceMetrics := metrics.ResourceMetrics().AppendEmpty()
	scopeMetrics := resourceMetrics.ScopeMetrics().AppendEmpty()
	// Keeps the track of recorded timestamps against given metric
	recordedTimestamps := make(map[string][]time.Duration, 0)
	t := format["timestamp"].(string)
	date, err := time.Parse(time.RFC3339, t)
	if err != nil {
		return pmetric.Metrics{}, err
	}
	unixMilli := float64(date.UnixMilli())
	startTimeNano := time.Duration(unixMilli) * time.Millisecond
	startTs := pcommon.Timestamp(startTimeNano)
	metricTs := pcommon.Timestamp(startTimeNano)
	resourceMetrics.Resource().Attributes().PutStr("resourceName", format["resourceName"].(string))

	for k, v := range format["values"].(map[string]interface{}) {
		metric := scopeMetrics.Metrics().AppendEmpty()
		metric.SetName(k)
		metric.SetDescription(k)
		recordedTimestamps[k] = append(recordedTimestamps[k], startTimeNano)
		metric.SetEmptyGauge()
		numberDataPoint := metric.Gauge().DataPoints().AppendEmpty()
		if _, ok := v.(map[string]interface{})["floatVal"]; ok {
			numberDataPoint.SetDoubleValue(v.(map[string]interface{})["floatVal"].(float64))
		} else if _, ok := v.(map[string]interface{})["intVal"]; ok {
			numberDataPoint.SetDoubleValue(v.(map[string]interface{})["intVal"].(float64))
		}
		numberDataPoint.SetTimestamp(metricTs)
		numberDataPoint.SetStartTimestamp(startTs)
		numberDataPoint.Attributes().PutStr("metricAttribute", format["metricAttribute"].(string))
		
	}
	return metrics, nil
}

func newCustomMetricsUnmarshaler(unmarshaler pmetric.Unmarshaler, encoding string) MetricsUnmarshaler {
	return customMetricsUnmarshaler{
		Unmarshaler: unmarshaler,
		encoding:    encoding,
	}
}
