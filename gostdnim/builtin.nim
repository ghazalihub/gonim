type Any* = ref object of RootObj
proc isNil*(e: any): bool =
  when e is ref:
    e == nil
  else:
    false
