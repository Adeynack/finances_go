# Import from Moneydance

## Run

```bash
go run cmd/import-from-moneydance/*.go
```

## Parameters

| Key                | Default         | Description                                                        |
| ------------------ | --------------- | ------------------------------------------------------------------ |
| `md-import-file`   | ./tmp/md.json   | Path to the JSON file to import, exported from Moneydance.         |
| `book-owner-email` | joe@example.com | E-Mail of the owner of the book to create.                         |
| `default-currency` | EUR             | Default currency of the book to create.                            |
| `auto-delete-book` | 0 (false)       | Destroy the book if it already exists (clean import from scratch). |

Here's a useful command during development. It is faster than destroying the existing book. It
simply truncates all of the development database, seeds it with the test fixtures and imports.

```bash
make db_seed build_import_from_moneydance && out/import-from-moneydance
```
