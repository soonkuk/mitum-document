# contest

contest is the simulation tool for mitum and it's children.

[![CircleCI](https://img.shields.io/circleci/project/github/spikeekips/contest/main.svg?style=flat-square&logo=circleci&label=circleci&cacheSeconds=60)](https://circleci.com/gh/spikeekips/contest/tree/main)
[![GoDoc](https://godoc.org/github.com/golang/gddo?status.svg)](https://pkg.go.dev/github.com/spikeekips/contest?tab=overview)
[![Go Report Card](https://goreportcard.com/badge/github.com/spikeekips/contest)](https://goreportcard.com/report/github.com/spikeekips/contest)
[![codecov](https://codecov.io/gh/spikeekips/contest/branch/master/graph/badge.svg)](https://codecov.io/gh/spikeekips/contest)
[![](http://tokei.rs/b1/github/spikeekips/contest?category=lines)](https://github.com/spikeekips/contest)

# Install

```sh
$ git clone https://github.com/spikeekips/contest
$ cd contest
$ go build -o ./contest
```

# Run

* Before running contest, you need to build mitum or mitum variants(ex. [mitum-currency](https://github.com/spikeekips/mitum-currency).
* Before running contest, check contest help, `$ contest --help`
* By default, contest looks for local mongodb(`mongodb://localhost:27017`)

```sh
$ ./contest run --log-level debug --exit-after 2m ./mitum-document ./scenario/standalone-run-create-blocksign-document.yml
```

* You can find some example scenarios for contest at `./scenario` in this repository.
