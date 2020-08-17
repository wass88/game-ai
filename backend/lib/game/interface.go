package game

import "github.com/wass88/gameai/lib/protocol"

type Game interface {
	Start(players []*CmdRW, sender IPlayoutSender) (*protocol.Result, error)
}

type IPlayoutSender interface {
	Update(result protocol.ResultA) error
	Complete(results []protocol.ResultPlayerA) error
}
