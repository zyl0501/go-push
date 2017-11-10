package message

import (
	"encoding/binary"
	"math"
	"io"
	"bytes"
	"bufio"
	"github.com/zyl0501/go-push/api/protocol"
	"github.com/zyl0501/go-push/api"
)

type ByteBufMessage struct {
	byteBufMessageCodec

	Pkt        protocol.Packet
	Connection api.Conn
}

func (msg *ByteBufMessage) GetConnection() api.Conn {
	return msg.Connection
}

func (msg *ByteBufMessage) DecodeBody() {
	packet := msg.GetPacket()

	//1.解密
	tmp := packet.Body;
	//2.解压

	if len(tmp) == 0 {
		//"message decode ex"
		return
	}

	packet.Body = tmp
	msg.decodeBaseMessage(packet.Body)
	packet.Body = nil // 释放内存
}

func (msg *ByteBufMessage) EncodeBody() {
	tmp := msg.encodeBaseMessage();
	if len(tmp) > 0 {
		//1.压缩
		//2.加密

		msg.Pkt.Body = tmp
	}
}

func (msg *ByteBufMessage) GetPacket() protocol.Packet {
	return msg.Pkt
}

func (msg *ByteBufMessage) Send() {
	msg.EncodeBody()
	writer := bufio.NewWriter(msg.GetConnection().GetConn())
	writer.Write(protocol.EncodePacket(msg.GetPacket()))
}

func (msg *ByteBufMessage) decodeBaseMessage(body []byte) {
	msg.decodeByteBufMessage(bytes.NewReader(body))
}

func (msg *ByteBufMessage) encodeBaseMessage() ([]byte) {
	buf := bytes.NewBuffer(make([]byte, 0))
	writer := bufio.NewWriter(buf)
	msg.encodeByteBufMessage(writer)
	writer.Flush()
	return buf.Bytes()
}

func (msg *ByteBufMessage) sendRaw() {
	msg.encodeRaw()
	writer := bufio.NewWriter(msg.GetConnection().GetConn())
	writer.Write(protocol.EncodePacket(msg.GetPacket()))
}

func (msg *ByteBufMessage) encodeRaw() {
	tmp := msg.encodeBaseMessage()
	msg.Pkt.Body = tmp
}

type byteBufMessageCodec interface {
	decodeByteBufMessage(reader io.Reader)
	encodeByteBufMessage(writer io.Writer)
}

//***********************Encode/Decode method***********************

func EncodeString(writer io.Writer, field string) {
	EncodeBytes(writer, []byte(field))
}

func EncodeBytes(writer io.Writer, field []byte) {
	fieldLen := len(field)
	if fieldLen == 0 {
		EncodeInt16(writer, 0)
	} else if fieldLen < math.MaxInt16 {
		EncodeInt16(writer, int16(fieldLen))
		encodeBytes(writer, field)
	} else {
		EncodeInt16(writer, math.MaxInt16)
		EncodeInt32(writer, int32(fieldLen-math.MaxInt16))
		encodeBytes(writer, field)
	}
}

func EncodeInt16(writer io.Writer, field int16) {
	encode(writer, field)
}

func EncodeInt32(writer io.Writer, field int32) {
	encode(writer, field)
}

func EncodeInt64(writer io.Writer, field int64) {
	encode(writer, field)
}

func EncodeByte(writer io.Writer, field byte) {
	encode(writer, field)
}

func encodeBytes(writer io.Writer, field []byte) {
	encode(writer, field)
}

func encode(writer io.Writer, field interface{}) {
	binary.Write(writer, binary.BigEndian, field)
}

func DecodeString(reader io.Reader) (string) {
	buf := DecodeBytes(reader)
	if buf == nil {
		var field string
		return field
	} else {
		return string(buf)
	}
}

func DecodeBytes(reader io.Reader) (field []byte) {
	var fieldLength int32
	fieldLength = int32(DecodeInt16(reader))
	if fieldLength == 0 {
		return field
	} else if fieldLength == math.MaxInt16 {
		fieldLength += DecodeInt32(reader)
	}
	return decodeBytes(reader, fieldLength)
}

func DecodeInt16(reader io.Reader) (field int16) {
	decode(reader, &field)
	return field
}
func DecodeInt32(reader io.Reader) (field int32) {
	decode(reader, &field)
	return field
}
func DecodeInt64(reader io.Reader) (field int64) {
	decode(reader, &field)
	return field
}
func DecodeByte(reader io.Reader) (field byte) {
	decode(reader, &field)
	return field
}

func decodeBytes(reader io.Reader, len int32) ([]byte) {
	field := make([]byte, len)
	decode(reader, &field)
	return field
}

func decode(reader io.Reader, field interface{}) {
	binary.Read(reader, binary.BigEndian, field)
}
