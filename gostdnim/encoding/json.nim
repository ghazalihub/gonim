import std/json

proc Marshal*(v: any): (seq[byte], ref Exception) =
  try:
    let s = $(%v)
    var b = newSeq[byte](s.len)
    for i in 0..<s.len: b[i] = byte(s[i])
    return (b, nil)
  except Exception as e:
    return (@[], (ref Exception)(msg: e.msg))

proc Unmarshal*[T](data: seq[byte], v: ptr T): ref Exception =
  try:
    var s = ""
    for b in data: s.add(char(b))
    v[] = parseJson(s).to(T)
    return nil
  except Exception as e:
    return (ref Exception)(msg: e.msg)
