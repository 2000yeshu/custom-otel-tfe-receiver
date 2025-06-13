package awscloudwatchmetricsreceiver

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	cloudwatchtypes "github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
)

func (sp *StatsProvider) getMetrics(ctx context.Context, cloudwatchMetrics *TFECloudwatchMetrics, nameConfig *NamedConfig) error {
	var dimensions = make([]cloudwatchtypes.DimensionFilter, 0)

	out, _ := sp.cw.ListMetrics(ctx, &cloudwatch.ListMetricsInput{
		Namespace: aws.String(nameConfig.Namespace),
		// MetricName: aws.String(metric),
		// Dimensions: []cloudwatchtypes.DimensionFilter{
		// 	{
		// 		Name:  aws.String("DBInstanceIdentifier"),
		// 		Value: aws.String("sandbox-tfe-rds-cluster-instance-0"),
		// 	},
		// },
	})

	for _, metric := range out.Metrics {
		for _, dim := range metric.Dimensions {
			// fmt.Println("Without metric Metric:", *metric.MetricName, "idx:", idx, "Dimension:", *dim.Name, "Value:", *dim.Value)
			dimensions = append(dimensions, cloudwatchtypes.DimensionFilter{
				Name:  dim.Name,
				Value: dim.Value,
			})
		}
	}
	// for _, metric := range nameConfig.MetricNames {
	// 	out, _ := sp.cw.ListMetrics(ctx, &cloudwatch.ListMetricsInput{
	// 		Namespace:  aws.String(nameConfig.Namespace),
	// 		MetricName: aws.String(metric),
	// 		// Dimensions: []cloudwatchtypes.DimensionFilter{
	// 		// 	{
	// 		// 		Name:  aws.String("DBInstanceIdentifier"),
	// 		// 		Value: aws.String("sandbox-tfe-rds-cluster-instance-0"),
	// 		// 	},
	// 		// },
	// 	})

	// 	for idx, metric := range out.Metrics {
	// 		for _, dim := range metric.Dimensions {
	// 			fmt.Println("Metric:", *metric.MetricName, "idx:", idx, "Dimension:", *dim.Name, "Value:", *dim.Value)
	// 			dimensions = append(dimensions, cloudwatchtypes.DimensionFilter{
	// 				Name:  dim.Name,
	// 				Value: dim.Value,
	// 			})
	// 		}
	// 	}
	// }

	var matchingDimensions = make([]cloudwatchtypes.DimensionFilter, 0)

	for _, dim := range dimensions {
		// Match the values with regex provided in nameConfig.Dimensions[0].Value
		ok, err := regexp.Match(nameConfig.Dimensions[0].Value, []byte(*dim.Value))
		if err != nil {
			return fmt.Errorf("error matching dimension value %s with regex %s: %v", *dim.Value, nameConfig.Dimensions[0].Value, err)
		}

		if ok && nameConfig.Dimensions[0].Name == *dim.Name {
			// if not present
			found := false
			for _, thisdim := range matchingDimensions {
				if *thisdim.Name == *dim.Name && *thisdim.Value == *dim.Value {
					found = true
					break
				}
			}
			if !found {
				matchingDimensions = append(matchingDimensions, dim)
			}
		}
	}

	fmt.Println("Matching Dimensions:", len(matchingDimensions))
	// Now for each matching dimension, create a MetricDataQuery
	// metricDataQueries := make([]cloudwatchtypes.MetricDataQuery, 0)
	for _, dim := range matchingDimensions {
		thisDimensionsBasedMetric := &DimensionBasedMetric{
			Key:          *dim.Name,
			Value:        *dim.Value,
			Prefix:       nameConfig.Prefix,
			MetricValues: make(map[string]cloudwatchtypes.MetricDataResult),
		}
		cloudwatchMetrics.DimensionBasedMetric = append(cloudwatchMetrics.DimensionBasedMetric, thisDimensionsBasedMetric)
		for _, metricName := range nameConfig.MetricNames {
			metricDataQuery := []cloudwatchtypes.MetricDataQuery{{
				Id: aws.String(rdsResource),
				MetricStat: &cloudwatchtypes.MetricStat{
					Metric: &cloudwatchtypes.Metric{
						Namespace:  aws.String(nameConfig.Namespace),
						MetricName: aws.String(metricName), // Assuming single metric for simplicity
						Dimensions: []cloudwatchtypes.Dimension{
							{
								Name:  dim.Name,
								Value: dim.Value,
							},
						},
					},
					Period: aws.Int32(int32(defaultPeriod.Seconds())),
					Stat:   aws.String(nameConfig.AwsAggregation),
				},
				ReturnData: aws.Bool(true),
			}}

			outputData, err := sp.cw.GetMetricData(ctx, &cloudwatch.GetMetricDataInput{
				MetricDataQueries: metricDataQuery,
				StartTime:         aws.Time(time.Now().Add(-time.Hour)), // 1 hour ago
				EndTime:           aws.Time(time.Now()),                 // Now
			})
			if err != nil {
				return err
			}
			thisDimensionsBasedMetric.MetricValues[metricName] = outputData.MetricDataResults[0]

		}
	}

	for _, metricDims := range cloudwatchMetrics.DimensionBasedMetric {
		fmt.Println("Metric:", metricDims.Key, "Value:", metricDims.Value)
		for metricName, metricData := range metricDims.MetricValues {
			fmt.Println("  Metric Name:", metricName)
			// fmt.Printf("    Value: %f at %s\n", value, metricData.Timestamps[idx])

			// Print no of values and timestamps
			fmt.Println("    No of Values:", len(metricData.Values))

		}
	}

	return nil

}
