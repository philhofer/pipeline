
build:
	@go build

## generate code in subdirectory,
## then run tests
test: build
	@./pipeline -methods=Merge,Fanout,Transform,Apply,Papply,Map,Pmap,Filter,Ptransform,SendAll,RecvAll,RecvN,Buffer -type=int,*string -o=gen-test/gen.go
	@go test -v ./gen-test

## removes generated code
clean:
	@rm gen-test/gen.go