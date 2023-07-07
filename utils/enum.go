package utils

type MatchingStrategy string

const (
	ViaTrim      MatchingStrategy = "trim"
	ViaUnmatched MatchingStrategy = "unmatched"
)

func (f MatchingStrategy) String() string {
	return string(f)
}

type StatusRequest string

const (
	Success StatusRequest = "success"
	Failed  StatusRequest = "failed"
)
