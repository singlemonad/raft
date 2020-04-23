package raft

import (
	"sort"
	"time"

	pb "github.com/singlemonad/raft/proto"
)

func (r *RaftNode) electionTimeout(msg *pb.Message) {
	r.StartElection()
}

func (r *RaftNode) pause(msg *pb.Message) {
	<-time.NewTimer(time.Duration(msg.GetPasueTime()) * time.Millisecond).C
}

func (r *RaftNode) vote(msg *pb.Message) {
	canVote := ((msg.GetTerm() > r.hardState.GetTerm() || r.hardState.GetVoteFor() == 0) &&
		r.raftLog.IsUptoDate(msg.GetLogIndex(), msg.GetLogTerm())) ||
		int64(r.hardState.GetVoteFor()) == msg.GetFrom()
	if canVote {
		r.electionTimeoutCounter = 0
		r.hardState.SetTerm(msg.GetTerm())
		r.hardState.SetVoteFor(int(msg.GetFrom()))
		r.Send(&pb.Message{Type: pb.MessageType_VoteReply, To: msg.GetFrom()})
		return
	}
	r.Send(&pb.Message{Type: pb.MessageType_VoteReply, To: msg.GetFrom(), Reject: true})
}

func (r *RaftNode) voteReply(msg *pb.Message) {
	// if we are not in election, ignore it.
	if r.election == nil {
		return
	}

	if msg.GetReject() {
		return
	}
	r.election.AddVote()
	if r.election.GetMajorityVote(r.cfg.NodeTotal) {
		r.becomeLeader()
	}
}

func (r *RaftNode) appendEntries(msg *pb.Message) {
	match := r.raftLog.Match(msg.GetLogIndex(), msg.GetLogTerm())
	if !match {
		r.Send(&pb.Message{Type: pb.MessageType_AppendEntriesReply, To: msg.GetFrom(), Reject: true})
		return
	}

	if len(msg.GetEntries()) == 0 {
		// update commit index
		if r.commitIndex < msg.GetCommitIndex() {
			r.commitIndex = min(msg.GetCommitIndex(), r.raftLog.LastIndex())
		}
		r.becomeFollower(msg)
		return
	}

	r.raftLog.Append(msg.GetEntries())
	r.commitIndex = min(msg.GetCommitIndex(), r.raftLog.LastIndex())
	r.Send(&pb.Message{Type: pb.MessageType_AppendEntriesReply, To: msg.GetFrom()})
}

func (r *RaftNode) appendEntriesReply(msg *pb.Message) {
	from := int(msg.GetFrom())

	if msg.GetReject() {
		r.followerIndex.SetNextIndex(from, r.followerIndex.GetNextIndex(from)-1)
		r.syncFollower(int(msg.GetFrom()))
		return
	}

	// update follower next index and match index
	r.followerIndex.SetNextIndex(from, r.raftLog.LastIndex()+1)
	r.followerIndex.SetMatchIndex(from, r.raftLog.LastIndex())

	lastIndex := r.raftLog.LastIndex()
	if r.raftLog.GetLogTerm(lastIndex) == r.hardState.GetTerm() {
		r.updateCommitIndex()
	}
}

func (r *RaftNode) updateCommitIndex() {
	matchArr := make([]int64, 0)
	for i, num := range r.followerIndex.matchIndex {
		if i != r.id {
			matchArr = append(matchArr, num)
		}
	}
	matchArr = append(matchArr, r.raftLog.LastIndex())
	sort.SliceIsSorted(matchArr, func(i, j int) bool {
		return matchArr[i] < matchArr[j]
	})
	r.commitIndex = matchArr[len(matchArr)/2+1]
}

func min(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}
