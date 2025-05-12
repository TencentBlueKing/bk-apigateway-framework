package probe

import "context"

// Result 健康探针结果
type Result struct {
	Name     string `json:"name"`
	Core     bool   `json:"core"`
	Healthy  bool   `json:"healthy"`
	Endpoint string `json:"endpoint"`
	Issue    string `json:"issue"`
}

// HealthProbe 健康探针
type HealthProbe interface {
	// Perform 执行探针
	Perform(ctx context.Context) *Result
}
