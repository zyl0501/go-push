package boot

import (
	"net"
	"os"
	log "github.com/alecthomas/log4go"
	"io"
)


func Init() {
	netListen, err := net.Listen("tcp", "localhost:9932")
	if err != nil {
		log.Error(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
	defer log.Info("server exit")
	defer netListen.Close()

	log.Info("Wait for Client")
	for {
		conn, err := netListen.Accept()
		if err != nil {
			continue
		}
		log.Info(conn.RemoteAddr().String(), "tcp connect success")

		handleConnection(conn)
	}
}

func Start(){

}

func handleConnection(conn net.Conn) {
	buffer := make([]byte, 2048)

	for {
		n, err := conn.Read(buffer)
		if err != nil {
			if err != io.EOF {
				log.Error(conn.RemoteAddr().String(), "connect error:", err)
			}
			break
		}
		log.Info(conn.RemoteAddr().String(), "receive data string:", string(buffer[:n]))
	}
}
