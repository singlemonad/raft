package raft

import pb "github.com/singlemonad/raft/proto"

// RaftNodeState is node state while node in raft cluster
type RaftNodeState struct {
	id           int
	role         pb.Role
	leader       int
	commitIndex  int64
	appliedIndex int64
	hardState    *HardState
}

// FollowerIndex is Leader record next index and match index for every follower
type FollowerIndex struct {
	nextIndex  []int64
	matchIndex []int64
}

func NewFollowerIndex(nodes int) *FollowerIndex {
	f := &FollowerIndex{
		nextIndex:  make([]int64, nodes+1),
		matchIndex: make([]int64, nodes+1),
	}
	for i := 0; i <= nodes; i++ {
		f.nextIndex[i] = 1
	}
	return f
}

func (f *FollowerIndex) GetNextIndex(id int) int64 {
	return f.nextIndex[id]
}

func (f *FollowerIndex) SetNextIndex(id int, next int64) {
	f.nextIndex[id] = next
}

func (f *FollowerIndex) GetMatchIndex(id int) int64 {
	return f.matchIndex[id]
}

func (f *FollowerIndex) SetMatchIndex(id int, match int64) {
	f.matchIndex[id] = match
}

// HardState need store to disk for every node
type HardState struct {
	term    int64
	voteFor int
}

func NewHardState() *HardState {
	return &HardState{}
}

func (h *HardState) GetTerm() int64 {
	return h.term
}

func (h *HardState) SetTerm(term int64) {
	h.term = term
}

func (h *HardState) GetVoteFor() int {
	return h.voteFor
}

func (h *HardState) SetVoteFor(voteFor int) {
	h.voteFor = voteFor
}

// Election is statistical info in a election
type Election struct {
	voteCounter int
}

func NewElection() *Election {
	return &Election{}
}

func (e *Election) AddVote() {
	e.voteCounter++
}

func (e *Election) GetMajorityVote(total int) bool {
	return e.voteCounter >= (total/2)+1
}
