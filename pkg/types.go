package root

import (
	"math/big"
	"os"
)

type EnvConfig struct {
	NodeEndpoint       string `json:"nodeEndpoint"`
	Bucket             string `json:"bucket"`
	Region             string `json:"region"`
	BlockEstimateCount string `json:"blockEstimateCount"`
}

type CachedBlockItem struct {
	BlockNum big.Int `json:"blockNum"`
	MinGas   big.Int `json:"minGas"`
}

type CachedItem struct {
	LastBlock int64             `json:"lastBlock"`
	Blocks    []CachedBlockItem `json:"blocks"`
}

type GasExpress struct {
	SafeLow  float64 `json:"safeLow"`
	Standard float64 `json:"standard"`
	Fast     float64 `json:"fast"`
	Fastest  float64 `json:"fastest"`
	BlockNum int64   `json:"blockNum"`
}

var Config = EnvConfig{
	Bucket:             EnvOrDefaultString("bucket", "gas-express"),
	Region:             EnvOrDefaultString("region", "us-east-1"),
	NodeEndpoint:       EnvOrDefaultString("nodeEndpoint", "http://localhost:8080"),
	BlockEstimateCount: EnvOrDefaultString("blockEstimateCount", "200"),
}

func EnvOrDefaultString(envVar string, defaultValue string) string {
	value := os.Getenv(envVar)
	if value == "" {
		return defaultValue
	}

	return value
}
