package awscloudwatchmetricsreceiver

import "time"

type EC2Stats struct {
	CPUUtilized []float64   `json:"cpu_utilized"`
	Timestamps  []time.Time `json:"timestamps"`
}

type EC2Metadata struct {
	InstanceID string `json:"instance_id"`
}
