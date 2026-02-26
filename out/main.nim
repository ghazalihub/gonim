import tables
proc main() =
  var m: Table[string, int] = initTable[string, int]()
  m["hello"] = 42
  echo(m["hello"])
