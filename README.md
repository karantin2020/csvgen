# CSVGEN 
Generates csv Unmarshaller and Marshaller funcs for go struct types

Generated functions MarshallCSV and UnmarshallCSV process builtin types,
call MarshallCSV and UnmarshallCSV for custom types
and process pointer types assuming that pointers were initiated (memory is
allocated).

So you have to verify data structures for null pointers before marshalling
and unmarshalling to prevent SEGFAULT

### Usage
	
	go build -o csvgen csvgen.go
	go generate // For example usage