package communicate

import (
	pb "github.com/singlemonad/raft/proto"
)

type Communicator interface {
	Send(message *pb.Message)
	Recv(message *pb.Message)

	GetMsgChannel() <-chan *pb.Message
}
