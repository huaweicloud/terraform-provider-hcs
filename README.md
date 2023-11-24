<!-- markdownlint-disable MD001 MD005 MD033 MD013 MD041 -->
<a href="https://www.huaweicloud.com/intl/en-us/product/huaweicloudstack.html"><img width="225px" height="38px" align="right" src="./docs/img/huaweicloudstack_log.png"></a>
<a href="https://www.huaweicloud.com/intl/en-us/product/huaweicloudstack.html"><img width="225px" height="38px" align="left" src="https://camo.githubusercontent.com/1a4ed08978379480a9b1ca95d7f4cc8eb80b45ad47c056a7cfb5c597e9315ae5/68747470733a2f2f7777772e6461746f636d732d6173736574732e636f6d2f323838352f313632393934313234322d6c6f676f2d7465727261666f726d2d6d61696e2e737667"></a>
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
