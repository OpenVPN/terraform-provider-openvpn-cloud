# Terraform Provider OpenVPN Cloud

<a href="https://www.terraform.io/" target="_blank">
  <img align="right" src="https://upload.wikimedia.org/wikipedia/commons/thumb/0/04/Terraform_Logo.svg/2560px-Terraform_Logo.svg.png" alt="Terraform" width="150px">
</a>

<a href="https://anna.money/?utm_source=terraform&utm_medium=docs" target="_blank">
  <img align="right" src="https://upload.wikimedia.org/wikipedia/commons/a/aa/ANNA_Money_Logo_PNG.png" alt="ANNA Money" width="80px">
</a>

<a href="https://openvpn.net/cloud-vpn/?utm_source=terraform&utm_medium=docs" target="_blank">
  <img align="right" src="https://upload.wikimedia.org/wikipedia/commons/thumb/f/f5/OpenVPN_logo.svg/2560px-OpenVPN_logo.svg.png" alt="OpenVPN" width="150px">
</a>

- [Website OpenVPN Cloud](https://openvpn.net/cloud-vpn/?utm_source=terraform&utm_medium=docs)
- [Terraform Registry](https://registry.terraform.io/providers/OpenVPN/openvpn-cloud/latest)

## Description

The Terraform provider for [OpenVPN Cloud](https://openvpn.net/cloud-vpn/?utm_source=terraform&utm_medium=docs) allows teams to configure and update OpenVPN Cloud project parameters via their command line.

## Maintainers

This provider plugin is maintained by:
-	OpenVPN team at [OpenVPN Cloud](https://openvpn.net/cloud-vpn/?utm_source=terraform&utm_medium=docs)
-	SRE Team at [ANNA Money](https://anna.money/?utm_source=terraform&utm_medium=docs) / [GitHub ANNA Money](http://github.com/anna-money/)
- [@patoarvizu](https://github.com/patoarvizu)

## Requirements

-	[Terraform](https://www.terraform.io/downloads.html) 0.12.x
-	[Go](https://golang.org/doc/install) 1.18 (to build the provider plugin)

## Building The Provider

Clone repository to: `$GOPATH/src/github.com/OpenVPN/terraform-provider-openvpn-cloud`

```sh
$ mkdir -p $GOPATH/src/github.com/OpenVPN; cd $GOPATH/src/github.com/OpenVPN
$ git clone git@github.com:OpenVPN/terraform-provider-openvpn-cloud.git
```

Enter the provider directory and build the provider

```sh
$ cd $GOPATH/src/github.com/OpenVPN/terraform-provider-openvpn-cloud
$ make build
```

## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (version 1.18+ is *required*). You'll also need to correctly setup a [GOPATH](http://golang.org/doc/code.html#GOPATH), as well as adding `$GOPATH/bin` to your `$PATH`.

To compile the provider, run `make build`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

```sh
$ make bin
...
$ $GOPATH/bin/terraform-provider-openvpn-cloud
...
```

In order to test the provider, you can simply run `make test`.

```sh
$ make test
```

In order to run the full suite of Acceptance tests, run `make testacc`.

*Note:* Acceptance tests create real resources, and often cost money to run.

```sh
$ make testacc
```

_**Please note:** This provider, like OpenVPN Cloud API, is in beta status. Report any problems via issue in this repo._
