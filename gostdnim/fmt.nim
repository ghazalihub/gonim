proc Println*(args: varargs[string, `$`]) =
  for arg in args:
    stdout.write($arg)
  stdout.write("\n")

proc Print*(args: varargs[string, `$`]) =
  for arg in args:
    stdout.write($arg)

proc Sprintf*(format: string, args: varargs[string, `$`]): string = ""
proc Errorf*(format: string, args: varargs[string, `$`]): ref Exception = (ref Exception)(msg: "")
