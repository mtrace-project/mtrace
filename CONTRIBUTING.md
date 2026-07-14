# Contributing to Mtracer

First off, thank you for taking the time to contribute! Contributions from the community help make **Mtracer** better for everyone.

Here is a guide to help you get started with contributing to the project.

---

## Code of Conduct

By participating in this project, you agree to maintain a respectful and welcoming environment for all contributors and maintainers.

## Getting Started

### Prerequisites

To build and test Mtracer locally, you will need:
* **Go**: version `1.26` or higher
* **Golangci-lint**: for static analysis and linting
* **Gofumpt**: for code formatting
* **Pre-commit**: for automating pre-commit validation

### Local Setup

1. **Clone the repository**:
   ```bash
   git clone https://github.com/mtracer-project/mtracer.git
   cd mtracer
   ```

2. **Download dependencies**:
   ```bash
   go mod download
   ```

3. **Install Pre-commit hooks**:
   We use pre-commit hooks to automate styling, linting, and commit message checks.
   ```bash
   # Install pre-commit via pip or your system package manager (e.g., brew install pre-commit)
   pip install pre-commit
   
   # Register the hooks
   pre-commit install --hook-type pre-commit --hook-type commit-msg
   
   # (Optional) Update the hooks to the latest versions
   pre-commit autoupdate
   ```

---

## Coding Guidelines

To keep the codebase clean and uniform, please adhere to the following guidelines before submitting your changes.

### Code Formatting
We use `gofumpt` (a stricter version of `gofmt`) to format the codebase. To format all files, run:
```bash
gofumpt -l -e -w .
```
This is automatically run as a pre-commit hook.

### Linting
We use `golangci-lint` to check code quality and enforce guidelines. To run the linter locally:
```bash
golangci-lint run
```
Linting rules are configured in the `.golangci.yml` file and will block build checks in CI if they fail.

### Running Tests
Make sure all existing and new tests pass before submitting a pull request.
```bash
go test -v -race ./...
```

---

## Commit Guidelines

We enforce the [Conventional Commits](https://www.conventionalcommits.org) specification for all commit messages. This helps automate our release generation and changelog updates.

### Format
Commit messages should follow the structure:
```
<type>(<scope>): <description>

[optional body]

[optional footer(s)]
```

* **Types** must be one of:
  * `feat`: A new feature (included in the release changelog).
  * `fix`: A bug fix (included in the release changelog).
  * `chore`: Maintenance tasks, dependencies updates, etc.
  * `ci`: Continuous integration configurations or scripts.
  * `docs`: Documentation changes.
  * `style`: Code style changes (whitespace, formatting, missing semi-colons).
  * `refactor`: A code change that neither fixes a bug nor adds a feature.
  * `test`: Adding missing tests or correcting existing tests.
* **Scope**: Optional, indicating the package or module affected (e.g., `parser`, `cmd`, `trace`).
* **Breaking Changes**: Must be indicated by appending a `!` to the type/scope (e.g., `feat(cmd)!: change input flag syntax`) or putting `BREAKING CHANGE:` in the footer.

---

## Release Process

* **Develop Branch**: Code pushed to `develop` triggers a **snapshot** release build using GoReleaser to verify compile and packaging workflows.
* **Main Branch & Tags**: Official production releases are created by tagging a commit with a version in `vX.Y.Z` format. While you can trigger a release from any branch, it is highly recommended to release from the `main` branch or merge into `main` immediately afterward.

