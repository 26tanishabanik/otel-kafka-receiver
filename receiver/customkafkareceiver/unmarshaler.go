package customkafkareceiver
import (
	"go.opentelemetry.io/collector/pdata/pmetric"
)

// MetricsUnmarshaler deserializes the message body
type MetricsUnmarshaler interface {
	// Unmarshal deserializes the message body into traces
	Unmarshal([]byte) (pmetric.Metrics, error)

	// Encoding of the serialized messages
	Encoding() string
}

func defaultMetricsUnmarshalers() map[string]MetricsUnmarshaler {
	custom := newCustomMetricsUnmarshaler(customMetricsUnmarshaler{}, "custom")
	return map[string]MetricsUnmarshaler{
		custom.Encoding(): custom,
	}
}
