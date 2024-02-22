# Finances

## How to develop?

### Start development process

```bash
make dev
```

This will start watching for modifications in `go` and `templ` (templates) files, and compile and start the server on modifications.

### Debug

The included _VSCode_ settings allow to attach to existing _Go_ process (the one started by `make dev`). Select the _Attach to Process_ configuration and upon start of the debugger, find the `bin/serve` process in the list.

Caveats:

- Needs to be redone every time the server restarts (on every modifications)
  - This is far from perfect, but so far what I could come up with. This dev setup is a work in progress.
