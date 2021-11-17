# distributed-services-in-go-proglog

## Generating the pb.go from proto file 

- Generating from log.proto at root level

```shell
make protogen
```

## Notes: 

- Logging related terminology used in this exercise
  - Record—the data stored in our log.
  - Store—the file we store records in.
  - Index—the file we store index entries in.
  - Segment—the abstraction that ties a store and an index together.
  - Log—the abstraction that ties all the segments together.