# goescape

Performance work in Go often entails avoiding heap allocations in hot parts of
the code. Once a heap allocation has been avoided (i.e. moved to the stack),
one has to worry about the next person accidentally undoing the optimization;
the rules of escape analysis are far from obvious and innocuous changes can
have drastic effects.

The escape analysis can be printed via `go build -gcflags=-m`, but this is
far from being comfortable. goescape provides a library around this that
can provide linting by allowing annotating variables which need to remain
on the stack in a fairly non-intrusive way:

```go
notAllowedToEscape := struct {
    goescape.Stack
    YourType
} {
    YourType: f(),
}
```

Once annotated, the code base can be kept regression free by adding a linter
of the form

```go
fs, err := goescape.Lint("./pkg/...", ".")
if err != nil {
    t.Fatal(err)
}
for _, f := range fs {
    // Output: some/file.go:12: unexpectedly escaped: notAllowedToEscape
    t.Error(f)
}
```

You can also use the `goescape` command, for example

```
$ go run cmd/goescape/main.go ./examples/...
examples/buffers.go:9: unexpectedly escaped: h1
examples/buffers.go:17: unexpectedly escaped: h2

exit status 1
```
