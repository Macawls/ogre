# Contributing to Ogre

## Reporting Issues

Open a GitHub issue with:
- What you expected to happen
- What actually happened
- HTML/CSS input that reproduces the issue
- Output format (SVG/PNG)

## Development

### Prerequisites

- Go 1.25+
- Bun (for running Satori comparison tests)

### Running tests

```bash
go test ./...
```

### Running Satori comparison tests

```bash
cd test/satori-reference && bun install && bun run generate.ts
cd test && go test -v -run TestCompareWithSatori
```

### Code style

- Idiomatic Go
- No unnecessary comments (code should be self-explanatory)
- Doc comments on all exported types and functions
- No external dependencies beyond golang.org/x/*

### Pull requests

1. Fork the repo
2. Create a branch from main
3. Make your changes
4. Run `go test ./...` and `go vet ./...`
5. Submit a PR with a description of the change

### Adding test fixtures

Test fixtures go in `test/fixtures/` as HTML files. After adding a fixture:

1. Run `cd test/satori-reference && bun run generate.ts` to create the Satori reference
2. Run `cd test && go test -run TestCompareWithSatori` to verify
