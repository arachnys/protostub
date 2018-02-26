# protostub

A tool for generating MyPy type stubs from a Protocol Buffer definition.

## Building

### Requirements
- Go
- make (optional)


### With go get
If you already have Go all setup in your `PATH`, then it is as simple as:

```
go get github.com/arachnys/protostub/cmd/protostub
```

### With Make
This approach might be best if you're less familiar with Go, and want it to 
*just work*. It requires no messing with `$GOPATH`.

```
git clone https://github.com/arachnys/protostub
cd protostub
make
```

The protostub binary should be in the `bin` folder.
