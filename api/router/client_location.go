package router

type ClientLocation struct {
	Host          string //长链接所在的机器IP
	Port          int    //长链接所在的机器端口
	OsName        string //客户端系统类型
	ClientVersion string //客户端版本
	DeviceId      string //客户端设备ID
	ConnId        string //链接ID
	ClientType    int    //客户端类型
}
