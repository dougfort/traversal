# Change Log

## v0.2.0

### Changed

- Changed the definition of the `Traversal` object
  - Make the interna1 representation a slice to represent a JSON Array
  - Fix methods and tests to support the new definition

## v0.1.1

### Added

- Ability to get raw message from a `json.RawMessage` type
- Ability to retrieve a slice from the next traversal step using `ArraySlice()`

## v0.1.0

Initial release. Basic functionality with simple tools.
