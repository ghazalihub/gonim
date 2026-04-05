import std/os
import streams

type File* = ref object
  stream*: Stream

proc Create*(name: string): (File, ref Exception) =
  try:
    let f = File(stream: newFileStream(name, fmWrite))
    return (f, nil)
  except Exception as e:
    return (nil, (ref Exception)(msg: e.msg))

proc Close*(f: File) =
  if not f.isNil and not f.stream.isNil:
    f.stream.close()

proc Remove*(name: string) =
  try:
    removeFile(name)
  except:
    discard
