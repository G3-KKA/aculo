package broker

import (
	"strconv"
	"sync/atomic"
)

//go:generate mockery --filename=mock_topic_name_generator.go --name=TopicNameGenerator --dir=. --structname MockTopicNameGenerator  --inpackage=true
type TopicNameGenerator interface {
	Generate() string
}

// # Default log topic name generator.
//
// Produces names like logtopic1,logtopic2, .. logtopicN.
type logTopicNamegen struct {
	counter atomic.Int32
}

// Produces names like logtopic1,logtopic2, .. logtopicN.
func (namegen *logTopicNamegen) Generate() string {
	idx := namegen.counter.Load()
	for !namegen.counter.CompareAndSwap(idx, idx+1) {
		idx = namegen.counter.Load()
	}
	return "logtopic" + strconv.Itoa(int(idx))
}
