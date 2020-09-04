# JSONC - JSON with Comments
[![GoDoc](https://godoc.org/github.com/msolo/jsonc?status.svg)](https://godoc.org/github.com/msolo/jsonc)

JSONC allows parsing of a chunk of JSON that contain helpful comemnts.


Sample JSON content:
```java
/*
Preamble with fanfare.
*/

{
  // Line comment.
  "x": "a string", // Trailing comment.
  "y": 1.0,
  "z": null,
  "array": [],
  "dict": {}  // Wish we could have a trailing comma here.
}
// Postamble.
```

Sample usage:
```go
v := make(map[string]interface{})
f, _ := os.Open("sample.jsonc")
dec := jsonc.NewDecoder(f)
if err := dec.Decode(&v); err != nil {
  return err
}
```
