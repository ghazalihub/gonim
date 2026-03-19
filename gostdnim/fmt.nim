type Any* = any
proc Println*(args: varargs[string, `$`]) = discard
proc Printf*(format: string, args: varargs[string, `$`]) = discard
proc Sprintf*(format: string, args: varargs[string, `$`]): string = ""
proc Print*(args: varargs[string, `$`]) = discard
