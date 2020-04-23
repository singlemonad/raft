package raft

import (
	"flag"
	"math/rand"
	"testing"
	"time"

	"github.com/singlemonad/raft"
	"github.com/singlemonad/raft/communicate"
	"github.com/stretchr/testify/assert"
)

var (
	nodeTotalFlag          int
	electionTimeoutFlag    int
	requestVoteTimeoutFlag int
)

func init() {
	flag.IntVar(&nodeTotalFlag, "nodetotal", 0, "total node in raft")
	flag.IntVar(&electionTimeoutFlag, "election_timeout", 0, "election timout Millisecond")
	flag.IntVar(&requestVoteTimeoutFlag, "requestvote_timeout", 0, "request vote timeout Millisecond")
}

func waitElectionFinish(silenceTime time.Duration) {
	time.Sleep(silenceTime)
}

var (
	gComms      map[int]communicate.Communicator
	gRaftServer map[int]*RaftServer
)

func initCommunicator() {
	gComms = make(map[int]communicate.Communicator)
	for i := 1; i <= nodeTotalFlag; i++ {
		gComms[i] = communicate.NewLocalCommunicator(i, gComms)
	}
}

func initRaftServer() {
	gRaftServer = make(map[int]*RaftServer)
	for i := 1; i <= nodeTotalFlag; i++ {
		gRaftServer[i] = NewRaftServer(RaftServerCfg{raftCfg: raft.RaftNodeCfg{
			Id:                 i,
			NodeTotal:          nodeTotalFlag,
			ElectionTimeout:    electionTimeoutFlag + int(rand.NewSource(time.Now().UnixNano()).Int63()%int64((nodeTotalFlag*2))),
			RequestVoteTimeout: requestVoteTimeoutFlag,
			Communicator:       gComms[i],
		}})
	}
}

func getLeaderNum() int {
	for num, s := range gRaftServer {
		if s.raftNode.IsLeader() {
			return num
		}
	}
	return 0
}

func TestPropose(t *testing.T) {
	initCommunicator()
	initRaftServer()

	waitElectionFinish(time.Millisecond * 20)

	leaderIdx := getLeaderNum()
	assert.NotEqual(t, 0, leaderIdx)

	var err error
	leader := gRaftServer[leaderIdx]
	err = leader.Put("name", "yang")
	assert.Equal(t, nil, err)
	//time.Sleep(time.Millisecond * 50)

	err = leader.Put("name-1", "yang-1")
	assert.Equal(t, nil, err)
	//time.Sleep(time.Millisecond * 50)

	err = leader.Put("name-2", "yang-2")
	assert.Equal(t, nil, err)
	//time.Sleep(time.Millisecond * 50)

	var value string
	value, err = leader.Get("name")
	assert.Equal(t, nil, err)
	assert.Equal(t, "yang", value)
	//time.Sleep(time.Millisecond * 50)

	err = leader.Delete("name")
	assert.Equal(t, nil, err)
	//time.Sleep(time.Millisecond * 50)

	value, err = leader.Get("name")
	assert.Equal(t, nil, err)
	assert.Equal(t, "", value)

	leader.Stop()
	waitElectionFinish(time.Millisecond * 20)
	delete(gRaftServer, leaderIdx)

	leaderIdx = getLeaderNum()
	assert.NotEqual(t, 0, leaderIdx)
	leader = gRaftServer[leaderIdx]

	err = leader.Put("name-3", "yang-3")
	assert.Equal(t, nil, err)
	value, err = leader.Get("name-3")
	assert.Equal(t, nil, err)
	assert.Equal(t, "yang-3", value)

	leader.Stop()
	waitElectionFinish(time.Millisecond * 20)
	delete(gRaftServer, leaderIdx)

	leaderIdx = getLeaderNum()
	assert.NotEqual(t, 0, leaderIdx)

	leader = gRaftServer[leaderIdx]
	value, err = leader.Get("name-2")
	assert.Equal(t, nil, err)
	assert.Equal(t, "yang-2", value)

	err = leader.Put("name-4", "yang-4")
	assert.Equal(t, nil, err)
	value, err = leader.Get("name-4")
	assert.Equal(t, nil, err)
	assert.Equal(t, "yang-4", value)

	waitElectionFinish(time.Millisecond * 10)
}

func TestInitElection(t *testing.T) {
	initCommunicator()
	initRaftServer()

	waitElectionFinish(time.Millisecond * 20)

	leaderIdx := getLeaderNum()
	assert.NotEqual(t, 0, leaderIdx)

	var err error
	leader := gRaftServer[leaderIdx]
	err = leader.Put("name", "yang")
	assert.Equal(t, nil, err)
}
