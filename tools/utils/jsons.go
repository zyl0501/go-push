package utils

import (
	"encoding/json"
	log "github.com/alecthomas/log4go"
)

func ToJson(v interface{}) []byte {
	data, err := json.Marshal(v)
	if err != nil {
		log.Error("Jsons.toJson ex, bean=%v, err=%v", v, err);
		return nil
	} else {
		return data
	}
}

func FromJson(data []byte, v interface{}) {
	err := json.Unmarshal(data, v)
	if err != nil {
		log.Error("Jsons.fromJson ex, json="+string(data)+", v=%v, err=%v", v, err);
	}
}
