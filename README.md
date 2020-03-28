# Traverse JSON Data

## Synopsis

Traverse is a `go` language utility for 'traversing' a block of JSON encoded text
and returning a selected piece of content.

Starting with a `go` [RawMessage](<https://golang.org/pkg/encoding/json/#RawMessage>),
traverse subdivides `RawMessage`s until we have the one we want.

We specify the traversal using the [Builder Pattern]( <https://en.wikipedia.org/wiki/Builder_pattern>).
The [ZeroLog](https://github.com/rs/zerolog) logger also makes excellent use of this pattern.

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
        m, err := tr.GetMapFromRawMessage(r)
        if err != nil {
            return false
        }
        n, err := tr.GetStringFromRawMessage(m["name"])
        if err != nil {
            return false
        }
        return n == "BMW"
    }

    tr.Start([]byte(data)).
            ObjectKey("cars").
            ArrayPredicate(predicate).
            ObjectKey("models").
            End(os.stdout)
```

This example prints stdout:

```bash
[ "320", "X3", "X5" ]
```