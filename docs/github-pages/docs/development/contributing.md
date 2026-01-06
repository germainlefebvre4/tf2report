---
sidebar_position: 6
---

# Contributing

Guidelines for contributing to tf2report.

## Welcome!

We appreciate your interest in contributing to tf2report! This document provides guidelines and best practices for contributing.

## Code of Conduct

- Be respectful and inclusive
- Focus on constructive feedback
- Help others learn and grow

## Getting Started

1. **Fork the repository** on GitHub
2. **Clone your fork** locally
3. **Create a feature branch** from `main`
4. **Make your changes**
5. **Test thoroughly**
6. **Submit a pull request**

## Development Setup

```bash
git clone https://github.com/your-username/tf2report.git
cd tf2report
make deps
make build
make test
```

See [Development Getting Started](./getting-started.md) for details.

## Making Changes

### Branch Naming

Use descriptive branch names:

- `feat/add-xml-format` - New features
- `fix/parser-error` - Bug fixes
- `docs/update-readme` - Documentation
- `test/add-coverage` - Tests
- `refactor/simplify-config` - Refactoring

### Commit Messages

Follow conventional commits:

```
<type>: <description>

[optional body]

[optional footer]
```

**Types:**
- `feat` - New feature
- `fix` - Bug fix
- `docs` - Documentation
- `test` - Tests
- `refactor` - Code refactoring
- `chore` - Maintenance

**Examples:**
```
feat: add XML output format

fix: handle empty plan files correctly

docs: update configuration examples
```

### Coding Standards

Follow guidelines in `.github/instructions/`:

1. **Go Code** - Idiomatic Go practices
2. **CLI Code** - Clear help text, flag naming
3. **Makefile** - Proper phony targets, tabs

**Key rules:**
- Run `make fmt` before committing
- Run `make test` to verify changes
- Document all exported functions
- Handle errors properly
- Write tests for new features

### Testing

1. **Write tests** for new features
2. **Update tests** for changes
3. **Run all tests** before submitting

```bash
make test
make test-coverage
```

### Documentation

Update documentation when:
- Adding new features
- Changing CLI options
- Modifying configuration
- Updating APIs

Files to update:
- `README.md` - Main documentation
- `docs/*.md` - Detailed guides
- Code comments - Exported items
- Examples - Usage examples

## Pull Request Process

### Before Submitting

- [ ] Code follows style guidelines
- [ ] All tests pass
- [ ] New tests added for features
- [ ] Documentation updated
- [ ] Commits are clean and descriptive
- [ ] No merge conflicts

### Submitting PR

1. **Push your branch** to your fork
2. **Create pull request** against `main`
3. **Fill out PR template** completely
4. **Link related issues** if any

### PR Description

Include:
- **What** - What does this change
- **Why** - Why is it needed
- **How** - How it works
- **Testing** - How you tested it
- **Screenshots** - If UI changes

**Example:**
```markdown
## Description
Adds XML output format support.

## Motivation
Users requested XML format for integration with legacy tools.

## Changes
- Added `pkg/report/xml.go` with XMLFormatter
- Updated `NewFormatter()` to support XML
- Added tests in `pkg/report/xml_test.go`
- Updated documentation

## Testing
- Unit tests pass
- Tested with sample plans
- Verified XML output validates

## Related Issues
Closes #123
```

### Review Process

1. **Automated checks** must pass (CI/CD)
2. **Code review** by maintainers
3. **Address feedback** promptly
4. **Update as needed**
5. **Merge** when approved

## Types of Contributions

### Bug Reports

Open an issue with:
- Clear title
- Steps to reproduce
- Expected vs actual behavior
- Environment details
- Sample input if applicable

### Feature Requests

Open an issue with:
- Clear description
- Use case
- Proposed solution
- Alternatives considered

### Code Contributions

- Bug fixes
- New features
- Performance improvements
- Test improvements
- Documentation updates

### Documentation

- Fix typos
- Improve clarity
- Add examples
- Update outdated info

## Development Guidelines

### Adding Features

1. **Discuss first** - Open an issue
2. **Design** - Plan the approach
3. **Implement** - Write code
4. **Test** - Add comprehensive tests
5. **Document** - Update docs
6. **Submit** - Create PR

### Code Review

When reviewing:
- Be constructive
- Ask questions
- Suggest improvements
- Appreciate contributions

When receiving review:
- Be open to feedback
- Ask for clarification
- Make requested changes
- Thank reviewers

## Release Process

Maintainers handle releases:

1. Update version
2. Update CHANGELOG
3. Create release tag
4. Build binaries
5. Publish release

## Questions?

- Check existing documentation
- Search existing issues
- Open a new issue
- Join discussions

## License

By contributing, you agree that your contributions will be licensed under the project's license.

## Thank You!

Every contribution helps make tf2report better. Thank you for your time and effort!

## Next Steps

- **[Development Getting Started](./getting-started.md)** - Setup development environment
- **[API Reference](./api-reference.md)** - Understand the codebase
- **[Testing](./testing.md)** - Write effective tests
