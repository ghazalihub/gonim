type AnyX* = auto
proc isNilX*(e: auto): bool =
  when e is ref:
    e == nil
  elif e is ptr:
    e == nil
  else:
    false

template handleRecover*(body: untyped) =
  try:
    body
  except Exception as eX:
    discard

proc `$`*(e: ref Exception): string =
  if e.isNil: "nil" else: e.msg
