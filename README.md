[![](https://godoc.org/github.com/fiatjaf/goleveldown?status.svg)](http://godoc.org/github.com/fiatjaf/goleveldown)
[![travis ci badge](https://travis-ci.org/fiatjaf/goleveldown.svg?branch=master)](https://travis-ci.org/fiatjaf/goleveldown)

## goleveldown

This package implements the [levelup](https://github.com/fiatjaf/levelup) interface.

## how to use

```go
package main

import (
    "github.com/fiatjaf/goleveldown"
    "github.com/fiatjaf/levelup"
    "github.com/fiatjaf/levelup/stringlevelup"
)

func main() {
    bdb := goleveldown.NewDatabase("/tmp/leveldownexample")
    defer bdb.Erase()

    db := stringlevelup.StringDB(bdb)

    fmt.Println("setting key1 to x")
    db.Put("key1", "x")
    res, _ := db.Get("key1")
    fmt.Println("setting key2 to 2")
    fmt.Println("res at key2: ", res)
    db.Put("key2", "y")
    res, _ = db.Get("key2")
    fmt.Println("res at key2: ", res)
    fmt.Println("deleting key1")
    db.Del("key1")
    res, _ = db.Get("key1")
    fmt.Println("res at key1: ", res)

    fmt.Println("batch")
    db.Batch([]levelup.Operation{
        stringlevelup.Put("key2", "w"),
        stringlevelup.Put("key3", "z"),
        stringlevelup.Del("key1"),
        stringlevelup.Put("key1", "t"),
        stringlevelup.Put("key4", "m"),
        stringlevelup.Put("key5", "n"),
        stringlevelup.Del("key3"),
    })
    res, _ = db.Get("key1")
    fmt.Println("res at key1: ", res)
    res, _ = db.Get("key2")
    fmt.Println("res at key2: ", res)
    res, _ = db.Get("key3")
    fmt.Println("res at key3: ", res)

    fmt.Println("reading all")
    iter := db.ReadRange(nil)
    for ; iter.Valid(); iter.Next() {
        fmt.Println("row: ", iter.Key(), " ", iter.Value())
    }
    fmt.Println("iter error: ", iter.Error())
    iter.Release()
}
```

if you don't call `stringlevelup.StringDB` on the object returned by `NewDatabase` you still can use the same methods, only replacing all `string` arguments with `[]byte` (returned values will also be bytes). I've just used the string approach here for readability, but it is slower.
