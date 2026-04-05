import strutils
import gostdnim/builtin

method Error*(e: ref Exception): string {.base.} =
  if e.isNil: "nil" else: e.msg

proc Println*(args: varargs[string, `$`]): (int, ref Exception) =
  var s = ""
  for i, arg in args:
    if i > 0: s.add(" ")
    s.add($arg)
  stdout.write(s & "\n")
  return (s.len + 1, nil)

proc Print*(args: varargs[string, `$`]): (int, ref Exception) =
  var s = ""
  for arg in args:
    s.add($arg)
  stdout.write(s)
  return (s.len, nil)

proc Sprintf*(format: string, args: varargs[string, `$`]): string =
  var res = ""
  var argIdx = 0
  var i = 0
  while i < format.len:
    if format[i] == '%' and i + 1 < format.len:
      let spec = format[i+1]
      if spec == '%':
        res.add('%')
      elif argIdx < args.len:
        res.add($args[argIdx])
        argIdx.inc
      i += 2
    else:
      res.add(format[i])
      i.inc
  return res

proc Printf*(format: string, args: varargs[string, `$`]): (int, ref Exception) =
  let s = Sprintf(format, args)
  stdout.write(s)
  return (s.len, nil)

proc Errorf*(format: string, args: varargs[string, `$`]): ref Exception =
  (ref Exception)(msg: Sprintf(format, args))
