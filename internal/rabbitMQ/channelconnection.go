package rabbitMQ

import(
	"github.com/NeowayLabs/wabbit"
)

func ChannelConnection(conn wabbit.Conn) (wabbit.Channel, error) {
	ch, err := conn.Channel()
	return ch, err
}
