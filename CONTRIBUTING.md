<!-- Contributing from template (https://github.com/docker/docker.github.io/blob/master/CONTRIBUTING.md) -->

## Contributing

We value your contributions, and we want to make it as easy
as possible to work in this repository. One of the first things to decide is
which branch to base your work on. If you get confused, just ask and we will
help. If a reviewer realizes you have based your work on the wrong branch, we'll
let you know so that you can rebase it.

>**Note**: To contribute code to Ritchie projects, see the
[Ritchie community guidelines](https://docs.ritchiecli.io/v/doc-english/community) as well as the 
[Open source contribution guidelines](https://opensource.guide/how-to-contribute/).

## New features for the project

Ritchie is composed of 3 projects which release at different times. **If, and only if,
your pull request relates to a currently unreleased feature of a project, base
your work on that project's `vnext` branch.** 

These branches were created bycloning `master` and then importing a project's `master` branch's 
implementation into it, in a way that preserved the commit history. 

When a project has a release, its `vnext` branch will be merged into `master`.

- **[vnext-cli](https://github.com/ZupIT/ritchie-cli/tree/vnext-cli):**
  implementation for upcoming features in the [ritchie/cli](https://github.com/ZupIT/ritchie-cli)
  project

- **[vnext-server](https://github.com/ZupIT/ritchie-cli/tree/vnext-server):**
  implementation for upcoming features in the [ritchie/server](https://github.com/ZupIT/ritchie-server)
  project
  
- **[vnext-formulas](https://github.com/ZupIT/ritchie-cli/tree/vnext-formulas):**
  implementation for upcoming features in the [ritchie/formulas](https://github.com/ZupIT/ritchie-formulas)
  project

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

- A Jenkins pipeline runs for each PR that is against one of our long-lived
  branches like `master` and the `vnext` branches, and deploys the result of
  your PR to a staging site. The URL will be available at the bottom of the PR
  in the **Conversation** view. Check the staging site for problems and fix them
  if necessary. Reviewers will check the staging site too.

If you can think of other ways we could streamline the review process, let us
know.

## Style guide

Ritchie does not currently maintain a style guide. Use your best judgment, and
try to follow the example set by the existing files.