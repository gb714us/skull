package main

import (
	"flag"
	"fmt"
	"math/rand/v2"

	"github.com/gb714us/skull/pkg/game"
	"github.com/gb714us/skull/pkg/player"
	"github.com/gb714us/skull/pkg/skull"
)

var aiPlayerCount int
var humanPlayerCount int

func init() {
	flag.IntVar(&humanPlayerCount, "human", 0, "human player count")
	flag.IntVar(&aiPlayerCount, "ai", 0, "ai player count")

	flag.Parse()

	if humanPlayerCount <= 0 {
		aiPlayerCount = 6
	}

	if humanPlayerCount > 6 {
		// default to 6 max.
		humanPlayerCount = 6
		aiPlayerCount = 0
	}

}

func main() {
	availableColors := map[skull.Color]struct{}{
		skull.ColorBlue:   {},
		skull.ColorGreen:  {},
		skull.ColorRed:    {},
		skull.ColorYellow: {},
		skull.ColorPurple: {},
		skull.ColorGray:   {},
	}

	skullGame := game.NewGame()

	// for i := 0; i < humanPlayerCount; i++ {
	// 	color := getRandomColor(availableColors)
	// 	skullGame.AddPlayer(skull.NewHumanPlayer(color))
	// }

	for i := 0; i < aiPlayerCount; i++ {
		color := getRandomColor(availableColors)
		delete(availableColors, color)
		skullGame.AddPlayer(player.NewAIPlayer(rand.New(rand.NewPCG(rand.Uint64(), rand.Uint64())), color))
	}

	winner := skullGame.Start()
	fmt.Printf("Winner: %s", winner.GetSkullColor())
}

func getRandomColor(ac map[skull.Color]struct{}) skull.Color {
	// write code that chooses random color from ac

	keys := make([]skull.Color, 0, len(ac))
	for k := range ac {
		keys = append(keys, k)
	}

	return keys[rand.IntN(len(keys))]
}
