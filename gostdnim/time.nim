type Time* = object
proc Now*(): Time = Time()
proc Format*(t: Time, layout: string): string = ""
const RFC3339* = ""
