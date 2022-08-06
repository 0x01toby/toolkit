package event

type SortListenerItems []*ListenerItem

func (s SortListenerItems) Len() int {
	return len(s)
}

func (s SortListenerItems) Less(i, j int) bool {
	return s[i].Priority > s[j].Priority
}

func (s SortListenerItems) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
