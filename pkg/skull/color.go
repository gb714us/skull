package skull

type Color int

const (
	ColorRed Color = 1 << iota
	ColorBlue
	ColorGreen
	ColorYellow
	ColorPurple
	ColorGray
)

func (s Color) String() string {
	switch s {
	case ColorRed:
		return "red"
	case ColorBlue:
		return "blue"
	case ColorGreen:
		return "green"
	case ColorYellow:
		return "yellow"
	case ColorPurple:
		return "purple"
	case ColorGray:
		return "brown"
	default:
		return "unknown"
	}
}
