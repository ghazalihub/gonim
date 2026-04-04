type Any* = any
proc Println*(args: varargs[string, `$`]) =
  for arg in args:
    stdout.write(arg)
  stdout.write("\n")

proc Printf*(format: string, args: varargs[string, `$`]) = discard
proc Sprintf*(format: string, args: varargs[string, `$`]): string = ""
proc Print*(args: varargs[string, `$`]) =
  for arg in args:
    stdout.write(arg)
