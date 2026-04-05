import times

type Time* = times.DateTime

proc Now*(): Time = times.now().utc

proc Format*(t: Time, layout: string): string =
  # Layout mapping is complex, let's just use a default for now
  t.format("yyyy-MM-dd'T'HH:mm:ss'Z'")

const RFC3339* = "yyyy-MM-ddTHH:mm:ssZ"
