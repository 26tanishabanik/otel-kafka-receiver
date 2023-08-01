package customkafkareceiver 

import (
	"context"
	"time"

	"go.opencensus.io/stats/view"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/receiver"

	"github.com/open-telemetry/opentelemetry-collector-contrib/exporter/kafkaexporter"
	"github.com/26tanishabanik/otel-kafka-receiver/receiver/customkafkareceiver/internal/metadata"
)

const (
	defaultTopic         = "otlp_spans"
	defaultEncoding      = "otlp_proto"
	defaultBroker        = "localhost:9092"
	defaultClientID      = "otel-collector"
	defaultGroupID       = defaultClientID
	defaultInitialOffset = offsetLatest

	// default from sarama.NewConfig()
	defaultMetadataRetryMax = 3
	// default from sarama.NewConfig()
	defaultMetadataRetryBackoff = time.Millisecond * 250
	// default from sarama.NewConfig()
	defaultMetadataFull = true

	// default from sarama.NewConfig()
	defaultAutoCommitEnable = true
	// default from sarama.NewConfig()
	defaultAutoCommitInterval = 1 * time.Second
)

// FactoryOption applies changes to kafkaExporterFactory.
type FactoryOption func(factory *kafkaReceiverFactory)

// WithMetricsUnmarshalers adds MetricsUnmarshalers.
func WithMetricsUnmarshalers(metricsUnmarshalers ...MetricsUnmarshaler) FactoryOption {
	return func(factory *kafkaReceiverFactory) {
		for _, unmarshaler := range metricsUnmarshalers {
			factory.metricsUnmarshalers[unmarshaler.Encoding()] = unmarshaler
		}
	}
}

// NewFactory creates Kafka receiver factory.
func NewFactory(options ...FactoryOption) receiver.Factory {
	_ = view.Register(MetricViews()...)

	f := &kafkaReceiverFactory{
		metricsUnmarshalers: defaultMetricsUnmarshalers(),
	}
	for _, o := range options {
		o(f)
	}
	return receiver.NewFactory(
		metadata.Type,
		createDefaultConfig,
		receiver.WithMetrics(f.createMetricsReceiver, metadata.MetricsStability),
	)
}

func createDefaultConfig() component.Config {
	return &Config{
		Topic:         defaultTopic,
		Encoding:      defaultEncoding,
		Brokers:       []string{defaultBroker},
		ClientID:      defaultClientID,
		GroupID:       defaultGroupID,
		InitialOffset: defaultInitialOffset,
		Metadata: kafkaexporter.Metadata{
			Full: defaultMetadataFull,
			Retry: kafkaexporter.MetadataRetry{
				Max:     defaultMetadataRetryMax,
				Backoff: defaultMetadataRetryBackoff,
			},
		},
		AutoCommit: AutoCommit{
			Enable:   defaultAutoCommitEnable,
			Interval: defaultAutoCommitInterval,
		},
		MessageMarking: MessageMarking{
			After:   false,
			OnError: false,
		},
	}
}

type kafkaReceiverFactory struct {
	metricsUnmarshalers map[string]MetricsUnmarshaler
}

func (f *kafkaReceiverFactory) createMetricsReceiver(
	_ context.Context,
	set receiver.CreateSettings,
	cfg component.Config,
	nextConsumer consumer.Metrics,
) (receiver.Metrics, error) {
	c := cfg.(*Config)
	r, err := newMetricsReceiver(*c, set, f.metricsUnmarshalers, nextConsumer)
	if err != nil {
		return nil, err
	}
	return r, nil
}
