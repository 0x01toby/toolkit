package polling

import "strings"

type Action string

var (
	ContinuePolling = Action("ContinuePolling")
	WaitingBlocks   = Action("WaitingBlocks")
	Done            = Action("PollingDone")
)

type NextPollingAction struct {
	iStart uint64
	iEnd   uint64
	action Action
}

func NewNextPollingAction(start, end uint64, action Action) *NextPollingAction {
	return &NextPollingAction{iStart: start, iEnd: end, action: action}
}

func (n *NextPollingAction) StartHeight() uint64 {
	return n.iStart
}

func (n *NextPollingAction) EndHeight() uint64 {
	return n.iEnd
}

func (n *NextPollingAction) IsContinuePolling() bool {
	return strings.EqualFold(string(n.action), string(ContinuePolling))
}

func (n *NextPollingAction) IsDone() bool {
	return strings.EqualFold(string(n.action), string(Done))
}

func (n *NextPollingAction) IsWaiting() bool {
	return strings.EqualFold(string(n.action), string(WaitingBlocks))
}
