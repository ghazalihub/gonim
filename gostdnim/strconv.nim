import strutils

proc Itoa*(i: int): string = $i
proc Atoi*(s: string): (int, ref Exception) =
  try:
    return (parseInt(s), nil)
  except:
    return (0, (ref Exception)(msg: "invalid int"))
