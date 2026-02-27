type Type* = object
proc TypeOf*(v: any): Type = Type()
proc Name*(t: Type): string = ""
