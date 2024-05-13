package skull

type CardType int

const (
	CardTypeFlower = iota
	CardTypeSkull
)

func (c CardType) String() string {
	switch c {
	case CardTypeFlower:
		return "🌺"
	case CardTypeSkull:
		return "💀"
	default:
		return "unknown"
	}
}
