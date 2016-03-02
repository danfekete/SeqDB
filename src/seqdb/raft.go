package seqdb

import "math/rand"

type RaftStateType int

const (
	RAFT_STATE_LEADER = 1
	RAFT_STATE_CANDIDATE = 2
	RAFT_STATE_FOLLOWER = 3
)

type RaftMessage struct {

}

var (
	state RaftStateType // current state of our Raft server
 	minElectionTimeout int = 100
	maxElectionTimeout int = 300
	electionTimeout int
)

func NewRaftServer() {
	electionTimeout = rand.Intn(maxElectionTimeout - minElectionTimeout) + minElectionTimeout
}