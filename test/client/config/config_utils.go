package config

import (
	"github.com/zyl0501/go-push/tools/config"
	"github.com/zyl0501/go-push/tools/utils"
	"time"
)

func GetHeartbeat(expireHeartbeat time.Duration) time.Duration {
	min := config.CC.Core.MinHeartbeat
	max := config.CC.Core.MaxHeartbeat
	return time.Duration(utils.MaxInt64(int64(min), utils.MinInt64(int64(expireHeartbeat), int64(max))))
}
