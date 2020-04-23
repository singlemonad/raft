package log

import (
	pb "github.com/singlemonad/raft/proto"
)

type RaftLog interface {
	LastTerm() int64
	LastIndex() int64
	Match(logIndex, logTerm int64) bool
	IsUptoDate(lastIndex, lastTerm int64) bool
	GetLogTerm(index int64) int64
	GetLog(index int64) *pb.Entry
	GetLogs(startIndex int64) []*pb.Entry
	Append(entries []*pb.Entry)
}

// MemoryLog is raft log in memory, not persistent to disk
// entries[0] not use, log start at entries[1]
type MemoryLog struct {
	len     int // len is log length, equal len(entries) - 1
	entries []*pb.Entry
}

func NewMemoryLog() RaftLog {
	l := &MemoryLog{
		entries: make([]*pb.Entry, 0),
	}
	l.entries = append(l.entries, &pb.Entry{Term: 0})
	return l
}

func (l *MemoryLog) LastTerm() int64 {
	return l.entries[l.len].GetTerm()
}

func (l *MemoryLog) LastIndex() int64 {
	return int64(l.len)
}

func (l *MemoryLog) Match(logIndex, logTerm int64) bool {
	return l.entries[logIndex].GetTerm() == logTerm
}

func (l *MemoryLog) IsUptoDate(lastIndex, lastTerm int64) bool {
	if lastIndex > l.LastIndex() {
		return true
	} else if lastIndex < l.LastIndex() {
		return false
	} else {
		if lastTerm >= l.LastTerm() {
			return true
		}
		return false
	}
}

func (l *MemoryLog) GetLogTerm(index int64) int64 {
	return l.entries[index].GetTerm()
}

func (l *MemoryLog) GetLog(index int64) *pb.Entry {
	return l.entries[index]
}

func (l *MemoryLog) GetLogs(startIndex int64) []*pb.Entry {
	return l.entries[startIndex : l.len+1]
}

func (l *MemoryLog) Append(entries []*pb.Entry) {
	l.entries = append(l.entries, entries...)
	l.len += len(entries)
}
