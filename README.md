<p align='center'>
  <img src="https://raw.githubusercontent.com/milsim-tools/pincer/refs/heads/main/assets/logo-color.svg" alt="Milsim.tools logo" style="width: 250px; height: 250px"></img>

  <h1 align='center'>Pincer</h1>

  <p align='center'>
    A modern military simulation unit management platform built with Go and 
    gRPC. Inspired by Unit Commander and PERSCOM, Pincer provides milsim 
    communities with tools to manage units, personnel, and operations through 
    a robust API-first architecture.
  </p>
</p>

## Features

- **Unit Management**: Create and manage military simulation units with 
  hierarchical structures
- **User System**: Comprehensive user profiles with preferences and role 
  management
- **gRPC API**: High-performance, type-safe API built with Protocol Buffers
- **Modular Design**: Clean separation of concerns with dedicated modules for 
  users, units, and members

## Architecture

Pincer is a complete Go application providing gRPC services for military 
simulation unit management:

- **Protocol Buffers**: Type-safe API definitions in `api/`
- **Generated Code**: Auto-generated Go packages in `pkg/api/gen/`
- **Business Logic**: Service implementations in domain packages
- **Database**: PostgreSQL with GORM for persistence
- **CLI Interface**: Command-line application for running services

### API Modules

- `milsimtools.users.v1` - User accounts, profiles, and preferences
- `milsimtools.units.v1` - Military unit structures and management
- `milsimtools.members.v1` - Unit membership and personnel management

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

2. Install dependencies and generate code:
   ```bash
   just gen
   ```

3. Set up environment configuration:
   ```bash
   just dotenv
   ```

### Available Commands

- `just build` - Compile the pincer binary
- `just run` - Build and run the application locally
- `just gen` - Generate protobuf Go code (includes clean)
- `just proto-lint` - Lint protobuf files
- `just proto-breaking` - Check for breaking changes against main branch
- `just fmt` - Format Go code
- `just lint` - Lint Go code with go vet
- `just check` - Compile everything to verify builds
- `just test` - Run Go tests

### Code Generation

The project uses [buf](https://buf.build) to generate Go code from Protocol 
Buffer definitions. Generated files are located in `pkg/api/gen/` and should 
not be edited manually.

## Project Structure

```
├── api/                    # Protocol Buffer definitions
│   └── milsimtools/
│       ├── members/v1/     # Member management APIs
│       ├── units/v1/       # Unit management APIs
│       └── users/v1/       # User management APIs
├── cmd/pincer/             # CLI application entry point
├── pkg/
│   ├── api/gen/            # Generated Go code
│   ├── members/            # Member service implementation
│   ├── units/              # Unit service implementation
│   ├── users/              # User service implementation
│   └── pincer/             # Core application logic
├── internal/               # Internal packages
├── buf.yaml                # Buf configuration
├── buf.gen.yaml           # Code generation config
└── Justfile               # Development commands
```

## Contributing

1. Make changes to `.proto` files in `api/` or Go code in `pkg/`
2. Run `just proto-lint` to validate protobuf changes
3. Run `just gen` to regenerate Go code from protobuf definitions
4. Run `just fmt && just lint && just test` to verify code quality
5. Test your changes and submit a pull request

### Protobuf Style Guide

- Use `proto3` syntax
- Package naming: `milsimtools.{module}.v1`
- Field naming: `snake_case`
- Enum values: `SCREAMING_SNAKE_CASE`
- Services: PascalCase with `Service` suffix

## License

This project is licensed under the Apache-2.0 License - see the 
[LICENSE](./LICENSE) file for details.
