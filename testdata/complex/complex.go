package complex

type Point struct {
	X, Y int
}

func (p Point) Sum() int {
	return p.X + p.Y
}

type Summer interface {
	Sum() int
}

func GetSum(s Summer) int {
	return s.Sum()
}

func Main() {
	p := Point{X: 10, Y: 20}
	s := GetSum(p)
	println(s)
}
