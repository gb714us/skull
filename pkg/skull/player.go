package skull

// Discarder chooses a card for losing skulled player to discard
type Discarder interface {
	// Select returns an index to discard from a player hand from 0 to cardCount - 1
	Select(handCount int) int
}

type Player interface {
	PlayCard() (CardType, error)
	CanPlayCard() bool
	HasPlayedCards() bool

	GetDiscarder() Discarder
	DiscardCard(cs Discarder) error

	AnnounceBid() bool

	// Bid must return an int > highestBid
	// if Bid is 0, player is sitting out.
	Bid(highestBid int, playedCards map[Color]int) int

	// SelectPlayer returns an index to a  player  from 0 to numOfPlayers - 1
	SelectOpposingColor(playedCards map[Color]int) Color

	// Reset gives all played card back
	Reset()

	// Restart gives players all discarded and played cards back
	Restart()

	//GetSkullType returns the SkullType of the player
	GetSkullColor() Color
}
