package config

import (
	"time"
	"github.com/go-akka/configuration"
	"fmt"
	"os"
	"log"
)

var CC = load()

type ccObj struct {
	Core struct {
		MinHeartbeat             time.Duration
		MaxHeartbeat             time.Duration
		MaxHeartbeatTimeoutTimes int
		SessionExpireTime        time.Duration
	}
	Security struct {
		PublicKey    string
		PrivateKey   string
		AesKeyLength int
	}
	Net struct {
		ConnectServerBindIp   string
		ConnectServerBindPort int
	}
}

func load() *ccObj {
	filename := "push.conf"
	_, err := os.Stat(filename)
	if err != nil {
		log.Fatalf("load config.yml error. %v", err)
	}

	cfg := configuration.LoadConfig(filename)
	CC := ccObj{}

	CC.Core.MinHeartbeat = cfg.GetTimeDuration("mp.core.min-heartbeat")
	CC.Core.MaxHeartbeat = cfg.GetTimeDuration("mp.core.max-heartbeat")
	CC.Core.MaxHeartbeatTimeoutTimes = int(cfg.GetInt32("mp.core.max-hb-timeout-times"))
	CC.Core.SessionExpireTime = cfg.GetTimeDuration("mp.core.session-expired-time")

	CC.Security.PublicKey = cfg.GetString("mp.security.public-key")
	CC.Security.PrivateKey = cfg.GetString("mp.security.private-key")
	CC.Security.AesKeyLength = int(cfg.GetInt32("mp.security.aes-key-length"))

	CC.Net.ConnectServerBindPort = int(cfg.GetInt32("mp.net.connect-server-port"))
	CC.Net.ConnectServerBindIp = cfg.GetString("mp.net.connect-server-bind-ip")

	fmt.Printf("config: %+v", CC)
	return &CC
}

var (
	MinHeartbeat             = 5 * time.Second
	MaxHeartbeat             = 5 * time.Second
	SessionExpireTime        = 10 * time.Second
	MaxHeartbeatTimeoutTimes = 3 //心跳检测，超时重试次数
	AesKeyLength             = 16
	PublicKey                = `
-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDZsfv1qscqYdy4vY+P4e3cAtmv
ppXQcRvrF1cB4drkv0haU24Y7m5qYtT52Kr539RdbKKdLAM6s20lWy7+5C0Dgacd
wYWd/7PeCELyEipZJL07Vro7Ate8Bfjya+wltGK9+XNUIHiumUKULW4KDx21+1NL
AUeJ6PeW+DAkmJWF6QIDAQAB
-----END PUBLIC KEY-----
`
	PrivateKey = `
-----BEGIN RSA PRIVATE KEY-----
MIICXQIBAAKBgQDZsfv1qscqYdy4vY+P4e3cAtmvppXQcRvrF1cB4drkv0haU24Y
7m5qYtT52Kr539RdbKKdLAM6s20lWy7+5C0DgacdwYWd/7PeCELyEipZJL07Vro7
Ate8Bfjya+wltGK9+XNUIHiumUKULW4KDx21+1NLAUeJ6PeW+DAkmJWF6QIDAQAB
AoGBAJlNxenTQj6OfCl9FMR2jlMJjtMrtQT9InQEE7m3m7bLHeC+MCJOhmNVBjaM
ZpthDORdxIZ6oCuOf6Z2+Dl35lntGFh5J7S34UP2BWzF1IyyQfySCNexGNHKT1G1
XKQtHmtc2gWWthEg+S6ciIyw2IGrrP2Rke81vYHExPrexf0hAkEA9Izb0MiYsMCB
/jemLJB0Lb3Y/B8xjGjQFFBQT7bmwBVjvZWZVpnMnXi9sWGdgUpxsCuAIROXjZ40
IRZ2C9EouwJBAOPjPvV8Sgw4vaseOqlJvSq/C/pIFx6RVznDGlc8bRg7SgTPpjHG
4G+M3mVgpCX1a/EU1mB+fhiJ2LAZ/pTtY6sCQGaW9NwIWu3DRIVGCSMm0mYh/3X9
DAcwLSJoctiODQ1Fq9rreDE5QfpJnaJdJfsIJNtX1F+L3YceeBXtW0Ynz2MCQBI8
9KP274Is5FkWkUFNKnuKUK4WKOuEXEO+LpR+vIhs7k6WQ8nGDd4/mujoJBr5mkrw
DPwqA3N5TMNDQVGv8gMCQQCaKGJgWYgvo3/milFfImbp+m7/Y3vCptarldXrYQWO
AQjxwc71ZGBFDITYvdgJM1MTqc8xQek1FXn1vfpy2c6O
-----END RSA PRIVATE KEY-----
`
)
