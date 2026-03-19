type MyError* = ref object of CatchableError
  msg*: string
proc WriteString*(w: any, s: string): (int, MyError) = (0, nil)
