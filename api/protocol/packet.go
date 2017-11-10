package protocol

import (
	"encoding/binary"
	"bytes"
	"sync/atomic"
)

//length(4)+cmd(1)+cc(2)+flags(1)+sessionId(4)+lrc(1)+body(n)

const (
	HeadLength      byte = 13 //packet头部的长度
	BodyLength      byte = 4
	CmdLength       byte = 1
	CCLength        byte = 2
	FlagsLength     byte = 1
	SessionIdLength byte = 4
	LrcLength       byte = 1

	BodyLenIndex   = 0
	CmdIndex       = 4
	CCIndex        = 5
	FlagsIndex     = 6
	SessionIdIndex = 7
	LrcIndex       = 11

	FLAG_CRYPTO    = 1
	FLAG_COMPRESS  = 2
	FLAG_BIZ_ACK   = 4
	FLAG_AUTO_ACK  = 8
	FLAG_JSON_BODY = 16
)

var(
	id_seq uint32
)

type Packet struct {
	Cmd       uint8  `json:"cmd"`         //命令
	Cc        uint16 `json:"cc"`          //校验码 暂时没有用到
	Flags     uint8  `json:"flags"`       //特性，如是否加密，是否压缩等
	SessionId uint32   `json:"sessionId"` // 会话id。客户端生成。
	Lrc       uint8  `json:"lrc"`         // 校验，纵向冗余校验。只校验head
	Body      []byte`json:"body"`
}

//Packet Cmd
const (
	HEARTBEAT            byte = 1 + iota //1
	HANDSHAKE                            //2
	LOGIN                                //3
	LOGOUT                               //4
	BIND                                 //5
	UNBIND                               //6
	FAST_CONNECT                         //7
	PAUSE                                //8
	ERROR                                //9
	OK                                   //10
	HTTP_PROXY                           //11
	KICK                                 //12
	GATEWAY_KICK                         //13
	PUSH                                 //14
	GATEWAY_PUSH                         //15
	NOTIFICATION                         //16
	GATEWAY_NOTIFICATION                 //17
	CHAT                                 //18
	GATEWAY_CHAT                         //19
	GROUP                                //20
	GATEWAY_GROUP                        //21
	ACK                                  //22
	NACK                                 //23
	UNKNOWN              = -1            //-1
)

func (packet *Packet) GetBodyLength() uint32 {
	if packet.Body == nil {
		return 0
	} else {
		return uint32(len(packet.Body))
	}
}

func (packet *Packet) HasFlag(flag byte) bool {
	return (packet.Flags & flag) != 0;
}

func DecodePacket(buf []byte) (Packet, uint32) {
	bodyLength := buf[0:4]
	cmd := buf[4:5]
	cc := buf[5:7]
	flags := buf[7:8]
	sessionId := buf[8:12]
	lrc := buf[12:13]

	var body []byte
	if len(buf) == int(binary.BigEndian.Uint32(bodyLength)+13) {
		body = buf[13:]
	}
	return Packet{
		cmd[0],
		binary.BigEndian.Uint16(cc),
		flags[0],
		binary.BigEndian.Uint32(sessionId),
		lrc[0],
		body}, binary.BigEndian.Uint32(bodyLength)
}

func EncodePacket(packet Packet) []byte {
	bodyLength := packet.GetBodyLength()

	buf := make([]byte, 13+bodyLength)
	copy(buf[0:4], int32ToBytes(bodyLength))
	buf[4] = packet.Cmd
	copy(buf[5:7], int16ToBytes(packet.Cc))
	buf[7] = packet.Flags
	copy(buf[8:12], int32ToBytes(packet.SessionId))
	buf[12] = packet.Lrc
	if bodyLength > 0 {
		copy(buf[13:], packet.Body)
	}
	return buf
}

func GetSessionId() uint32{
	return atomic.AddUint32(&id_seq, 1)
}

//整形转换成字节
func int32ToBytes(n uint32) []byte {
	tmp := n
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, tmp)
	return bytesBuffer.Bytes()
}

func int16ToBytes(n uint16) []byte {
	tmp := n
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, tmp)
	return bytesBuffer.Bytes()
}

//字节转换成整形
func bytesToInt32(b []byte) uint32 {
	bytesBuffer := bytes.NewBuffer(b)
	var tmp uint32
	binary.Read(bytesBuffer, binary.BigEndian, &tmp)
	return tmp
}

func bytesToInt16(b []byte) uint16 {
	bytesBuffer := bytes.NewBuffer(b)
	var tmp uint16
	binary.Read(bytesBuffer, binary.BigEndian, &tmp)
	return tmp
}
