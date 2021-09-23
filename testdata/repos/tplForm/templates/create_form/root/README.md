<!-- markdownlint-disable MD013-->

# Ritchie Formula Repo

![Rit banner](/docs/img/ritchie-banner.png)

## Documentation

This repository contains Ritchie formulas which can be executed by the [ritchie-cli](https://github.com/ZupIT/ritchie-cli).

For more information, check the [Ritchie CLI documentation](https://docs.ritchiecli.io)

## Import repository to use formulas

To import this repository, you need [Ritchie CLI installed](https://docs.ritchiecli.io/getting-started/installation)

Then, you can use the `rit add repo` as below:

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
 Set the priority: X (the lower the number, the higher the priority)
```

Or by using this command:

```bash
rit add repo --provider="Github" --name="{{some_repo_name}}" --repoUrl="{{this_repo_url}}" --priority=1
```

Finally, you can check if the repository has been imported correctly by executing the `rit list repo` command.

## Contribute to this repository with your formulas

### Creating formulas

1. Fork and clone the repository
2. Create a branch: `git checkout -b <branch_name>`
3. Check the step by step of [how to create formulas on Ritchie](https://docs.ritchiecli.io/tutorials/formulas/how-to-create-formulas)
4. Add your formulas to the repository and commit your implementation: `git commit -s -m '<commit_message>'`
5. Push your branch: `git push origin <project_name>/<location>`
6. Open a pull request on the repository for analysis.

### Updating Formulas

1. Fork and clone the repository
2. Create a branch: `git checkout -b <branch_name>`
3. Add the cloned repository to your workspaces (`rit add workspace`) with a highest priority (for example: 1).
4. Check the step by step of [how to implement formulas on Ritchie](https://docs.ritchiecli.io/tutorials/formulas/how-to-implement-a-formula)
and commit your implementation: `git commit -m '<commit_message>`
5. Push your branch: `git push origin <project_name>/<location>`
6. Open a pull request on the repository for analysis.

- [Contribute to the Ritchie community](https://github.com/ZupIT/ritchie-formulas/blob/master/CONTRIBUTING.md)
