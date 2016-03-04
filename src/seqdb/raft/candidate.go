package raft
import (
	"net/rpc"
	"log"
)

func (r *RaftServer) elect() {

	conn, err := rpc.DialHTTP("tcp", "localhost:12345")

	if err != nil {
		log.Println("Cannot connect to RPC server", err)
	}

	var reply = &RequestVoteReply{}
	var args = &RequestVoteArgs{
		term: r.currentTerm,
		candidateId: r.myId,
		lastLogIndex: 0,
		lastLogTerm: 0,
	}
	err = conn.Call("RaftServer.RequestVote", args, reply)

	if err != nil {
		log.Println("RequestVote RPC call failed", err)
	}


}