<div align="center">

![grafton](.github/grafton.png)

**Manifold's provider validation tool**

[Manifold's Provider Documentation](https://provider-api-docs.manifold.co/?tab=provider) |
[Code of Conduct](./.github/CONDUCT.md) |
[Contribution Guidelines](./.github/CONTRIBUTING.md) |
[CHANGELOG](./CHANGELOG.md)

[![GitHub release](https://img.shields.io/github/tag/manifoldco/grafton.svg?label=latest)](https://github.com/manifoldco/grafton/releases)
[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg)](https://godoc.org/github.com/manifoldco/grafton)
[![Travis](https://img.shields.io/travis/manifoldco/grafton/master.svg)](https://travis-ci.org/manifoldco/grafton)
[![Go Report Card](https://goreportcard.com/badge/github.com/manifoldco/grafton)](https://goreportcard.com/report/github.com/manifoldco/grafton)
[![License](https://img.shields.io/badge/license-BSD-blue.svg)](./LICENSE.md)

</div>

The [Manifold's Provider Documentation](https://provider-api-docs.manifold.co/?tab=provider) is
the best place to get started on adding your product to the Manifold's Marketplace.

## Introduction

Grafton is a super simple CLI tool used by service providers (and Manifold) to
test their integrations with Manifold.

Using this CLI tool, a provider can test the following:

- Provisioning/Deprovisioning of Resources
- Provisioning/Deprovisioning of Credentials
- Rotation of Credentials
- Resizing of Resources
- Pulling [Resource Measures](#excluding-features) (optional)

Requests generated by Grafton are the exact same shape and style as those made
by the Provisioning service. Using this tool, a provider should be able to
implement and test their integration without needing real requests from Manifold.

## Installation

Precompiled binaries are build for each release of Grafton for the following
platforms:
- `linux / amd64`
- `darwin (macos) / amd64`
- `windows / amd64`

The zip files of these binaries are in the downloads section of each release
(i.e. the
[latest release](https://github.com/manifoldco/grafton/releases/latest)).

To install, download and unzip the appropriate binary for your system, and run
the `grafton` program inside. You may want to add it to a directory in your
`PATH` for ease of use.

## Usage

1. **Generate Keypair**

To use, a provider must generate a local public keypair used for signing
ephemereal live keypairs. This allows `Grafton` to generate and sign requests
using the same algorithms as Manifold.

```
$ grafton generate
```

When a keypair is generated it will be written to a `masterkey.json` file in
the current working directory.

The public key contained in this file can be used by the implemented service
to verify the authenticity of the requests made by Grafton


2. **Run Tests**

Now that a key has been generated, you can run the tests which will provision a
resource, create credentials, create another set of credentials, and then
deprovision those before resizing the resource. Finally, Grafton will
deprovision the resource.

```
grafton test --product=bonnets --plan=small --region=aws::us-east-1 \
    --client-id=21jtaatqj8y5t0kctb2ejr6jev5w8 \
    --client-secret=3yTKSiJ6f5V5Bq-kWF0hmdrEUep3m3HKPTcPX7CdBZw \
    --connector-port=3001 \
    --new-plan=large \
    http://localhost:4567
```

![Grafton test output](.github/grafton-test-output.png)

3. **Mini-Marketplace Connector API**

Grafton provides a simple version of Manifold Marketplace that allows to test provision, deprovision and SSO manually.

```
grafton serve --product=bonnets --plan=simple-hood --region=east-coast --provider-api=http://yourlocalserver/v1
```

### Excluding Features

When testing it is possible to exclude of one more features from being run. To
disable Resource Measures and Resize, for example, you can pass:
`--exclude resource-measures --exclude plan-change`.

A full list of available tests you can exclude (which is all of them):
- `cleanup`
- `credentials`
- `resource-measures`
- `provision`
- `plan-change`
- `sso`
- `credential-rotation`

_Note_ : resource-measures is a test you are ONLY required to pass if you are using metered pricing. If you are not, you can exclude it.

## Developing

### Backward compatibility

We strongly enforce Grafton changes to be backward compatible, which means changes
to the API should not break existing clients. If a breaking change is necessary, it
should be done either as a new endpoint or a new spec version, eg: v2.

In order to guarantee compatibility, we expect any API change to also be implemented
for our [go example client](https://github.com/manifoldco/go-sample-provider) first,
before the change can be approved in Grafton. Our CI takes cares of getting the latest
version of the client and test against the new changes.

### OpenAPI spec

`provider.yaml` is a manual copy of an internally generated file. It, in
turn, is used to generate the client code under `generated`.

Until we automate more of this, the following has to be kept in mind:
- Any changes to `provider.yaml` should not be merged here directly, but
  instead ported to the internal file first.
- If the internal source of `provider.yaml` changes, the updates must be brought
  over here.
- Whenever `provider.yaml` changes, `make generated-clients` should be run, and the
  changes to the generated code should be checked in.

### Releasing

Releasing grafton uses [promulgate](https://github.com/manifoldco/promulgate).

1. Ensure all of the required code has been merged into master.
2. Determine whether or not the `CHANGELOG.md` is up to date and correct for
   this release, if it's not, update it!
3. Create a tag off master (following [semver](http://semver.org/)), which
   includes the up to date `CHANGELOG.md` (and matches the remote sha of master
   on github).
4. Done! Promulgate via Travis will take care of creating the binaries and
   uploading the zip files.
