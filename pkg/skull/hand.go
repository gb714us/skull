package skull

import (
	"errors"
)

var (
	ErrOutOfBound = errors.New("out of bounds")
	ErrCantPlay   = errors.New("can't play any cards")
)

// type HandManager interface {
// 	Play(position int) (CardType, error)
// 	CurrentSize() int
// 	CanPlay() bool
// 	Remove(position int) error
// 	MaxSize() int
// 	Reset()
// 	Restart()
// }

type Hand struct {
	removedStartIdx int // tells us the end index of the available cards (exclusive)
	playedStartIdx  int // tells us the end index of the unused cards (exclusive)
	cards           []CardType
}

func NewHand() Hand {
	return Hand{
		cards:           []CardType{CardTypeFlower, CardTypeFlower, CardTypeFlower, CardTypeSkull},
		removedStartIdx: 4,
		playedStartIdx:  4,
	}
}

func (h Hand) MaxSize() int {
	return h.removedStartIdx
}

func (h Hand) CurrentSize() int {
	return h.removedStartIdx
}

func (h *Hand) Reset() {
	h.playedStartIdx = h.removedStartIdx
}

func (h *Hand) Restart() {
	h.playedStartIdx = 4
	h.removedStartIdx = 4
}

func (h *Hand) Remove(i int) error {
	if i < 0 || i >= h.removedStartIdx {
		return ErrOutOfBound
	}

	// decrement the valid card index
	h.removedStartIdx--
	h.swap(i, h.removedStartIdx)
	return nil
}

func (h *Hand) swap(i, j int) {
	h.cards[i], h.cards[j] = h.cards[j], h.cards[i]
}

func (h *Hand) CanPlay() bool {
	return h.MaxSize() != 0 || h.playedStartIdx != 0
}

func (h Hand) Play(i int) (CardType, error) {
	if !h.CanPlay() || i >= h.playedStartIdx {
		return CardTypeSkull, ErrCantPlay
	}

	h.playedStartIdx--
	h.swap(i, h.playedStartIdx)
	return h.cards[h.playedStartIdx], nil
}
