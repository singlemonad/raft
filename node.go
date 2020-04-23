package raft

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang/protobuf/ptypes/any"
	"github.com/singlemonad/raft/communicate"
	"github.com/singlemonad/raft/log"
	pb "github.com/singlemonad/raft/proto"
	"go.uber.org/zap"
)

const (
	appliedLogChannelSize = 1000
)

type stepFunc func()

type RaftNodeCfg struct {
	Id                 int
	NodeTotal          int
	ElectionTimeout    int
	RequestVoteTimeout int
	Communicator       communicate.Communicator
}

type RaftNode struct {
	cfg    RaftNodeCfg
	logger *zap.SugaredLogger

	RaftNodeState

	raftLog     log.RaftLog
	stepFunc    stepFunc
	appliedLogC chan *pb.Entry
	exitC       chan struct{}

	electionTimeoutCounter int
	election               *Election
	followerIndex          *FollowerIndex
}

func NewRaftNode(cfg RaftNodeCfg) *RaftNode {
	nodeState := RaftNodeState{
		id:        cfg.Id,
		role:      pb.Role_Follower,
		hardState: NewHardState(),
	}
	logger, _ := zap.NewDevelopment()
	r := &RaftNode{
		cfg:           cfg,
		logger:        logger.Sugar(),
		RaftNodeState: nodeState,
		raftLog:       log.NewMemoryLog(),
		appliedLogC:   make(chan *pb.Entry, appliedLogChannelSize),
		exitC:         make(chan struct{}),
	}
	r.stepFunc = r.electionTimeoutCheck
	return r
}

func (r *RaftNode) Propose(body *any.Any) error {
	if r.role != pb.Role_Leader {
		return errors.New("propose to Follower or Candidate is invalid")
	}
	r.Send(&pb.Message{Type: pb.MessageType_AppendEntries, From: int64(r.id), To: int64(r.id), Entries: []*pb.Entry{&pb.Entry{Body: body}}})
	return nil
}

func (r *RaftNode) IsLeader() bool {
	return r.id == r.leader
}

func (r *RaftNode) Tick() {
	r.stepFunc()
}

func (r *RaftNode) GetApplyLogChannel() <-chan *pb.Entry {
	return r.appliedLogC
}

func (r *RaftNode) Stop() {
	close(r.exitC)
}

func (r *RaftNode) Run() {
	defer func() {
		r.logger.Infof("[Id: %d] exit raft loop.", r.id)
	}()

	ticker := time.NewTicker(time.Millisecond)
	for {
		var msg *pb.Message

		select {
		case msg = <-r.cfg.Communicator.GetMsgChannel():
			//r.logger.Debugf("%s receive message [%s]", r.nodeState(), msg.String())
			r.Step(msg)
		case <-ticker.C:
			r.fetchApplyLog()
		case <-r.exitC:
			return
		}
	}
}

// Pause RaftNode for a while, just for test
func (r *RaftNode) Pause(time int) {
	r.Send(&pb.Message{PasueTime: int64(time), To: int64(r.id)})
}

func (r *RaftNode) Send(msg *pb.Message) {
	msg.From = int64(r.id)
	msg.Term = r.hardState.GetTerm()
	r.cfg.Communicator.Send(msg)

	//r.logger.Debugf("%s send message [%s]", r.nodeState(), msg.String())
}

func (r *RaftNode) Step(msg *pb.Message) {
	if msg.GetFrom() == int64(r.id) && msg.GetType() == pb.MessageType_AppendEntries {
		r.handlePropose(msg.GetEntries())
		return
	}

	switch {
	case msg.GetTerm() < r.hardState.GetTerm():
		if msg.GetType() == pb.MessageType_Vote {
			// out of date vote msg，reject it
			r.Send(&pb.Message{Type: pb.MessageType_VoteReply, To: msg.GetFrom(), Reject: true})
		}
		r.logger.Infof("%s receive out of date message [%v]", r.nodeState(), msg)
		return
	case msg.GetTerm() > r.hardState.GetTerm():
		// msg term bigger than mine, update term and reset voterFor, than become follower
		if msg.GetType() == pb.MessageType_Vote ||
			msg.GetType() == pb.MessageType_VoteReply ||
			msg.GetType() == pb.MessageType_AppendEntriesReply {
			r.resetToFollower(msg.GetTerm())
		}

		// if AppendEntriesReply's term bigger than mine, return
		if msg.GetType() == pb.MessageType_AppendEntriesReply {
			return
		}
	}

	switch msg.GetType() {
	case pb.MessageType_ElectionTimeout:
		r.electionTimeout(msg)
	case pb.MessageType_Vote:
		r.vote(msg)
	case pb.MessageType_VoteReply:
		r.voteReply(msg)
	case pb.MessageType_AppendEntries:
		r.appendEntries(msg)
	case pb.MessageType_AppendEntriesReply:
		r.appendEntriesReply(msg)
	case pb.MessageType_Puase:
		r.pause(msg)
	}
}

func (r *RaftNode) StartElection() {
	r.becomeCandidate()
	r.startElection()
}

