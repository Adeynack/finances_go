# Finances

## How to develop?

### Start development process

```bash
make dev
```

This will start watching for modifications in `go` and `templ` (templates) files, and compile and start the server on modifications.

### Debug

The included _VSCode_ settings allow to attach to existing _Go_ process (the one started by `make dev`). Select the _Attach to Process_ configuration and upon start of the debugger, find the `out/serve` process in the list.

Caveats:

- Needs to be redone every time the server restarts (on every modifications)
  - This is far from perfect, but so far what I could come up with. This dev setup is a work in progress.

## Database

### Concept

The `gorm` ORM is used, with DAL generation for type safe querying.

- PROS
  - Automatic logging of all queries.
  - Type-Safe GO-ish access to the data.
- CONS
  - Not as performant as raw SQL queries.
  - Vendor-lock of the coding style.

Migrations are done using `golang-migrate`. However, during development, for facilitating prototyping and be fast to change things, in `APP_ENV=dev` the _Auto Migrate_ feature of `gorm` is used. In any other environment, a [custom dialector](https://gist.github.com/molind/a67100448b886b7257e30799e06a0718) for _Auto Migrate_ will be used, ensuring that the lastest migrations correspond to the latest database `gorm` would generate. This allows a tighter control of the migration in production, allowing data transformation and clean up, and not only DDL operations.

### Sequence vs Timestamp

Since `migrate` does not support inserting migrations before the last ran -- eg: pull request `A` creates a migration on June 6th, pull request `B` creates one on June 12th, but gets merged before `A`; the migration in `A` will never be ran because the last migration ran was created after -- I decided to go for sequencial migration numbers. This way, it will be (more) obvious before merging pull requests.

`TODO`: A CI check can even be introduce to make sure that a pull request does not duplicate a migration number.
