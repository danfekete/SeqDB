package raft

func (r *RaftServer) RequestVote(args *RequestVoteArgs, reply *RequestVoteReply) error {

	// start from a state where the vote is refused
	reply.currentTerm = r.currentTerm
	reply.granted = false

	// The follower is in newer term than the vote requester
	if args.term < r.currentTerm {
		return nil
	}

	if r.votedFor == -1 || r.votedFor == args.candidateId {

		// TODO: check for log health
		r.currentTerm = args.term
		reply.granted = true
		reply.currentTerm = r.currentTerm
		return nil
	}

	// something really bad happened
	return nil
}