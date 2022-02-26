# pbkit-go: Protobuf Toolkit Written in Go

## Challenges

### Language differences

- Inheritance of interface: There is no inheitance in Go. But it can be done by anonymous field as a composition.
- Union type: Replace with super type
- Generic: Generic is supported since 1.18
  - https://go.dev/doc/tutorial/generics
- No need type field for each interfaces

### Others

- Circular dependency between files, directories.
  - files: core/ast/index.ts vs core/ast/lexical-elements.ts
  - directories: core/ast vs core/parser
