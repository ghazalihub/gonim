package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"
	"unsafe"
)

// ================== Constants ==================

const (
	UntypedConst = 42
	TypedConst   int = 100
)

const (
	Red = iota
	Green
	Blue
)

// ================== Type Declarations ==================

type MyInt int                 // new type
type AliasInt = int            // alias
type Stringer interface {      // interface
	String() string
}

type Describer interface {
	Describe() string
}

type Person struct {
	Name string
	Age  int
}

type Employee struct {
	Person // embedded struct
	ID     int
}

// Custom error
type MyError struct {
	Msg string
}

func (e MyError) Error() string {
	return "MyError: " + e.Msg
}

// ================== Methods ==================

func (p Person) String() string {
	return fmt.Sprintf("Person(Name=%s, Age=%d)", p.Name, p.Age)
}

func (p *Person) Grow() {
	p.Age++
}

func (e Employee) Describe() string {
	return fmt.Sprintf("Employee(ID=%d, Name=%s)", e.ID, e.Name)
}

// ================== Generics ==================

type Number interface {
	int | float64
}

func Add[T Number](a, b T) T {
	return a + b
}

// ================== Functions ==================

func variadicSum(nums ...int) int {
	total := 0
	for _, v := range nums {
		total += v
	}
	return total
}

func namedReturn(a int) (result int) {
	result = a * 2
	return
}

func multipleReturn() (int, string) {
	return 7, "seven"
}

func mayFail(flag bool) error {
	if flag {
		return MyError{"something went wrong"}
	}
	return nil
}

func panicExample() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered from:", r)
		}
	}()
	panic("intentional panic")
}

// Higher order
func apply(f func(int) int, x int) int {
	return f(x)
}

// ================== init ==================

func init() {
	fmt.Println("Init function executed")
}

// ================== MAIN ==================

func main() {
	_ = errors.New("unused")
	// ===== Variables =====
	var x int = 10
	y := 20
	var z MyInt = 30
	_ = z

	// ===== Pointers =====
	ptr := &x
	*ptr = 50

	// ===== Arrays =====
	arr := [3]int{1, 2, 3}

	// ===== Slices =====
	slice := []int{4, 5, 6}
	slice = append(slice, 7)

	// ===== Maps =====
	m := map[string]int{"one": 1, "two": 2}
	m["three"] = 3
	delete(m, "two")

	// ===== Struct & Methods =====
	p := Person{"Alice", 25}
	p.Grow()

	e := Employee{
		Person: Person{"Bob", 30},
		ID:     123,
	}

	// ===== Interface usage =====
	var s Stringer = p
	fmt.Println(s.String())

	var d Describer = e
	fmt.Println(d.Describe())

	// ===== Type Assertion =====
	if v, ok := s.(Person); ok {
		fmt.Println("Type asserted:", v.Name)
	}

	// ===== Switch =====
	switch x {
	case 50:
		fmt.Println("x is 50")
	default:
		fmt.Println("x unknown")
	}

	// Type switch
	var any interface{} = 42
	switch v := any.(type) {
	case int:
		fmt.Println("int:", v)
	case string:
		fmt.Println("string:", v)
	}

	// ===== For loops =====
	for i := 0; i < 3; i++ {
		fmt.Print(i)
	}
	fmt.Println()

	for i := range arr {
		fmt.Print(arr[i])
	}
	fmt.Println()

	for k, v := range m {
		fmt.Println(k, v)
	}

	// ===== Labels =====
Outer:
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if j == 2 {
				break Outer
			}
		}
	}

	// ===== Defer =====
	defer fmt.Println("Deferred call")

	// ===== Closures =====
	counter := func() func() int {
		count := 0
		return func() int {
			count++
			return count
		}
	}()

	fmt.Println(counter())
	fmt.Println(counter())

	// ===== Generics =====
	fmt.Println(Add(5, 3))
	fmt.Println(Add(2.5, 1.5))

	// ===== Variadic =====
	fmt.Println(variadicSum(1, 2, 3, 4))

	// ===== Multiple return =====
	num, word := multipleReturn()
	fmt.Println(num, word)

	// ===== Error Handling =====
	if err := mayFail(true); err != nil {
		fmt.Println(err)
	}

	// ===== Panic & Recover =====
	panicExample()

	// ===== Reflection =====
	t := reflect.TypeOf(p)
	fmt.Println("Type:", t.Name())

	// ===== JSON =====
	data, _ := json.Marshal(p)
	fmt.Println(string(data))

	var newP Person
	json.Unmarshal(data, &newP)

	// ===== String Builder =====
	var builder strings.Builder
	builder.WriteString("Hello ")
	builder.WriteString("World")
	fmt.Println(builder.String())

	// ===== Sorting =====
	ints := []int{5, 2, 8}
	sort.Ints(ints)
	fmt.Println(ints)

	// ===== File I/O =====
	file, _ := os.Create("temp.txt")
	io.WriteString(file, "sample")
	file.Close()
	os.Remove("temp.txt")

	// ===== Time =====
	fmt.Println("Now:", time.Now().Format(time.RFC3339))

	// ===== Rune =====
	r := '世'
	fmt.Println("Rune:", r, string(r))

	// ===== Unsafe =====
	size := unsafe.Sizeof(x)
	fmt.Println("Size of x:", size)

	// ===== Anonymous Struct =====
	temp := struct {
		Field string
	}{"value"}
	fmt.Println(temp.Field)

	// ===== Function Type =====
	var f func(int) int = func(n int) int { return n * n }
	fmt.Println(f(5))

	// ===== Method Expression =====
	method := Person.String
	fmt.Println(method(p))

	// ===== strconv =====
	str := strconv.Itoa(123)
	num2, _ := strconv.Atoi(str)
	fmt.Println(str, num2)

	// ===== Named return =====
	fmt.Println(namedReturn(5))

	// ===== Blank Identifier =====
	_, _ = y, slice
}
