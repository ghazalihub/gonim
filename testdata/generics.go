package generics

func Identity[T any](x T) T {
	return x
}

type Container[T any] struct {
	Value T
}

func Main() {
	x := Identity[int](42)
	println(x)

	c := Container[string]{Value: "hello"}
	println(c.Value)
}
