package log

import (
	pb "github.com/singlemonad/raft/proto"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLog(t *testing.T) {
	log := NewRaftLog()
	assert.Equal(t, int64(-1), log.lastIndex())
	assert.Equal(t, int64(-1), log.lastTerm())
	assert.Equal(t, true, log.isUpToDate(-1, -1))
	entries := []*pb.Entry{
		&pb.Entry{
			Term: 1,
		},
		&pb.Entry{
			Term: 1,
		},
	}
	log.append(entries)
	assert.Equal(t, int64(1), log.lastIndex())
	assert.Equal(t, int64(1), log.lastTerm())

	entries = []*pb.Entry{
		&pb.Entry{
			Term: 1,
		},
		&pb.Entry{
			Term: 1,
		},
	}
	log.append(entries)
	assert.Equal(t, int64(3), log.lastIndex())
	assert.Equal(t, int64(1), log.lastTerm())
}
