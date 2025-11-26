# 🤝 Contributing to Dejavu

Terima kasih atas minat Anda untuk berkontribusi! Panduan ini akan membantu Anda memulai.

## Code of Conduct

Dengan berpartisipasi dalam project ini, Anda setuju untuk mematuhi code of conduct kami:
- Bersikap hormat dan profesional
- Menerima kritik yang konstruktif
- Fokus pada yang terbaik untuk komunitas
- Menunjukkan empati terhadap anggota komunitas lainnya

## How to Contribute

### Reporting Bugs

Sebelum membuat bug report:
1. Check existing issues untuk memastikan bug belum dilaporkan
2. Gunakan template bug report yang tersedia
3. Sertakan informasi lengkap:
   - Versi OS dan software
   - Steps to reproduce
   - Expected vs actual behavior
   - Screenshots jika applicable

### Suggesting Features

Feature requests sangat diterima! Pastikan:
1. Jelaskan use case yang jelas
2. Deskripsikan solusi yang Anda usulkan
3. Diskusikan alternatif yang sudah Anda pertimbangkan

### Pull Requests

#### Setup Development Environment

```bash
# Fork dan clone repository
git clone https://github.com/YOUR_USERNAME/dejavu.git
cd dejavu

# Create feature branch
git checkout -b feature/amazing-feature

# Setup dependencies
make infra
cd backend && go mod download
cd ../builder && go mod download
cd ../deployer && go mod download
cd ../frontend && npm install
```

#### Development Workflow

1. **Make Changes**
   - Write clean, readable code
   - Follow existing code style
   - Add tests for new features
   - Update documentation

2. **Test Your Changes**
   ```bash
   # Backend tests
   cd backend && go test ./...
   
   # Frontend tests
   cd frontend && npm test
   ```

3. **Commit Your Changes**
   ```bash
   git add .
   git commit -m "feat: add amazing feature"
   ```

   **Commit Message Format:**
   - `feat:` New feature
   - `fix:` Bug fix
   - `docs:` Documentation changes
   - `style:` Code style changes (formatting, etc)
   - `refactor:` Code refactoring
   - `test:` Adding/updating tests
   - `chore:` Maintenance tasks

4. **Push and Create PR**
   ```bash
   git push origin feature/amazing-feature
   ```
   
   Then create Pull Request di GitHub dengan:
   - Clear title dan description
   - Reference related issues
   - Screenshots untuk UI changes

#### Code Style Guidelines

**Go:**
```go
// Use gofmt for formatting
gofmt -w .

// Follow standard Go conventions
// - Use camelCase for variables
// - Use PascalCase for exported functions
// - Add comments for exported functions
```

**TypeScript/React:**
```typescript
// Use Prettier for formatting
npm run format

// Follow React best practices
// - Use functional components
// - Use TypeScript types
// - Avoid inline styles
```

**Git:**
- Keep commits atomic and focused
- Write descriptive commit messages
- Rebase before merging to keep history clean

### Project Structure

```
dejavu/
├── backend/          # Go API server
│   ├── cmd/          # Entry points
│   ├── internal/     # Business logic
│   └── pkg/          # Shared packages
├── builder/          # Build worker
├── deployer/         # Deploy worker
├── frontend/         # Next.js app
└── infra/            # Infrastructure configs
```

### Testing

#### Backend Tests

```go
// backend/internal/service/project_service_test.go
func TestCreateProject(t *testing.T) {
    // Arrange
    service := NewProjectService(mockRepo)
    
    // Act
    project, err := service.Create("user123", &req)
    
    // Assert
    assert.NoError(t, err)
    assert.NotNil(t, project)
}
```

#### Frontend Tests

```typescript
// frontend/components/Button.test.tsx
import { render, screen } from '@testing-library/react'
import { Button } from './Button'

test('renders button with text', () => {
  render(<Button>Click me</Button>)
  expect(screen.getByText('Click me')).toBeInTheDocument()
})
```

### Documentation

Update documentation when:
- Adding new features
- Changing APIs
- Updating configuration
- Improving developer experience

Files to update:
- `README.md` - Main documentation
- `API.md` - API changes
- `DEPLOYMENT.md` - Infrastructure changes
- Inline code comments

### Review Process

1. **Automated Checks**
   - CI/CD pipeline runs tests
   - Linting checks
   - Build verification

2. **Code Review**
   - At least 1 approval required
   - Address review comments
   - Keep discussions constructive

3. **Merge**
   - Squash and merge for clean history
   - Delete feature branch after merge

## Development Tips

### Local Development

```bash
# Watch mode for backend
cd backend
air # or go run cmd/api/main.go

# Watch mode for frontend
cd frontend
npm run dev

# View logs
make logs
```

### Debugging

**Backend:**
```go
// Add debug logs
log.Printf("Debug: %+v", data)

// Use debugger
dlv debug cmd/api/main.go
```

**Frontend:**
```typescript
// Browser DevTools
console.log('Debug:', data)

// React DevTools extension
```

### Common Issues

**Issue:** `go mod` errors
```bash
# Solution
go mod tidy
go mod download
```

**Issue:** Port already in use
```bash
# Solution
lsof -ti:8080 | xargs kill -9
```

**Issue:** Docker build fails
```bash
# Solution
docker system prune -a
```

## Community

- **Discord:** [Join our server](#)
- **Twitter:** [@dejavu_dev](#)
- **Email:** dev@dejavu.id

## Recognition

Contributors will be added to:
- Contributors list in README
- Release notes for significant contributions
- Hall of Fame for major features

## License

By contributing, you agree that your contributions will be licensed under the MIT License.

---

**Thank you for contributing to Dejavu! 🚀**

