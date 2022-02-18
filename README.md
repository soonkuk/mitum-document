### mitum-document

*mitum-document* is the data management case of mitum model, based on
[*mitum*](https://github.com/spikeekips/mitum) and [*mitum-currency*](https://github.com/spikeekips/mitum-currency).

#### Features,

* account: account address and keypair is not same.
* documentData: actual data stored in document.
* simple transaction: create document, update document, sign document.
* *mongodb*: as mitum does, *mongodb* is the primary storage.

#### Installation

> NOTE: at this time, *mitum* and *mitum-document* is actively developed, so
before building mitum-blocksign, you will be better with building the latest
mitum source.
> `$ git clone https://github.com/protoconNet/mitum-document`
>
> and then, add `replace github.com/spikeekips/mitum => <your mitum source directory>` to `go.mod` of *mitum-document*.

Build it from source
```sh
$ cd mitum-document
$ go build -ldflags="-X 'main.Version=v0.0.1'" -o ./mitum-document ./main.go
```

#### Run

At the first time, you can simply start node with example configuration.

> To start, you need to run *mongodb* on localhost(port, 27017).

```
$ ./mitum-document node init ./standalone.yml
$ ./mitum-document node run ./standalone.yml
```

> Please check `$ ./mbs --help` for detailed usage.

#### Test

```sh
$ go clean -testcache; time go test -race -tags 'test' -v -timeout 20m ./... -run .
```
