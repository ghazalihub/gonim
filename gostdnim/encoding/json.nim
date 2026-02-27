proc Marshal*(v: any): (seq[byte], error) = (newSeq[byte](), nil)
proc Unmarshal*(data: seq[byte], v: any): error = nil
