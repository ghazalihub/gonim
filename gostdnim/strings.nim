type Builder* = object
  s: string

proc WriteString*(b: var Builder, s: string): (int, ref Exception) =
  b.s.add(s)
  return (s.len, nil)

proc String*(b: Builder): string =
  b.s
