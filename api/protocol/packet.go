package protocol

type Packet struct {
	Cmd       byte  //命令
	cc        int16 //校验码 暂时没有用到
	flags     byte  //特性，如是否加密，是否压缩等
	sessionId int   // 会话id。客户端生成。
	lrc       byte  // 校验，纵向冗余校验。只校验head
	body      string
}
