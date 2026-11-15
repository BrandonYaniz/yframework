# yFramework

yFramework is a Go-based framework designed for building robust, maintainable applications with strong foundations in error handling, validation, logging, and database integration.

## Features

### Milestone 1: Foundation (Completed)
- **Repository bootstrap**: Established the module structure
- **Error model foundation**: Structured error handling with stable codes and public-safe messages
- **Validation foundation**: Reusable validation primitives with field errors and deterministic output  
- **Bootstrap logger**: Logger that outputs to stderr for early bootstrapping

### Milestone 2: MariaDB and Migrations (In Progress)
- MariaDB connection and configuration utilities
- Transaction helper functions
- Migration registry management
- Database integration components

## Architecture

yFramework follows strict architectural boundaries:
- No framework-owned application `main` or executable templates
- No service locator or generic dependency bag  
- No independent component worker pools
- No direct external transport from components
- MariaDB only as the relational database

## Getting Started

```go
import "github.com/BrandonYaniz/yframework"
```

## Roadmap to v1

1. Complete MariaDB integration and migration system
2. Add comprehensive documentation
3. Implement full test coverage for all packages
4. Add example applications demonstrating usage
5. Finalize API design and stability

## License

MIT License - see LICENSE file for details.