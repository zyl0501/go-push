package utils

import (
	"net"
	log "github.com/alecthomas/log4go"
)

func ReadData(conn net.Conn, length uint32) ([]byte, error){
	readLen := 0
	date := make([]byte, length)
	for {
		n, err := conn.Read(date[readLen: length])
		if err != nil {
			return nil, err
		} else {
			log.Debug("data length %d", length)
			if uint32(readLen)+uint32(n) < length {
				log.Debug("read data part %s", string(date[readLen:readLen+n]))
				readLen += n
			} else {
				log.Debug("read data complete %s", string(date))
				return date, nil
			}
		}
	}
}