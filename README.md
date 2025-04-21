# victorgo
# Victor SDK for Golang

Consider this to run a mini test:

## Quick Start

### Project Structure
```
victorgo/
├── examples
│   └── sdk_example
│       └── main.go // run this file
├── include/        // Header files
│   └── victor/
│       └── victor.h
└── lib/            // Library files
    └── libvictor.so (or .dylib, .dll, etc.) 
```

1. Ensure the Victor library (`libvictor.so`, `.dylib`, or `.dll`) is placed in the `lib/` directory
2. Run the example code:

```bash
go run ./examples/sdk_example/main.go
```
