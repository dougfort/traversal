# Traverse JSON Data

[![godoc](http://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](https://godoc.org/github.com/dougfort/traversal)

## Synopsis

Traverse is a `go` language utility for 'traversing' a block of JSON encoded text
and returning selected content.

Starting with a `go` [RawMessage](<https://golang.org/pkg/encoding/json/#RawMessage>),
traverse subdivides `RawMessage`s until we have the one we want.

We specify the traversal using the [Builder Pattern]( <https://en.wikipedia.org/wiki/Builder_pattern>).
The [ZeroLog](https://github.com/rs/zerolog) logger also makes excellent use of this pattern.

## Components

A traversal is **composed** of a chain chain of components. It must start with `Start` and
end with `End`. Otherwise you can snap the components together in whatever order fits the JSON
text.

* `Start` - initiize the traversl state
* `End` - extract the traversal state, or report error
* `ObjectKey` - pull an item out of a JSON Object, and make that the new state
* `ArraySingleton` - select the only item in a 1 item JSON Array and make that the new state
* `ArraySlice` - select a JSON array as a slice and make that the new state
* `ArrayPredicate` - select an item from a JSON Array based on a predicate function and make that the new state
* `Selector` - use a selector function to select whatever you want from the current state and make that the new state

## Error Handling

An error that occurs anywhere in the traversal chain will be reported as the return value from `End`

## Helper Functions

Traversal components use helper functions to manipulate `json.RawMessage` objects. These are available to you so you can write
your own components.

* `GetStringFromRawMessage`
* `GetBoolFromRawMessage`
* `GetInt32FromRawMessage`
* `GetSliceFromRawMessage`
* `GetMapFromRawMessage`
* `GetMsgFromRawMessage`

## Example

example copied from [w3schools.com](https://www.w3schools.com/js/js_json_arrays.asp)

```json
{
  "name":"John",
  "age":30,
  "cars": [
    { "name":"Ford", "models":[ "Fiesta", "Focus", "Mustang" ] },
    { "name":"BMW", "models":[ "320", "X3", "X5" ] },
    { "name":"Fiat", "models":[ "500", "Panda" ] }
  ]
}
```

 We want to extract an Array of the BMW models which John owns.

```golang
    predicate := func(r json.RawMessage) bool {
        m, err := traversal.GetMapFromRawMessage(r)
        if err != nil {
            return false
        }
        n, err := traversal.GetStringFromRawMessage(m["name"])
        if err != nil {
            return false
        }
        return n == "BMW"
    }

    traversal.Start([]byte(data)).
            ObjectKey("cars").
            ArrayPredicate(predicate).
            ObjectKey("models").
            End(os.stdout)
```

This example prints to stdout:

```bash
[ "320", "X3", "X5" ]
```
