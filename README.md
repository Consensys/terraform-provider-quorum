# Terraform Provider for Quorum

[Quorum](https://goquorum.com) is an Ethereum-based distributed ledger protocol with transaction/contract privacy and new consensus mechanisms.

Quorum is a fork of [go-ethereum](https://github.com/ethereum/go-ethereum) and is updated in line with go-ethereum releases.

This plugin provides number of data sources and resources which can be used to bootstrap a Quorum Network from scratch.

* Website: https://docs.goquorum.com 
* Slack Channel: https://bit.ly/quorum-slack

## Requirements

* Terraform 0.12.x
* Go 1.13.x (to build the provider plugin)

## Using the provider

See the [Quorum Provider documentation](website/docs) to get started using the Quorum provider.

Also check out some examples in [examples directory](examples).

## Developing the provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine.

*Note:* This project uses [Go Modules](https://blog.golang.org/using-go-modules) making it safe to work with it outside of your existing [`GOPATH`](http://golang.org/doc/code.html#GOPATH). The instructions that follow assume a directory in your home directory outside of the standard `GOPATH` (i.e `$HOME/development/terraform-providers/`).

Clone repository to: `$HOME/development/terraform-providers/`

```sh
$ mkdir -p $HOME/development/terraform-providers/; cd $HOME/development/terraform-providers/
$ git clone git@github.com:jpmorganchase/terraform-provider-quorum.git
...
```

Enter the provider directory and run `make` to compile the provider. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

```sh
$ make
...
```

## Testing the provider

In order to test the provider, you can run `go test ./quorum -v`

## Generate new documentation website

Go to the `website` folder and do:
```
go build gen.go
./gen
rm -f gen
```

This will update automatically the website documentation.
