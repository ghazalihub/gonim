type Builder* = object
proc WriteString*(b: var Builder, s: string) = discard
proc String*(b: Builder): string = ""
