package shared


const (
	HashPrefixLength = 5
)

type PState uint

const (
	Playing PState = iota
	Paused
	Stopped
)


