package config

import (
	"github.com/zyl0501/go-push/tools/config"
	"github.com/zyl0501/go-push/tools/utils"
	"time"
)

func GetHeartbeat(minHeartbeat time.Duration, maxHeartbeat time.Duration) time.Duration {
	return time.Duration(utils.MaxInt64(int64(config.MinHeartbeat), utils.MinInt64(int64(maxHeartbeat), int64(config.MaxHeartbeat))))
}
