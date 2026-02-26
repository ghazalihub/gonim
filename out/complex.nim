type Point* = object
    X*: int
    Y*: int

proc Sum*(p: Point): int =
  return (p.X + p.Y)

type Summer* = concept x
    x.Sum() is int

proc GetSum*(s: Summer): int =
  return s.Sum()

proc Main*() =
  var p = Point((X : 10), (Y : 20))
  var s = GetSum(p)
  echo(s)
