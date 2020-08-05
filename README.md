[![CircleCI](https://circleci.com/gh/ZupIT/ritchie-cli/tree/master.svg?style=svg)](https://circleci.com/gh/ZupIT/ritchie-cli) [![codecov](https://codecov.io/gh/zupit/ritchie-cli/branch/master/graph/badge.svg)](https://codecov.io/gh/zupit/ritchie-cli) [![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

<img class="special-img-class" src="/docs/img/ritchie-banner.png" />

## Ritchie - One CLI to rule them all

Ritchie is an open source framework that creates and tweaks a CLI for your team. It allows you to easily create, build and share formulas.

This repository contains the CLI core, which can execute formulas stored inside other repositories such as [ritchie-formulas](https://github.com/ZupIT/ritchie-formulas).


## Quick start

### Install rit

- Linux|MacOS

```bash
curl -fsSL https://commons-repo.ritchiecli.io/install.sh | bash
```

- Windows

Download the installer from (https://commons-repo.ritchiecli.io/latest/ritchiecli.msi)

### Initialize rit

Once you made it,  Ritchie will add all community formulas repository and create all the necessary configuration's files.

```bash
rit init
```

### Run your fist formula

After you finished the previous steps - installation and initialization -, you can run a "hello-world" formula to test Ritchie. 
As most of developers like coffee, we created an initial formula that "delivers coffee" to you. 

```bash
rit scaffold generate coffee-go
```


## Full Documentation

- [Gitbook](https://docs.ritchiecli.io)

## Ritchie Legacy-1.x

With the release of version 2.0.0 of Ritchie, the previous version (Ritchie 1.x) has been deprecated. Therefore, only bugs fixes will be implemented in this version.

The legacy code is available at [Ritchie Legacy-1.0.0](https://github.com/ZupIT/ritchie-cli/tree/legacy-1.0.0).

## Contributing

[Contribute to the Ritchie community](https://github.com/ZupIT/ritchie-cli/blob/master/CONTRIBUTING.md)


## Zup Products

- [Zup open source](https://opensource.zup.com.br)

