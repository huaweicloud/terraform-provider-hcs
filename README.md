<!-- markdownlint-disable MD001 MD005 MD033 MD013 MD041 -->
<a href="https://www.huaweicloud.com/intl/en-us/product/huaweicloudstack.html"><img width="225px" height="38px" align="right" src="./docs/img/huaweicloudstack_log.png"></a>
<a href="https://www.huaweicloud.com/intl/en-us/product/huaweicloudstack.html"><img width="225px" height="38px" align="left" src="https://private-user-images.githubusercontent.com/1907997/305699212-b44813ff-c0b6-47a0-9281-d3a5f2d18b6b.svg?jwt=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJnaXRodWIuY29tIiwiYXVkIjoicmF3LmdpdGh1YnVzZXJjb250ZW50LmNvbSIsImtleSI6ImtleTUiLCJleHAiOjE3MDgzMjk5OTIsIm5iZiI6MTcwODMyOTY5MiwicGF0aCI6Ii8xOTA3OTk3LzMwNTY5OTIxMi1iNDQ4MTNmZi1jMGI2LTQ3YTAtOTI4MS1kM2E1ZjJkMThiNmIuc3ZnP1gtQW16LUFsZ29yaXRobT1BV1M0LUhNQUMtU0hBMjU2JlgtQW16LUNyZWRlbnRpYWw9QUtJQVZDT0RZTFNBNTNQUUs0WkElMkYyMDI0MDIxOSUyRnVzLWVhc3QtMSUyRnMzJTJGYXdzNF9yZXF1ZXN0JlgtQW16LURhdGU9MjAyNDAyMTlUMDgwMTMyWiZYLUFtei1FeHBpcmVzPTMwMCZYLUFtei1TaWduYXR1cmU9YWQzNjVmMDI2ZWVmYjcwYWQzMjMxMzZjMTM4YjJkYmE2N2JjMjdlMTUyYzJlZTljNGUyYmIwNDFlY2FjZWZjOCZYLUFtei1TaWduZWRIZWFkZXJzPWhvc3QmYWN0b3JfaWQ9MCZrZXlfaWQ9MCZyZXBvX2lkPTAifQ.mY56JuCJCIMpyhDbvcKb4KfvoPWW2TOqSOetN2F7fm4"></a>
<br/><br/>
<!-- markdownlint-enable MD001 MD005 MD033  MD013 MD041 -->

Huawei Cloud Stack Provider
==============================

<!-- markdownlint-disable-next-line MD034 -->
* Website: https://www.terraform.io
* [![Documentation](https://img.shields.io/badge/documentation-blue)](https://registry.terraform.io/providers/huaweicloud/hcs/latest/docs)
* [![Gitter chat](https://img.shields.io/badge/chat-on_gitter-yellowgreen)](https://gitter.im/hashicorp-terraform/Lobby)
* Mailing list: [Google Groups](http://groups.google.com/group/terraform-tool)

Requirements
------------

* [Terraform](https://www.terraform.io/downloads.html) 0.12.x
* [Go](https://golang.org/doc/install) 1.18 (to build the provider plugin)

Compatibility with HCS
------------

| Provider Version | Huawei Cloud Stack Version  |
|------------------|-----------------------------|
| v2.3.x           | v8.3.0                      |

Building The Provider
---------------------

Clone repository to: `$GOPATH/src/github.com/huaweicloud/terraform-provider-hcs`

```sh
$ mkdir -p $GOPATH/src/github.com/huaweicloudstack; cd $GOPATH/src/github.com/huaweicloudstack
$ git clone https://github.com/huaweicloud/terraform-provider-hcs
```

Enter the provider directory and build the provider

```sh
$ cd $GOPATH/src/github.com/huaweicloudstack/terraform-provider-hcs
$ make build
```

Using the provider
------------------

Please see the documentation at [provider usage](docs/index.md).

Or you can browse the documentation within this repo [here](https://github.com/huaweicloud/terraform-provider-hcs/tree/master/docs).

Developing the Provider
-----------------------

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed
on your machine (version 1.14+ is *required*).
You'll also need to correctly setup a [GOPATH](http://golang.org/doc/code.html#GOPATH),
as well as adding `$GOPATH/bin` to your `$PATH`.

To compile the provider, run `make build`.
This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

```sh
$ make build
...
$ $GOPATH/bin/terraform-provider-hcs
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

License
-------

Terraform-Provider-HCS is under the Mozilla Public License 2.0. See the [LICENSE](LICENSE) file for details.
