# MyGit

A Git-like version control system implementation in Go, built for educational purposes. This project demonstrates the core concepts and internal architecture of Git by implementing its fundamental data structures and operations from scratch.

## Overview

MyGit is a simplified version control system that replicates Git's core functionality using content-addressed storage, SHA-1 hashing, and zlib compression. The implementation follows Git's three-tree architecture and provides a command-line interface for basic repository operations.

## Features

- **Repository Initialization**: Create a new Git-compatible repository structure
- **Object Database**: Store and retrieve blobs, trees, and commits using SHA-1 content addressing
- **Staging Area**: Index-based staging mechanism for preparing commits
- **Commit Creation**: Snapshot working directory state with commit objects
- **Content Inspection**: Read and display stored objects by their hash

## Architecture

### Three-Layer Design

The project follows a clean architecture pattern with three distinct layers:

1. **CLI Layer** (`cmd/mygit/`): Command-line interface and argument parsing
2. **Command Layer** (`internal/commands/`): User-facing command implementations
3. **Core Layer** (`pkg/`): Core functionality including object model, storage, and repository management

### Object Model

MyGit implements Git's object model with three fundamental types:

- **Blob**: Stores raw file content
- **Tree**: Represents directory structure with file modes, names, and object references
- **Commit**: Captures snapshot metadata including tree reference, parent commits, author information, and commit message

All objects follow the format: `<type> <size>\0<content>`

### Storage System

Objects are stored using content-addressed storage:

1. Serialization to Git object format
2. SHA-1 hash computation (40 hexadecimal characters)
3. Zlib compression
4. Storage in `.git/objects/<first-2-chars>/<remaining-38-chars>`

## Installation

```bash
# Clone the repository
git clone https://github.com/SteliosSpanos/mygit.git
cd mygit

# Build the executable
go build -o mygit cmd/mygit/main.go
```

## Usage

### Initialize a Repository

```bash
./mygit init
```

Creates a `.git` directory with the following structure:
```
.git/
├── HEAD              # Points to current branch
├── objects/          # Content-addressed object storage
└── refs/
    ├── heads/        # Branch references
    └── tags/         # Tag references
```

### Hash and Store Files

```bash
./mygit hash-object <file>
```

Reads a file, creates a blob object, stores it in the object database, and outputs the SHA-1 hash.

### Retrieve Objects

```bash
./mygit cat-file <hash>
```

Retrieves and displays an object from the database by its SHA-1 hash.

### Stage Files

```bash
./mygit add <file>
```

Adds a file to the staging area (index) by creating a blob and recording its hash, mode, and path.

### Create Commits

```bash
./mygit commit -m "commit message"
```

Creates a commit object from staged files, builds a tree structure, and updates the current branch reference.

## Project Structure

```
mygit/
├── cmd/mygit/              # Main entry point
│   └── main.go
├── internal/commands/      # Command implementations
│   ├── init.go
│   ├── hash_object.go
│   ├── cat_file.go
│   ├── add.go
│   └── commit.go
├── pkg/
│   ├── objects/           # Object model and serialization
│   │   ├── object.go
│   │   ├── blob.go
│   │   ├── tree.go
│   │   └── commit.go
│   ├── storage/           # Object storage and retrieval
│   │   └── storage.go
│   ├── repository/        # Repository initialization
│   │   └── repository.go
│   ├── index/             # Staging area management
│   │   └── index.go
│   ├── refs/              # Branch reference handling
│   │   └── refs.go
│   └── tree/              # Tree building utilities
│       └── builder.go
├── go.mod
├── CLAUDE.md              # Development guidelines
└── README.md
```

## Technical Details

### Content-Addressed Storage

MyGit uses SHA-1 hashing to create unique identifiers for all objects. Identical content produces identical hashes, enabling automatic deduplication and efficient storage.

### Object Serialization

Each object type implements the `Object` interface:

```go
type Object interface {
    Type() ObjectType
    Serialize() ([]byte, error)
    Deserialize(data []byte) error
}
```

### Index Format

The staging area uses a simplified text-based format:
```
<mode> <hash> <path>
```

Each line represents a staged file with its Unix permission mode, SHA-1 hash, and relative path from the repository root.

### Branch References

Branches are implemented as files in `.git/refs/heads/` containing the SHA-1 hash of the latest commit. The `HEAD` file contains a symbolic reference to the current branch.

## Implementation Notes

- Default branch name is `main` (configurable in `pkg/repository/repository.go`)
- Tree entries are stored in sorted order by name (Git requirement)
- Tree hashes are stored as 20-byte binary values, not 40-character hex strings
- File modes follow Unix conventions: `100644` (regular), `100755` (executable), `040000` (directory)

## Learning Objectives

This project demonstrates:

- Content-addressed storage systems
- Object serialization and deserialization
- Data compression techniques
- Hash-based data structures
- File system operations in Go
- Interface-based polymorphism
- Repository structure and version control concepts

## Limitations

This is an educational implementation with the following limitations:

- No support for subdirectories in commits (flat structure only)
- Simplified index format (no metadata like timestamps or file size)
- No branch merging or conflict resolution
- No remote repository operations (clone, push, pull)
- No packfiles or delta compression
- No garbage collection for unreferenced objects

## Future Enhancements

Potential additions to extend functionality:

- `status` command to show working directory state
- `log` command for commit history visualization
- `branch` and `checkout` commands for branch management
- `diff` command for comparing file versions
- Nested directory support in tree objects
- Merge functionality with conflict detection

## Requirements

- Go 1.25.5 or higher

## License

This project is intended for educational purposes.

## Contributing

This is a learning project. Feel free to fork and experiment with your own implementations.

## References

- [Pro Git Book - Git Internals](https://git-scm.com/book/en/v2/Git-Internals-Plumbing-and-Porcelain)
- [Git Object Model](https://git-scm.com/book/en/v2/Git-Internals-Git-Objects)
- [Building Git by James Coglan](https://shop.jcoglan.com/building-git/)
