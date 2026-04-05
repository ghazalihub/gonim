import std/typetraits

type Type* = object
  name*: string

template TypeOf*(v: any): Type =
  Type(name: v.typeof.name)

proc Name*(t: Type): string = t.name
