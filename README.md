# fm [![Go Report Card](https://goreportcard.com/badge/github.com/enocom/fm)](https://goreportcard.com/report/github.com/enocom/fm) [![GoDoc](https://godoc.org/github.com/enocom/fm?status.svg)](https://godoc.org/github.com/enocom/fm)

The letters `fm` are short for the Chinese word _fangmao_ 仿冒, which literally means "to imitate and obscure", or "counterfeit." It is also a tool written in Go for generating spy implementations of interfaces.

*Note*: the use of the word "spy" is deliberate. See [here](https://8thlight.com/blog/uncle-bob/2014/05/14/TheLittleMocker.html) for more.

## Background

Writing a spy generator is the "Hello, World!" of AST parsing and generating in Go. There are many full featured libraries that do the same thing and better. For example, see [Counterfeiter](https://github.com/maxbrunsfeld/counterfeiter), [Hel](https://github.com/nelsam/hel), or [GoMock](https://github.com/golang/mock). The code here represents my own minimalist approach to the problem of generating test doubles.

## TODO

- add unit tests
- support embedded interfaces
- update import paths when necessary
- add comment to generated code identifying it as such
- add ci
