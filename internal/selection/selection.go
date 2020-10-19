package selection

import (
	"math/rand"
	"time"
)

type SelectionStrategy interface {
	Code() int
}

type RoundRobinSelectionStrategy struct {
	responses []int

	cursor int
}

func NewRoundRobinSelectionStrategy(r []int) *RoundRobinSelectionStrategy {
	return &RoundRobinSelectionStrategy{
		responses: r,
	}
}

func (s *RoundRobinSelectionStrategy) Code() int {
	code := s.responses[s.cursor]
	s.incCursor()
	return code
}

func (s *RoundRobinSelectionStrategy) incCursor() int {
	if s.cursor == (len(s.responses) - 1) {
		s.cursor = 0
		return s.cursor
	}

	s.cursor += 1
	return s.cursor
}

type RandomSelectionStrategy struct {
	responses []int
	rando     *rand.Rand
}

func NewRandomSelectionStrategy(r []int) *RandomSelectionStrategy {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	return &RandomSelectionStrategy{
		responses: r,
		rando:     r1,
	}
}

func (s *RandomSelectionStrategy) Code() int {
	return s.responses[s.rando.Intn(len(s.responses))]
}
