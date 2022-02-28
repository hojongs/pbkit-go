# pbkit-go: Protobuf Toolkit Written in Go

## Challenges

### Language differences

- Inheritance of interface: There is no inheitance in Go. But it can be done by anonymous field as a composition.
- Union type: Replace with super type
- Generic: Generic is supported since 1.18
  - https://go.dev/doc/tutorial/generics
- No need type field for each interfaces
- There are no Iterable.map, Collection.shift methods.
- There is no built-in string template in Go
- Entrypoint in Go is main function in package main
  - it should be independent module (이건 pbkit이 특이한 것)

### Others

- Circular dependency between files, directories.
  - files: core/ast/index.ts vs core/ast/lexical-elements.ts
  - directories: core/ast vs core/parser
- Some code seems to need refactoring (spent a lot of time to understand code)
  - circular dependency
  - core dir is shared by everywhere: ambiguous module boundary
- Too coupled code
  - analyzeDeps()에서 호출하는 getPollapoYml()는 cacheDeps()와 couple 되어 있음 (cache를 지우면 pollapo yml file not found error가 발생함)
- Lack of documentation: --clean option -> description: Don't use cache -> behavior: remove cache
  - 각 package에 대한 description, cache dir location 등 documentation
