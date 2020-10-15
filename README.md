[![CircleCI](https://circleci.com/gh/ZupIT/ritchie-cli/tree/master.svg?style=svg)](https://circleci.com/gh/ZupIT/ritchie-cli) 
[![codecov](https://codecov.io/gh/ZupIT/ritchie-cli/branch/master/graph/badge.svg)](https://codecov.io/gh/ZupIT/ritchie-cli)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

<img class="special-img-class" src="/docs/img/ritchie-banner.png"  alt="Ritchie logo with the phrase: Keep it simple"/>

## Your customizable automation tool

**Ritchie CLI** is an open source product that allows you to create, store and share any kind of automations, executing them through command lines, to run operations or start workflows.

This repository contains the CLI core, which can execute **formulas** stored inside other repositories such as [ritchie-formulas](https://github.com/ZupIT/ritchie-formulas) or [ritchie-formulas-demo](https://github.com/ZupIT/ritchie-formulas-demo)

In Ritchie's context, **a formula is a script** that can be executed automatically or interactively through a command line. 

Adapting an existing script to Ritchie structure allows you to run it **locally** or through **Docker**, and to share it on a **Github** or **Gitlab** repository.

<img class="special-img-class" src="/docs/img/formulas-explanation.png" alt="Formulas explanation"/>

## Full Documentation

[Gitbook](https://docs.ritchiecli.io)

## Quick start

### Install rit

- Linux|MacOS

```bash
curl -fsSL https://commons-repo.ritchiecli.io/install.sh | bash
```

- Windows

Download the installer from (https://commons-repo.ritchiecli.io/latest/ritchiecli.msi)

### Initialize rit

```bash
rit init
```

*Note: You need to import the **commons** repository to be able to create formulas.*

### Run your first formula

To access the ["hello-world" formula]((https://github.com/ZupIT/ritchie-formulas-demo/tree/master/demo/hello-world)), you'll need to add the [ritchie-formulas-demo](https://github.com/ZupIT/ritchie-formulas-demo) repository locally. To do so, you can use the `rit add repo` command, or execute the command line below:

```bash
echo '{"provider":"Github", "name":"demo", "version":"2.0.0", "url":"https://github.com/ZupIT/ritchie-formulas-demo", "token": null, "priority":1}' | rit add repo --stdin
```

Then, you'll be able to execute the "hello-world" formula through the command line below:

```bash
rit demo hello-world
```

## Ritchie Legacy-1.x

With the release of version 2.0.0 of Ritchie, the previous version (Ritchie 1.x) has been deprecated. Therefore, only bugs fixes will be implemented in this version.

The legacy code is available at [Ritchie Legacy-1.0.0](https://github.com/ZupIT/ritchie-cli/tree/legacy-1.0.0).

## Contributing

[Contribute to the Ritchie community](https://github.com/ZupIT/ritchie-cli/blob/master/CONTRIBUTING.md)

## Zup Products

[Zup open source](https://opensource.zup.com.br)
