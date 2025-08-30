# Pincer

A modern Milsim unit management platform built with Go and gRPC, inspired by Unit Commander and PERSCOM. Pincer provides military simulation communities with tools to manage units, personnel, and operations through a robust API-first architecture.

## Features

- **Unit Management**: Create and manage military simulation units with hierarchical structures
- **User System**: Comprehensive user profiles with preferences and role management
- **gRPC API**: High-performance, type-safe API built with Protocol Buffers
- **Modular Design**: Clean separation of concerns with dedicated modules for users, units, and core functionality

## Architecture

This repository contains the core API definitions and generated Go code for the Pincer platform:

- **Protocol Buffers**: Type-safe API definitions in `api/`
- **Generated Code**: Auto-generated Go packages in `pkg/api/gen/`
- **Core Services**: User management, unit management, and shared utilities

### API Modules

- `milsimtools.users.v1` - User accounts, profiles, and preferences
- `milsimtools.units.v1` - Military unit structures and management
- `milsimtools.core.v1` - Shared status codes and error handling

## Development

### Prerequisites

- Go 1.24+
- [buf](https://buf.build) - Protocol buffer toolchain
- [just](https://github.com/casey/just) - Command runner

### Setup

1. Clone the repository:
   ```bash
   git clone https://github.com/milsim-tools/pincer.git
   cd pincer
   ```

2. Generate Go code from protobuf definitions:
   ```bash
   just gen
   ```

### Available Commands

- `just gen` - Generate protobuf Go code (includes clean)
- `just clean` - Remove generated files
- `just proto-lint` - Lint protobuf files
- `just proto-breaking` - Check for breaking changes against main branch

### Code Generation

The project uses [buf](https://buf.build) to generate Go code from Protocol Buffer definitions. Generated files are located in `pkg/api/gen/` and should not be edited manually.

## Project Structure

```
├── api/                    # Protocol Buffer definitions
│   └── milsimtools/
│       ├── core/v1/       # Core status and error types
│       ├── units/v1/      # Unit management APIs
│       └── users/v1/      # User management APIs
├── pkg/
│   ├── api/gen/           # Generated Go code
│   └── pincer/            # Core Go packages
├── buf.yaml               # Buf configuration
├── buf.gen.yaml          # Code generation config
└── Justfile              # Development commands
```

## Contributing

1. Make changes to `.proto` files in `api/`
2. Run `just proto-lint` to validate changes
3. Run `just gen` to regenerate Go code
4. Test your changes and submit a pull request

### Protobuf Style Guide

- Use `proto3` syntax
- Package naming: `milsimtools.{module}.v1`
- Field naming: `snake_case`
- Enum values: `SCREAMING_SNAKE_CASE`
- Services: PascalCase with `Service` suffix

## License

This project is licensed under the Apache-2.0 License - see the [LICENSE](./LICENSE) file for details.
