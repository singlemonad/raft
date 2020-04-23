package raft

import (
	"fmt"
	"sync"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/singlemonad/raft"
	pb "github.com/singlemonad/raft/proto"
)

type RaftServerCfg struct {
	raftCfg raft.RaftNodeCfg
}

// not thread safe, only support serial operation
type RaftServer struct {
	mutex    sync.Mutex
	kvs      map[string]string
	errorC   chan error
	exitC    chan interface{}
	raftNode *raft.RaftNode
}

func NewRaftServer(cfg RaftServerCfg) *RaftServer {
	s := &RaftServer{
		mutex:    sync.Mutex{},
		kvs:      make(map[string]string),
		errorC:   make(chan error),
		exitC:    make(chan interface{}),
		raftNode: raft.NewRaftNode(cfg.raftCfg),
	}
	go s.raftNode.Run()
	go func() {
		for {
			select {
			case <-time.NewTicker(time.Millisecond).C:
				s.raftNode.Tick()
			case <-s.exitC:
				return
			}
		}
	}()
	go s.fetchCommitLog()

	return s
}

func (s *RaftServer) Put(key, value string) error {
	kv := &pb.KV{Key: key, Value: value, Op: pb.OperateType_Put}
	body, err := ptypes.MarshalAny(kv)
	if err != nil {
		return err
	}

	s.raftNode.Propose(body)
	select {
	case err := <-s.errorC:
		return err
	case <-time.After(time.Second * 5):
		return fmt.Errorf("put kv timeout")
	}
}

func (s *RaftServer) Get(key string) (string, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if val, ok := s.kvs[key]; ok {
		return val, nil
	}
	return "", nil
}

func (s *RaftServer) Delete(key string) error {
	kv := &pb.KV{Key: key, Op: pb.OperateType_Delete}
	body, err := ptypes.MarshalAny(kv)
	if err != nil {
		return err
	}

	s.raftNode.Propose(body)
	select {
	case err := <-s.errorC:
		return err
	case <-time.After(time.Second * 5):
		return fmt.Errorf("put kv timeout")
	}
}

func (s *RaftServer) Stop() {
	s.raftNode.Stop()
	close(s.exitC)
}

func (s *RaftServer) fetchCommitLog() {
	for {
		select {
		case entry := <-s.raftNode.GetApplyLogChannel():
			kv := &pb.KV{}
			ptypes.UnmarshalAny(entry.GetBody(), kv)

			s.mutex.Lock()
			switch kv.Op {
			case pb.OperateType_Put:
				s.kvs[kv.Key] = kv.Value
			case pb.OperateType_Delete:
				delete(s.kvs, kv.Key)
			default:
			}
			s.mutex.Unlock()

			if s.raftNode.IsLeader() {
				s.errorC <- nil
			}
		case <-s.exitC:
			return
		}
	}
}
