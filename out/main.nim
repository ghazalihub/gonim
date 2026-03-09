import tables
import gostdnim/encoding/json
import gostdnim/errors
import gostdnim/fmt
import gostdnim/io
import gostdnim/os
import gostdnim/reflect
import gostdnim/sort
import gostdnim/strconv
import gostdnim/strings
import gostdnim/time
import gostdnim/unsafe
const UntypedConst: untyped int = 42
const TypedConst: int = 100
const Red: untyped int = 0
const Green: untyped int = 1
const Blue: untyped int = 2
type MyInt* = int
type AliasInt* = command-line-arguments.AliasInt
type Stringer* = concept x
    x.String() is string

type Describer* = concept x
    x.Describe() is string

type Person* = object
    Name*: string
    Age*: int

type Employee* = object
    Person*: Person
    ID*: int

type MyError* = object
    Msg*: string

proc Error*(e: MyError): string =
  return ("MyError: " + e.Msg)

proc String*(p: Person): string =
  return fmt.Sprintf("Person(Name=%s, Age=%d)", p.Name, p.Age)

proc Grow*(p: ref Person) =
  discard

proc Describe*(e: Employee): string =
  return fmt.Sprintf("Employee(ID=%d, Name=%s)", e.ID, e.Name)

type Number* = concept x

proc Add*[T](a: T, b: T): T =
  return (a + b)

proc variadicSum(nums: varargs[int]): int =
  var total: int = 0
  for _, v in nums:
  total += v
  return total

proc namedReturn(a: int): int =
  result = (a * 2)
  return result

proc multipleReturn(): (int, string) =
  return 7, "seven"

proc mayFail(flag: bool): ref Exception =
  if flag:
    return MyError("something went wrong")
  return nil

proc panicExample() =
  defer: proc(...) = discard()
  raise("intentional panic")

proc apply(f: proc(int): int, x: int): int =
  return f(x)

proc init() =
  fmt.Println("Init function executed")

proc main() =
  discard errors.New("unused")
  var x: int = 10
  var y: int = 20
  var z: MyInt = 30
  discard z
  var ptr: ref int = addr(x)
  ptr[] = 50
  var arr: array[3, int] = array[3, int](1, 2, 3)
  var slice: seq[int] = seq[int](4, 5, 6)
  slice = slice.add(7)
  var m: Table[string, int] = Table[string, int]("one": 1, "two": 2)
  m["three"] = 3
  del(m, "two")
  var p: Person = Person("Alice", 25)
  p.Grow()
  var e: Employee = Employee(Person: Person("Bob", 30), ID: 123)
  var s: Stringer = p
  fmt.Println(s.String())
  var d: Describer = e
  fmt.Println(d.Describe())
  if ok:
    fmt.Println("Type asserted:", v.Name)
  case x:
  of 50:
  fmt.Println("x is 50")
  else:
  fmt.Println("x unknown")
  var any: concept = 42
  if any is int:
  fmt.Println("int:", v)
elif any is string:
  fmt.Println("string:", v)
  discard # unknown stmt
  fmt.Println()
  for i in arr:
  fmt.Print(arr[i])
  fmt.Println()
  for k, v in m:
  fmt.Println(k, v)
  defer: fmt.Println("Deferred call")
  var counter: proc(): int = proc(...) = discard()
  fmt.Println(counter())
  fmt.Println(counter())
  fmt.Println(Add(5, 3))
  fmt.Println(Add(2.5, 1.5))
  fmt.Println(variadicSum(1, 2, 3, 4))
  var (num, word): int = multipleReturn()
  fmt.Println(num, word)
  if (err != nil):
    fmt.Println(err)
  panicExample()
  var t: reflect.Type = reflect.TypeOf(p)
  fmt.Println("Type:", t.Name())
  var (data, _): seq[byte] = json.Marshal(p)
  fmt.Println(string(data))
  var newP: Person
  json.Unmarshal(data, addr(newP))
  var builder: strings.Builder
  builder.WriteString("Hello ")
  builder.WriteString("World")
  fmt.Println(builder.String())
  var ints: seq[int] = seq[int](5, 2, 8)
  sort.Ints(ints)
  fmt.Println(ints)
  var (file, _): ref os.File = os.Create("temp.txt")
  io.WriteString(file, "sample")
  file.Close()
  os.Remove("temp.txt")
  fmt.Println("Now:", time.Now().Format(time.RFC3339))
  var r: int32 = '世'
  fmt.Println("Rune:", r, string(r))
  var size: uintptr = unsafe.Sizeof(x)
  fmt.Println("Size of x:", size)
  var temp: tuple[Field: string] = tuple[Field: string]("value")
  fmt.Println(temp.Field)
  var f: proc(int): int = proc(...) = discard
  fmt.Println(f(5))
  var method: proc(Person): string = Person.String
  fmt.Println(method(p))
  var str: string = strconv.Itoa(123)
  var (num2, _): int = strconv.Atoi(str)
  fmt.Println(str, num2)
  fmt.Println(namedReturn(5))
  (_, _) = (y, slice)
