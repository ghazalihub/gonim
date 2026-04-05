import gostdnim/os as gos
import streams

proc WriteString*(f: gos.File, s: string): (int, ref Exception) =
  try:
    if not f.isNil and not f.stream.isNil:
      f.stream.write(s)
      return (s.len, nil)
    return (0, (ref Exception)(msg: "invalid file"))
  except Exception as e:
    return (0, (ref Exception)(msg: e.msg))
