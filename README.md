[![CircleCI](https://circleci.com/gh/ZupIT/ritchie-cli/tree/master.svg?style=svg)](https://circleci.com/gh/ZupIT/ritchie-cli) 
[![codecov](https://codecov.io/gh/ZupIT/ritchie-cli/branch/master/graph/badge.svg)](https://codecov.io/gh/ZupIT/ritchie-cli)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

<img class="special-img-class" src="/docs/img/ritchie-banner.png"  alt="Ritchie logo with the phrase: Keep it simple"/>

# Summary

- [All your automations in one place](https://github.com/ZupIT/ritchie-cli#all-your-automations-in-one-place)
  - [Usage](https://github.com/ZupIT/ritchie-cli#usage)
  - [A customizable CLI automation tool](https://github.com/ZupIT/ritchie-cli#a-customizable-cli-automation-tool)
- [Full Documentation](https://github.com/ZupIT/ritchie-cli#-full-documentation)
- [Quick Start](https://github.com/ZupIT/ritchie-cli#-quick-start-)
  - [Install rit latest version](https://github.com/ZupIT/ritchie-cli#1%EF%B8%8F⃣-install-rit-latest-version)
  - [Initialize rit locally](https://github.com/ZupIT/ritchie-cli#2%EF%B8%8F⃣-initialize-rit-locally)
  - [Add your first formulas repository](https://github.com/ZupIT/ritchie-cli#3%EF%B8%8F⃣-add-your-first-formulas-repository)
  - [Run the Hello World formula](https://github.com/ZupIT/ritchie-cli#4%EF%B8%8F⃣-run-the-hello-world-formula)
- [Contributing](https://github.com/ZupIT/ritchie-cli#-contributing-to-ritchie)
  - [Zup Open Source projects](https://github.com/ZupIT/ritchie-cli#zup-open-source-projects)
  - [Zup Open Source forum](https://github.com/ZupIT/ritchie-cli#zup-open-source-forum)

## All your automations in one place

**Ritchie CLI** is an open source project that allows to **create**, **store** and **share** automations, executing them through command lines.

### Usage

<p align="center">
  <a rel="noopener" target="_blank"><img width="600px" src="https://user-images.githubusercontent.com/22433243/121697854-d95a8400-caa3-11eb-9fd3-a3776f6e5e12.mp4" alt="video containing the command demonstration"></a>
</p>

### A customizable CLI automation tool

This repository contains the CLI core, which can execute **formulas** stored inside other repositories such as [ritchie-formulas](https://github.com/ZupIT/ritchie-formulas) or [ritchie-formulas-demo](https://github.com/ZupIT/ritchie-formulas-demo).

In Ritchie's context, **a formula is a script** that can be executed automatically or interactively through a command line.

Adapting an existing script to Ritchie structure allows you to run it **locally** or through **Docker**, and to share it on a **Git** repository.

<img class="special-img-class" src="/docs/img/formulas-explanation.png" alt="Formulas explanation"/>

## 📚 [Full Documentation](https://docs.ritchiecli.io)

[![Documentation](/docs/img/documentation-ritchie.png)](https://docs.ritchiecli.io)

## 🚀 Quick start 🤖

### 1️⃣ Install rit latest version

#### Linux

```bash
curl -fsSL https://commons-repo.ritchiecli.io/install.sh | bash
```

#### MacOS

```bash
curl -fsSL https://commons-repo.ritchiecli.io/install.sh | bash
```

#### Windows

- Download the installer from [ritchiecli.msi](https://commons-repo.ritchiecli.io/latest/ritchiecli.msi)

- Using Winget:

```bash
winget install Ritchie-CLI
```

You can also download rit **packages** or **specific versions** according to the OS [on the documentation](https://docs.ritchiecli.io/getting-started/install-cli)

### 2️⃣ Initialize rit locally

```bash
rit init
```

***Note**: You need to import the **commons** repository to be able to create formulas.*

***Sharing metrics anonymously** will help us improving the tool.
For any question, check our [privacy policy](https://www.zup.com.br/politica-de-privacidade/politica-ritchie#politicas).*

### 3️⃣ Add your first formulas repository

To access the ["hello-world" formula]((https://github.com/ZupIT/ritchie-formulas-demo/tree/master/demo/hello-world)), you need to add the [ritchie-formulas-demo](https://github.com/ZupIT/ritchie-formulas-demo) repository locally.

To do so, you can use the `rit add repo` command **manually** on your terminal, or execute the command line below with **input flags**:

```bash
rit add repo --provider="Github" --name="demo" --repoUrl="https://github.com/ZupIT/ritchie-formulas-demo" --priority=1
```

### 4️⃣ Run the Hello World formula

Execute the "hello-world" formula through the command line below:

```bash
rit demo hello-world
```

***Note**: This formula has been implemented using **Golang**, so to use it **locally** you'll need Golang to be installed on your machine. If you don't have or don't want to install Golang, you can use the same command with **Docker**:*

```bash
rit demo hello-world --docker
```

## 🤝 Contributing to Ritchie

- [Guidelines](https://github.com/ZupIT/ritchie-cli/blob/master/CONTRIBUTING.md)

- [Developer Guide](https://github.com/ZupIT/ritchie-cli/blob/master/DEVELOPER_GUIDE.md)

- [Documentation repository](https://github.com/ZupIT/docs-ritchie)

### [Zup Open Source Projects](https://opensource.zup.com.br)

[![Zup open source](/docs/img/zup-open-source.png)](https://opensource.zup.com.br)

### [Zup Open Source Forum](https://forum.zup.com.br/c/en/9)

[![Zup forum](/docs/img/zup-forum-topics.png)](https://forum.zup.com.br/c/en/9)
