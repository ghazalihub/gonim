proc Create*(name: string): (File, error) = (File(), nil)
proc Remove*(name: string): error = nil
type File* = object
proc Close*(f: File) = discard