func (r *RaftNode) fetchApplyLog() {
	for i := r.appliedIndex + 1; i <= r.commitIndex; i++ {
		entry := r.raftLog.GetLog(int64(i))
		r.appliedLogC <- entry
		r.appliedIndex++
	}
}

func (r *RaftNode) resetToFollower(term int64) {
	r.role = pb.Role_Follower
	r.electionTimeoutCounter = 0
	r.hardState.SetTerm(term)
	r.hardState.SetVoteFor(0)
}

// Leader handle user propose
func (r *RaftNode) handlePropose(entries []*pb.Entry) {
	for _, entry := range entries {
		entry.Term = r.hardState.GetTerm()
	}
	r.raftLog.Append(entries)
	r.syncAllFollower()
}

func (r *RaftNode) syncFollower(followerID int) {
	preLogIndex := r.followerIndex.GetNextIndex(followerID) - 1
	r.Send(&pb.Message{Type: pb.MessageType_AppendEntries, To: int64(followerID), LogIndex: preLogIndex, LogTerm: r.raftLog.GetLogTerm(preLogIndex), CommitIndex: r.commitIndex, Entries: r.raftLog.GetLogs(preLogIndex)})
}

func (r *RaftNode) syncAllFollower() {
	// get need sync followers
	needSyncFollowers := make([]int, 0)
	for i := 1; i <= r.cfg.NodeTotal; i++ {
		if i == r.id {
			continue
		}

		if r.followerNeedAppendLog(i) {
			needSyncFollowers = append(needSyncFollowers, i)
		}
	}

	// sync followers
	for _, followerID := range needSyncFollowers {
		preLogIndex := r.followerIndex.GetNextIndex(followerID) - 1
		r.Send(&pb.Message{Type: pb.MessageType_AppendEntries, To: int64(followerID), LogIndex: preLogIndex, LogTerm: r.raftLog.GetLogTerm(preLogIndex), CommitIndex: r.commitIndex, Entries: r.raftLog.GetLogs(preLogIndex + 1)})
	}
}

func (r *RaftNode) followerNeedAppendLog(followerID int) bool {
	if r.followerIndex.GetNextIndex(followerID) <= r.raftLog.LastIndex() {
		return true
	}
	return false
}

func (r *RaftNode) leaderHeartbeat() {
	for followerID := 1; followerID <= r.cfg.NodeTotal; followerID++ {
		if followerID == r.id {
			continue
		}

		r.Send(&pb.Message{Type: pb.MessageType_AppendEntries, To: int64(followerID), CommitIndex: r.commitIndex})
	}
}

func (r *RaftNode) becomeLeader() {
	r.logger.Debugf("%s become leader", r.nodeState())

	r.restElection()
	r.leader = r.id
	r.role = pb.Role_Leader
	r.stepFunc = r.leaderHeartbeat
	r.followerIndex = NewFollowerIndex(r.cfg.NodeTotal)
	for i := 1; i <= r.cfg.NodeTotal; i++ {
		if i != r.id {
			r.followerIndex.SetNextIndex(i, r.raftLog.LastIndex()+1)
			r.followerIndex.SetMatchIndex(i, 0)
		}
	}
	r.leaderHeartbeat()
}

func (r *RaftNode) restElection() {
	r.election = nil
}

func (r *RaftNode) becomeCandidate() {
	r.election = NewElection()
	r.role = pb.Role_Candidate
	r.election.AddVote()
	r.hardState.SetVoteFor(r.id)
	r.electionTimeoutCounter = 0
	r.stepFunc = r.electionTimeoutCheck
	r.hardState.SetTerm(r.hardState.GetTerm() + 1)
}

func (r *RaftNode) electionTimeoutCheck() {
	r.electionTimeoutCounter++
	if r.electionTimeoutCounter > r.cfg.ElectionTimeout {
		r.Send(&pb.Message{Type: pb.MessageType_ElectionTimeout, To: int64(r.id)})
		r.electionTimeoutCounter = 0
	}
}

func (r *RaftNode) startElection() {
	r.electionTimeoutCounter = 0
	for i := 1; i <= r.cfg.NodeTotal; i++ {
		if i != r.id {
			r.Send(&pb.Message{Type: pb.MessageType_Vote, To: int64(i), LogIndex: r.raftLog.LastIndex(), LogTerm: r.raftLog.LastTerm()})
		}
	}
}

func (r *RaftNode) becomeFollower(msg *pb.Message) {
	r.restElection()
	r.electionTimeoutCounter = 0
	r.leader = int(msg.GetFrom())
	r.stepFunc = r.electionTimeoutCheck
	r.hardState.SetTerm(msg.GetTerm())
	r.hardState.SetVoteFor(r.leader)
}

func (r *RaftNode) nodeState() string {
	return fmt.Sprintf("[Role: %s Id: %d Term: %d Leader: %d VoteFor: %d LastIndex: %d LastTerm: %d]",
		r.role.String(), r.id, r.hardState.GetTerm(), r.leader, r.hardState.GetVoteFor(), r.raftLog.LastIndex(), r.raftLog.LastTerm())
}
