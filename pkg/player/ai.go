package player

import (
	"errors"
	"math/rand/v2"

	"github.com/gb714us/skull/internal/utils"
	"github.com/gb714us/skull/pkg/skull"
)

var _ skull.Player = (*AIPlayer)(nil)

type AIPlayer struct {
	r  *rand.Rand
	h  skull.Hand
	rd randomDiscarder
	st skull.Color
}

type randomDiscarder struct {
	r *rand.Rand
}

func (rd randomDiscarder) Select(numOfCards int) int {
	return utils.GetRandomBoundedInt(rd.r, numOfCards)
}

func NewAIPlayer(r *rand.Rand, st skull.Color) *AIPlayer {
	return &AIPlayer{
		r:  r,
		h:  skull.NewHand(),
		rd: randomDiscarder{r},
		st: st,
	}
}

func (sp AIPlayer) PlayCard() (skull.CardType, error) {
	if !sp.h.CanPlay() {
		return skull.CardTypeSkull, skull.ErrCantPlay
	}

	cardIdx := utils.GetRandomBoundedInt(sp.r, sp.h.CurrentSize())
	return sp.h.Play(cardIdx)
}

func (sp AIPlayer) GetDiscarder() skull.Discarder {
	return sp.rd
}

func (sp AIPlayer) DiscardCard(d skull.Discarder) error {
	if sp.h.MaxSize() == 0 {
		return errors.New("no more cards to discard")
	}

	i := d.Select(sp.h.CurrentSize())
	return sp.h.Remove(i)
}

func (sp AIPlayer) Reset() {
	sp.h.Reset()
}

func (sp AIPlayer) Restart() {
	sp.h.Restart()
}

// Bid implements Player.
func (sp *AIPlayer) Bid(highestBid int, playedCards map[skull.Color]int) int {
	// get random bool to see if player is sitting out
	// more weighted to sit out
	if sp.r.IntN(3) != 1 {
		return 0
	}

	// bids atleast 1 higher than highestBid
	// adds a random number between 0 and 1 so its not only 1 number.
	return highestBid + (sp.r.IntN(2) + 1)
}

// HasDiscardedAllCards implements Player.
func (sp *AIPlayer) CanPlayCard() bool {
	return sp.h.CanPlay()
}

// HasPlayedCards implements Player.
func (sp *AIPlayer) HasPlayedCards() bool {
	return sp.h.CurrentSize() < sp.h.MaxSize()
}

// SelectPlayer implements Player.
func (sp *AIPlayer) SelectOpposingColor(playedCards map[skull.Color]int) skull.Color {
	// TODO: Build heuristics on how to select player since the number of played cards can change
	// the reasoning of who to choose

	// select random player
	keys := make([]skull.Color, 0, len(playedCards))
	for k := range playedCards {
		if k == sp.st {
			continue
		}
		keys = append(keys, k)
	}

	return keys[sp.r.IntN(len(keys))]
}

// GetSkullType implements Player.
func (sp *AIPlayer) GetSkullColor() skull.Color {
	return sp.st
}

// AnnounceBid implements skull.Player.
func (sp *AIPlayer) AnnounceBid() bool {
	// 20% of announcing the bid
	return sp.r.IntN(5) == 0
}
