package communicate

import (
	"fmt"

	pb "github.com/singlemonad/raft/proto"
)

const (
	msgChannelSize = 100
)

type localCommunicator struct {
	msgC   chan *pb.Message
	others map[int]Communicator
}

func NewLocalCommunicator(nodeTotal int, others map[int]Communicator) Communicator {
	return &localCommunicator{
		msgC:   make(chan *pb.Message, msgChannelSize),
		others: others,
	}
}

func (c *localCommunicator) Send(msg *pb.Message) {
	targetID := int(msg.GetTo())
	if _, ok := c.others[targetID]; !ok {
		panic(fmt.Sprintf("can not communicate to %d", targetID))
	}
	c.others[targetID].Recv(msg)
}

func (c *localCommunicator) Recv(msg *pb.Message) {
	c.msgC <- msg
}

func (c *localCommunicator) GetMsgChannel() <-chan *pb.Message {
	return c.msgC
}
