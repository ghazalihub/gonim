type MyError* = ref object of CatchableError
  msg*: string
proc Marshal*(v: auto): (seq[byte], MyError) = (newSeq[byte](), nil)
proc Unmarshal*(data: seq[byte], v: auto): MyError = nil
