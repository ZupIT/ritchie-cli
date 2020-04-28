<!-- Contributing from template (https://github.com/docker/docker.github.io/blob/master/CONTRIBUTING.md) -->

## Contributing

We value your contributions, and we want to make it as easy
as possible to work in this repository. One of the first things to decide is
which branch to base your work on. If you get confused, just ask and we will
help. If a reviewer realizes you have based your work on the wrong branch, we'll
let you know so that you can rebase it.

>**Note**: To contribute code to Ritchie projects, see the
[Ritchie community guidelines](https://docs.ritchiecli.io/community) as well as the 
[Open source contribution guidelines](https://opensource.guide/how-to-contribute/) and our 
[Code of Conduct](https://github.com/ZupIT/ritchie-cli/blob/master/CODE_OF_CONDUCT.md).

## New features for the project

Ritchie is composed of 3 projects which release at different times. 

**Always base your work on the project's `master` branch, naming your new branch according to the following guide :**

<img class="special-img-class" src="/docs/img/git-branchs.png" /> 

**Examples : `feature/name` or `fix/name`**

## Collaborate on a pull request

Unless the PR author specifically disables it, you can push commits into another
contributor's PR. You can do it from the command line by adding and fetching
their remote, checking out their branch, and adding commits to it. Even easier,
you can add commits from the Github web UI, by clicking the pencil icon for a
given file in the **Files** view.

If a PR consists of multiple small addendum commits on top of a more significant
one, the commit will usually be "squash-merged", so that only one commit is
merged in. On occasion this is not appropriate and all commits will be kept
separate when merging.ficant one, the commit will usually be "squash-merged", so that only one commit is merged in. 
On occasion this is not appropriate and all commits will be kept separate when merging.

## Pull request guidelines

Help us review your PRs more quickly by following these guidelines.

- Try not to touch a large number of files in a single PR if possible.

- Don't change whitespace or line wrapping in parts of a file you are not
  editing for other reasons. Make sure your text editor is not configured to
  automatically reformat the whole file when saving.

- Reviewers will check the staging site and contact you to fix eventual problems.

If you can think of other ways we could streamline the review process, let us
know.

## Style guide

Ritchie does not currently maintain a style guide. Use your best judgment, and
try to follow the example set by the existing files.