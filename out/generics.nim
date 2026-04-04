import gostdnim/builtin

import gostdnim/fmt

proc Identity*[T](x: T): T =
  return x

type Container*[T] = object
    Value*: T


proc Main*() =
  var x = Identity[int](42)
  echo(x)
  var c = Container[string](Value: "hello")
  echo(c.Value)
