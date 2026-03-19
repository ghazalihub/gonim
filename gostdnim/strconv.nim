type MyError* = ref object of CatchableError
  msg*: string
proc Itoa*(i: int): string = ""
proc Atoi*(s: string): (int, MyError) = (0, nil)
