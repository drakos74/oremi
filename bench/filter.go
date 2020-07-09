package bench

type FilterType int

const (
	IN FilterType = iota + 1
	OUT
)

type Filter struct {
	Type   FilterType
	filter map[string]float64
}

func Include(f map[string]float64) Filter {
	return Filter{
		Type:   IN,
		filter: f,
	}
}

func Exclude(f map[string]float64) Filter {
	return Filter{
		Type:   OUT,
		filter: f,
	}
}
