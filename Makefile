
build:
	@go build

## write all the methods to _gen.go
gen: build
	@./pipeline -methods=Merge,Fanout,Transform,Apply,Papply,Map,Pmap,Filter,Ptransform,SendAll,RecvAll,RecvN,Buffer -type=int,*string

## run 'gen', then move the file into
## the gen-test subdirectory, then run the tests
## there.
test: build gen
	@mv _gen.go gen-test/gen.go
	@go test -v ./gen-test