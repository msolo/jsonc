# JSONC - JSON with Comments
[![GoDoc](https://godoc.org/github.com/msolo/jsonc?status.svg)](https://godoc.org/github.com/msolo/jsonc)

JSONC allows parsing chunks of JSON that contain helpful comments.

The original motivation was to have usable config files without having to resort to things like YAML that are staggeringly complex despite apparent simplicity.


## Sample JSONC Snippet
```java
/*
You can see that comments are safely valid in any sane place.

You can put a lengthy preamble or embed a poem if necessary.
*/

{
  // Line comment.
  "x": "a string", // Trailing comment.
  "y": 1.0,
  "z": null,
  "quoted-range": "/* this is not a comment *",
  "quoted-line": "// this is also not a comment",
  // "a": "value temporarily removed for debugging or idle curiosity",
  "array": [],
  "dict": {}  // Wish we could have a trailing comma here,
              // but that's a problem for a better parser.
}
// Postamble.
```

## Sample Usage in Go
```go
v := make(map[string]interface{})
f, _ := os.Open("sample.jsonc")
dec := jsonc.NewDecoder(f)
if err := dec.Decode(&v); err != nil {
  return err
}
```

## Command Line Tool

`jsonc-strip` is unix-esque tool to filter out comments so standard tools like `jq` are still useful.

```
go install github.com/msolo/jsonc/cmd/jsonc-strip

jsonc-strip < sample.jsonc

jsonc-strip < sample.jsonc | jq .x
"a string"

```
