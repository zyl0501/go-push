package api

const(
	STATUS_NEW byte = 0
	STATUS_CONNECTED byte = 1
	STATUS_DISCONNECTED byte = 2
)

type Connection interface {
	Init()

	IsConnected() bool
}

type ConnectionManager interface {


	Connection get(Channel channel);

Connection removeAndClose(Channel channel);

void add(Connection connection);

int getConnNum();

void init();

void destroy();
}