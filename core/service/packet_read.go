package service

import (
	"github.com/zyl0501/go-push/tools/utils"
	"github.com/zyl0501/go-push/api/protocol"
	"net"
)

func ReadPacket(conn net.Conn) (*protocol.Packet, error){
	header, err := utils.ReadData(conn, uint32(protocol.HeadLength))
	if err != nil {
		return nil, err
	}
	packet, bodyLen := protocol.DecodePacket(header)

	if bodyLen > 0 {
		body, err := utils.ReadData(conn, bodyLen)
		if err != nil {
			return nil, err
		}
		packet.Body = body
	}
	return &packet, nil
}
