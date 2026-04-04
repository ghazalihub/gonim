type AnyX* = auto
proc isNil*(e: auto): bool =
  when e is ref:
    e == nil
  else:
    false

proc Println*(args: varargs[string, `$`]) =
  for arg in args:
    stdout.write($arg)
  stdout.write("\n")

proc Print*(args: varargs[string, `$`]) =
  for arg in args:
    stdout.write($arg)
