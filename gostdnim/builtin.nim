type MyError* = ref object of CatchableError
  msg*: string

proc isNil*(e: MyError): bool = e == nil
