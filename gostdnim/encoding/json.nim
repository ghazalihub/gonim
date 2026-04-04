proc Marshal*(v: any): (seq[byte], ref Exception) = (@[], nil)
proc Unmarshal*(data: seq[byte], v: any) = discard
