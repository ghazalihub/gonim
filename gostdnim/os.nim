type MyError* = ref object of CatchableError
  msg*: string
type File* = ref object
proc Create*(name: string): (File, MyError) = (nil, nil)
proc Remove*(name: string): MyError = nil
proc Close*(f: File) = discard
