type File* = ref object
proc Create*(name: string): (File, ref Exception) = (nil, nil)
proc Close*(f: File) = discard
proc Remove*(name: string) = discard
