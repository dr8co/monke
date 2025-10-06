# Contributing to Monke

Thanks for your interest in contributing! This document explains the recommended workflow,
coding conventions, and how to submit changes so they can be reviewed and merged quickly.

## Quick summary

- Open an issue to discuss larger changes before starting work.
- Fork the repository and create a branch for your change.
- Follow the project's code style and run tests locally.
- Submit a pull request (PR) with a clear description and related issue (if any).

## Reporting issues

- Use the issue tracker to report bugs or request features.
- Include a minimal reproduction case, the behavior you expect, and any relevant environment details (Go version, OS).

## Development workflow

1. Fork the repository.
2. Create a feature branch:
   - git checkout -b feat/short-description
3. Make small, focused commits with clear messages.
4. Rebase or merge main as needed to keep your branch up-to-date.
5. Push and open a PR against the main branch.

## Testing

- Run tests locally before opening a PR:
  - go test ./...
- If you add new behavior, include unit tests that cover it.
- Contributors should update existing tests or add new tests as needed when introducing new features, fixing bugs that change behavior, or modifying public APIs.

## Code style and formatting

- This is a Go project. Please:
  - run go fmt / goimports on changed files
  - keep code simple and idiomatic
- Prefer small, well-scoped changes and maintain readability.

## Commit messages

- Use concise, descriptive messages. Example:
  - feat(parser): add support for ternary operator
  - fix(evaluator): correct integer division semantics

## Pull requests

- Provide a clear description of what the PR changes and why.
- Link to any related issue(s).
- Add tests for new behavior where appropriate.
- Expect code review and iterate on requested changes.
- The maintainers will merge when CI passes and changes are approved.

## Documentation

- Update docs/ (or README.md) for user-visible changes.
- Keep examples and usage up-to-date.

## CI and checks

- The project uses standard Go tooling and GitHub Actions for CI.
- Ensure `go test ./...` passes locally; CI will run tests on PRs.

## Licensing & CLA

- By contributing, you agree that your contributions are licensed under this repository's license (see LICENSE).

## Code of Conduct

- Be respectful and constructive in discussions and reviews.

Thank you for helping improve Monke!
