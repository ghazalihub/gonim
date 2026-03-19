type MyError* = ref object of CatchableError

proc New*(msg: string): MyError =
  new(result)
