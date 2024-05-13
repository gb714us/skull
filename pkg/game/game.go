package game

import (
	"container/ring"
	"fmt"
	"math/rand/v2"
	"strings"

	"github.com/gb714us/skull/internal/stack"
	"github.com/gb714us/skull/pkg/skull"
)

type Game struct {
	wins              map[skull.Color]int
	lastWinningBidder *ring.Ring
	playerBufferMap   map[skull.Color]skull.Player

	board      map[skull.Color]*stack.Stack[skull.CardType]
	eliminated map[skull.Color]bool
}

func NewGame() *Game {
	return &Game{
		wins:              map[skull.Color]int{},
		lastWinningBidder: nil,
		playerBufferMap:   map[skull.Color]skull.Player{},
		board:             map[skull.Color]*stack.Stack[skull.CardType]{},
		eliminated:        map[skull.Color]bool{},
	}
}

func (g *Game) AddPlayer(p skull.Player) {
	if _, ok := g.playerBufferMap[p.GetSkullColor()]; ok {
		panic("player already exists. cannot continue game")
	}

	// add plauer now, but build ring once the game begins
	g.playerBufferMap[p.GetSkullColor()] = p
	g.board[p.GetSkullColor()] = stack.New[skull.CardType]()
}

func (g *Game) Start() skull.Player {
	// create eliminated
	playerCount := len(g.playerBufferMap)
	startingPlayerIdx := rand.IntN(len(g.playerBufferMap))
	var i = 0
	currentRing := ring.New(playerCount)
	for color, player := range g.playerBufferMap {
		currentRing.Value = player
		// init elimination map
		g.eliminated[color] = false
		if i == startingPlayerIdx {
			g.lastWinningBidder = currentRing
		}

		currentRing = currentRing.Next()
		i++
	}

	currentPlayer := currentRing.Value.(skull.Player)
	for g.wins[currentPlayer.GetSkullColor()] != 2 {
		highestBidder, win := g.startRound()

		g.lastWinningBidder = highestBidder
		currentPlayer = highestBidder.Value.(skull.Player)
		if win {
			// mark the win if the player won.
			g.wins[currentPlayer.GetSkullColor()]++
		}

		g.resetBoard()
		if !currentPlayer.CanPlayCard() && !win {
			// mark the player as eliminated
			g.eliminated[currentPlayer.GetSkullColor()] = true
			g.lastWinningBidder = g.lastWinningBidder.Next()
		}
	}

	return currentPlayer
}

// runs a single round of the game and returns the winner
func (g *Game) startRound() (*ring.Ring, bool) { // skull.Player {
	current := g.lastWinningBidder
	// initial card setting loop
	current.Do(func(c any) {
		player := c.(skull.Player)
		if eliminated := g.eliminated[player.GetSkullColor()]; eliminated {
			// this player is eliminated from the game
			return
		}

		card, err := player.PlayCard()
		if err != nil {
			// this should never happen.. but just in case
			panic(err)
		}

		// add card to board
		g.board[player.GetSkullColor()].Push(card)
		g.printBoard(false)
	})

	// playing card loop
	for {
		player := current.Value.(skull.Player)
		if eliminated := g.eliminated[player.GetSkullColor()]; eliminated {
			// skip over this player
			continue
		}

		// if player can't play a card, break out of the loop and
		// start the bidding phase
		if !player.CanPlayCard() || player.AnnounceBid() {
			break
		}

		card, err := player.PlayCard()
		if err != nil {
			// this should never happen.. but just in case
			panic(err)
		}

		// add card to board
		g.board[player.GetSkullColor()].Push(card)
		g.printBoard(false)
	}

	// get final state of the board
	colorCardCount := g.getCardsCountByPlayer()

	// bid loop
	highestBid := 0
	pushed := map[skull.Color]bool{}
	remainingBidders := len(g.playerBufferMap) - len(g.eliminated)
	for remainingBidders > 1 {
		player := current.Value.(skull.Player)
		// if player can't play a card, break out of the loop and
		// start the bidding phase
		if _, ok := pushed[player.GetSkullColor()]; ok || g.eliminated[player.GetSkullColor()] {
			continue
		}

		bid := player.Bid(highestBid, colorCardCount)
		if bid == 0 {
			// player is sitting out
			pushed[player.GetSkullColor()] = true
			remainingBidders--
		} else {
			highestBid = bid
			current = current.Next()
		}
	}

	win := true
	highestBiddingPlayer := current.Value.(skull.Player)
	// card selection loop
	for highestBid > 0 && win {
		color := highestBiddingPlayer.SelectOpposingColor(colorCardCount)
		colorCardStack := g.board[color]

		cardType, err := colorCardStack.Pop()
		if err != nil {
			panic(err)
		}

		fmt.Printf("Player %s popped %s's %s\n", highestBiddingPlayer.GetSkullColor(), color, cardType)
		g.printBoard(true)
		if colorCardStack.Size() == 0 {
			// removed so they no longer show up on the game boar
			delete(colorCardCount, color)
		}

		if cardType == skull.CardTypeSkull {
			win = false
			// allow the chosen color player to discard
			chosenColorPlayerDiscarder := g.playerBufferMap[color].GetDiscarder()
			highestBiddingPlayer.DiscardCard(chosenColorPlayerDiscarder)
		} else {
			// let the player keep discarding
			highestBid--
		}
	}

	return current, win
}

func (g *Game) resetBoard() {
	for color, stack := range g.board {
		// clear the board stack
		stack.Clear()

		// reset players hand without giving back
		// those that have been discarded
		g.playerBufferMap[color].Reset()
	}
}

func (g *Game) getCardsCountByPlayer() map[skull.Color]int {
	result := map[skull.Color]int{}
	for color, stack := range g.board {
		if stack.Size() > 0 {
			result[color] = stack.Size()
		}
	}
	return result
}

func (g *Game) printBoard(showCardType bool) {
	for color, stack := range g.board {
		fmt.Printf("%s: ", color)
		fmt.Printf("%s\n", strings.Repeat("*", stack.Size()))
	}

	fmt.Println("+++++++++")
}
