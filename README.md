# Traverse JSON Data

Use Case: process a subset JSON encoded text.

The go RawMessage <https://golang.org/pkg/encoding/json/#RawMessage> enables stepping through JSON without parsing all of it.

This package facilitates the stepping by implementing a form of the Builder Pattern <https://en.wikipedia.org/wiki/Builder_pattern>
