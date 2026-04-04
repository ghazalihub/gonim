import gostdnim/builtin

import gostdnim/strconv

import gostdnim/strings

import gostdnim/time

import tables

import gostdnim/encoding/json

import gostdnim/errors

import gostdnim/fmt

import gostdnim/io

import gostdnim/os

import gostdnim/reflect

import gostdnim/sort

import gostdnim/unsafe

const UntypedConst = 42

const TypedConst: int = 100

const Red = 0

const Green = 1

const Blue = 2

type MyInt* = int

type AliasInt* = int

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
  return ("MyError: " & e.Msg)

proc String*(p: Person): string =
  return fmt.Sprintf("Person(Name=%s, Age=%d)", p.Name, p.Age)

proc Grow*(p: ref Person) =
  p.Age.inc

proc Describe*(e: Employee): string =
  return fmt.Sprintf("Employee(ID=%d, Name=%s)", e.ID, e.Name)

type Number* = concept x


proc Add*[T](a: T, b: T): T =
  return (a + b)

proc variadicSum(nums: varargs[int]): int =
  var total: int = 0
  for i, v in nums:
    total += v
  return total

proc namedReturn(a: int): int =
  result = (a * 2)
  return result

proc multipleReturn(): (int, string) =
  return (7, "seven")

proc mayFail(flag: bool): ref Exception =
  if flag:
    return MyError("something went wrong")
  return nil

proc panicExample() =
  defer:
    (proc(): void =
      var r: Any = getCurrentExceptionMsg()
      if not r.isNil:
        fmt.Println("Recovered from:", r))()
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
  var ptr_: ref int = addr(x)
  ptr_[] = 50
  var arr: array[3, int] = [1, 2, 3]
  var slice: seq[int] = @ [4, 5, 6]
  slice.add(7)
  var m: Table[string, int] = {"one": 1, "two": 2}.toTable
  m["three"] = 3
  del(m, "two")
  var p: Person = Person("Alice", 25)
  p.Grow()
  var e: Employee = Employee(Person: Person("Bob", 30), ID: 123)
  var s: Stringer = p
  fmt.Println(s.String())
  var d: Describer = e
  fmt.Println(d.Describe())
  var (v, ok) = (if s is Person: (Person(s), true) else: (default(Person), false))
  if ok:
    fmt.Println("Type asserted:", v.Name)
  case x
  of 50:
    fmt.Println("x is 50")
  else:
    fmt.Println("x unknown")
  var any: Any = 42
  if any is int:
    let v = int(any)
    fmt.Println("int:", v)
  elif any is string:
    let v = string(any)
    fmt.Println("string:", v)
  var i: int = 0
  while (i < 3):
    fmt.Print(i)
    i.inc
  fmt.Println()
  for i, v in arr:
    fmt.Print(arr[i])
  fmt.Println()
  for k, v in m:
    fmt.Println(k, v)
  defer:
    fmt.Println("Deferred call")
  var counter: proc(): int = (proc(): proc(): int =
    var count: int = 0
    return (proc(): int =
      count.inc
      return count))()
  fmt.Println(counter())
  fmt.Println(counter())
  fmt.Println(Add(5, 3))
  fmt.Println(Add(2.5, 1.5))
  fmt.Println(variadicSum(1, 2, 3, 4))
  var (num, word) = multipleReturn()
  fmt.Println(num, word)
  var err: ref Exception = mayFail(true)
  if not err.isNil:
    fmt.Println(err)
  panicExample()
  var t: reflect.Type = reflect.TypeOf(p)
  fmt.Println("Type:", t.Name())
  block:
    var (data, _tmp1) = json.Marshal(p)
    discard _tmp1
  fmt.Println(string(data))
  var newP: Person
  json.Unmarshal(data, addr(newP))
  var builder: strings.Builder
  builder.WriteString("Hello ")
  builder.WriteString("World")
  fmt.Println(builder.String())
  var ints: seq[int] = @ [5, 2, 8]
  sort.Ints(ints)
  fmt.Println(ints)
  block:
    var (file, _tmp1) = os.Create("temp.txt")
    discard _tmp1
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
  var f: proc(int): int = (proc(n: int): int = (n * n))
  fmt.Println(f(5))
  var method_: proc(Person): string = Person.String
  fmt.Println(method_(p))
  var str: string = strconv.Itoa(123)
  block:
    var (num2, _tmp1) = strconv.Atoi(str)
    discard _tmp1
  fmt.Println(str, num2)
  fmt.Println(namedReturn(5))
  block:
    (_tmp0, _tmp1) = (y, slice)
    discard _tmp0
    discard _tmp1
