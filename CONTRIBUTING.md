# Contributing to CrawlScraper

Thanks for your interest in contributing! 🎉

## Quick Start

1. **Fork** the repository
2. **Clone** your fork:
   ```bash
   git clone https://github.com/YOUR_USERNAME/CrawlScraper.git
   cd CrawlScraper
   ```
3. **Install dependencies**:
   ```bash
   go mod download
   ```
4. **Make your changes**
5. **Test**:
   ```bash
   go run example_single_url.go //change the url as per requirement
   ```
6. **Submit a Pull Request**

## How to Contribute

### Found a Bug?
- Check if it's already reported in [Issues](https://github.com/harshitbansal184507/CrawlScraper/issues)
- If not, create a new issue with:
  - What happened
  - What you expected
  - Code to reproduce it

### Want to Add a Feature?
- Open an issue first to discuss it
- Get feedback before starting work
- Submit a PR when ready

### Improving Documentation?
- Fix typos, add examples, clarify explanations
- Documentation PRs are always welcome!

## Coding Guidelines

### Keep it Simple
```go
// ✅ Good - clear and simple
func ValidateURL(url string) error {
    if url == "" {
        return fmt.Errorf("URL cannot be empty")
    }
    return nil
}


### Format Your Code
```bash
go fmt ./...
```

### Write Tests
```go
func TestValidateURL(t *testing.T) {
    err := ValidateURL("")
    if err == nil {
        t.Error("expected error for empty URL")
    }
}
```

## Commit Messages

Keep them clear:
```bash
git commit -m "fix(HTTP CLIENT):timeout handling in HTTP client"
git commit -m "feat(Parser):Add support for custom headers"
git commit -m "docs(Readme) :Update README with installation steps"


## Pull Request Process

1. Create a new branch:
   ```bash
   git checkout -b fix/your-fix
   ```
2. Make your changes
3. Push to your fork:
   ```bash
   git push origin fix/your-fix
   ```
4. Create a Pull Request with:
   - Clear title
   - What changed
   - Why it's needed

## Code of Conduct

- Be respectful
- Be helpful

## Questions?

- Check existing [Issues](https://github.com/harshitbansal184507/CrawlScraper/issues)
- Create a new issue with the `question` label
- We're happy to help!

---

**Thank you for contributing!** 🕷️