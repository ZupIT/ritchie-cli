# Ritchie Formula Repo

![Rit banner](/docs/img/ritchie-banner.png)

## Documentation

[Contribute to the Ritchie community](https://github.com/ZupIT/ritchie-formulas/blob/master/CONTRIBUTING.md)

This repository contains rit formulas which can be executed by the [ritchie-cli](https://github.com/ZupIT/ritchie-cli).

- [Gitbook](https://docs.ritchiecli.io)

## Use Formulas

```bash
 rit add repo
 Select your provider:
  > Github
    Gilab
 Repository name: {{some_repo_name}}
 Repository URL: {{this_repo_url}}
 Is a private repository?
    no
  > yes
 Personal access tokens: {{git_personal_token}}
 Select a tag version:
  > 1.0.1
    1.0.0
 Set the priority: 2
```

## Build and test formulas locally

```bash
 rit build formula
```

## Contribute to the repository with your formulas

1. Fork the repository
2. Create a branch: `git checkout -b <branch_name>`
3. Check the step by step of [how to create formulas on Ritchie](https://docs.ritchiecli.io/getting-started/creating-formulas)
4. Add your formulas to the repository
and commit your implementation: `git commit -m '<commit_message>`
5. Push your branch: `git push origin <project_name>/<location>`
6. Open a pull request on the repository for analysis.
