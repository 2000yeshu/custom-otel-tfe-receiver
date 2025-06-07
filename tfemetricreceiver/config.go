// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

// TODO: Add more validations related to the  namespaces
// and metrics for that namespaces and dimensions for those metrics allowed

package awscloudwatchmetricsreceiver // import "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/awscloudwatchmetricsreceiver"

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"
)

var defaultPollInterval = 5 * time.Minute
var defaultPeriod = 60 * time.Second

// Config is the overall config structure for the awscloudwatchmetricsreceiver
type Config struct {
	Region       string         `mapstructure:"region"`
	Profile      string         `mapstructure:"profile"`
	IMDSEndpoint string         `mapstructure:"imds_endpoint"`
	PollInterval time.Duration  `mapstructure:"poll_interval"`
	Metrics      *MetricsConfig `mapstructure:"metrics"`

	// prevent unkeyed literal initialization
	_ struct{}
}

// MetricsConfig is the configuration for the metrics part of the receiver
// added this so we could expand to other inputs such as autodiscover
type MetricsConfig struct {
	Names []*NamedConfig `mapstructure:"named"`

	// prevent unkeyed literal initialization
	_ struct{}
}

// NamesConfig is the configuration for the metric namespace and metric names
// https://docs.aws.amazon.com/AmazonCloudWatch/latest/monitoring/cloudwatch_concepts.html
type NamedConfig struct {
	// Namespace  string `mapstructure:"namespace"`
	Type        string   `mapstructure:"type"`
	MetricNames []string `mapstructure:"metric_names"`
	ARN         string   `mapstructure:"arn"`

	// Defaults to 60 seconds
	Period         time.Duration `mapstructure:"period"`
	AwsAggregation string        `mapstructure:"aws_aggregation"`
	// Dimensions     []MetricDimensionsConfig `mapstructure:"dimensions"`
}

// MetricDimensionConfig is the configuration for the metric dimensions
// type MetricDimensionsConfig struct {
// 	Name  string `mapstructure:"Name"`
// 	Value string `mapstructure:"Value"`

// 	// prevent unkeyed literal initialization
// 	_ struct{}
// }

var (
	errNoMetricsConfigured   = errors.New("no metrics configured")
	errNoNamespaceConfigured = errors.New("no metric namespace configured")
	errNoRegion              = errors.New("no region was specified")
	errInvalidPollInterval   = errors.New("poll interval is incorrect, it must be a duration greater than one second")

	// https://docs.aws.amazon.com/cli/latest/reference/cloudwatch/get-metric-data.html
	errEmptyDimensions      = errors.New("dimensions name and value is empty")
	errTooManyDimensions    = errors.New("you cannot define more than 30 dimensions for a metric")
	errInvalidPeriod        = errors.New("period must be between 60 and 21600 seconds, inclusive")
	errDimensionColonPrefix = errors.New("dimension name cannot start with a colon")

	errInvalidAwsAggregation = errors.New("invalid AWS aggregation")
)

func (cfg *Config) Validate() error {
	if cfg.Region == "" {
		return errNoRegion
	}

	if cfg.IMDSEndpoint != "" {
		_, err := url.ParseRequestURI(cfg.IMDSEndpoint)
		if err != nil {
			return fmt.Errorf("unable to parse URI for imds_endpoint: %w", err)
		}
	}

	if cfg.PollInterval < time.Second {
		return errInvalidPollInterval
	}
	var errs error
	errs = errors.Join(errs, cfg.validateMetricsConfig())
	return errs
}

func (cfg *Config) validateMetricsConfig() error {
	if cfg.Metrics == nil {
		return errNoMetricsConfigured
	}
	return cfg.validateNamedConfig()
}

func (cfg *Config) validateNamedConfig() error {
	if cfg.Metrics.Names == nil {
		return errNoMetricsConfigured
	}
	return cfg.validateDimensionsConfig()
}

func (cfg *Config) validateDimensionsConfig() error {
	var errs error

	metricsNames := cfg.Metrics.Names
	for _, name := range metricsNames {
		if name.Type == "" {
			return errNoNamespaceConfigured
		}
		err := validateAwsAggregation(name.AwsAggregation)
		if err != nil {
			return err
		}
		if len(name.MetricNames) == 0 {
			return errNoMetricsConfigured
		}
	}
	return errs
}

// https://docs.aws.amazon.com/AmazonCloudWatch/latest/monitoring/Statistics-definitions.html
func validateAwsAggregation(agg string) error {
	switch {
	case agg == "SampleCount":
		return nil
	case agg == "Sum":
		return nil
	case agg == "Average":
		return nil
	case agg == "Minimum":
		return nil
	case agg == "Maximum":
		return nil
	case strings.HasPrefix(agg, "p"):
		return nil
	case strings.HasPrefix(agg, "TM"):
		return nil
	case agg == "IQM":
		return nil
	case strings.HasPrefix(agg, "PR"):
		return nil
	case strings.HasPrefix(agg, "TC"):
		return nil
	case strings.HasPrefix(agg, "TS"):
		return nil
	default:
		return errInvalidAwsAggregation
	}
}
