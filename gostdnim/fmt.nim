proc Println*(args: varargs[untyped]) = discard
proc Printf*(format: string, args: varargs[untyped]) = discard
proc Sprintf*(format: string, args: varargs[untyped]): string = ""
proc Print*(args: varargs[untyped]) = discard
