package config

import (
	"github.com/zyl0501/go-push/tools/config"
	"github.com/zyl0501/go-push/tools/utils"
)

func GetHeartbeat(minHeartbeat int32, maxHeartbeat int32) int32 {
	return int32(utils.MaxInt(int(config.MinHeartbeat), utils.MaxInt(int(maxHeartbeat), int(config.MaxHeartbeat))))
}
