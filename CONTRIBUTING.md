<!-- Contributing from template (https://github.com/docker/docker.github.io/blob/master/CONTRIBUTING.md) -->

## Contributing

We value your contributions, and we want to make it as easy
as possible to work in this repository. One of the first things to decide is
which branch to base your work on. If you get confused, just ask and we will
help. If a reviewer realizes you have based your work on the wrong branch, we'll
let you know so that you can rebase it.

>**Note**: To contribute code to Ritchie projects, see the
[Ritchie Community Guidelines](https://docs.ritchiecli.io/faq#community) as well as the 
[Open source contribution guidelines](https://opensource.guide/how-to-contribute/) and our 
[Code of Conduct](https://github.com/ZupIT/ritchie-cli/blob/master/CODE_OF_CONDUCT.md).

## New features for the project

Ritchie is composed of 3 projects which release at different timelines.

**Always base your work on the project's `master` branch, naming your new branch according to the following guide :**

![Git flow](https://user-images.githubusercontent.com/22433243/125490327-6a00861e-98ff-4b86-ba46-70aa90906edc.png)

**Examples : `feature/name` or `fix/name`**

## Collaborate on a pull request

Unless the PR author specifically disables it, you can push commits into another
contributor's PR. You can do it from the command line by adding and fetching
their remote, checking out their branch, and adding commits to it. Even easier,
you can add commits from the Github web UI, by clicking the pencil icon for a
given file in the **Files** view.

If a PR consists of multiple small addendum commits on top of a more significant
one, the commit will usually be "squash-merged", so that only one commit is
merged in. Unless the new addendum commit is a significant one, the commit will usually be "squash-merged", so that only one commit is merged in.
On occasion this is not appropriate and all commits will be kept separate when merging.

## Pull request guidelines

Help us review your PRs more quickly by following these guidelines.

- Try not to touch a large number of files in a single PR if possible.

- Provide a clear description of the motivation of your PR, this is a large
  project, so context is important

- Don't change whitespace or line wrapping in parts of a file you are not
  editing for other reasons / unknowingly. Make sure your text editor is not configured to
  automatically reformat the whole file when saving.

- Reviewers will check the staging site and contact you to fix any problems.

- If you agree with the suggested comment, just resolve or react with an emoji to it. No need to write a confirmation. 
The code owner's mailbox will thank you later :smiley:

If you can think of other ways we could streamline the review process, let us
know.

## Pull review guidelines

For those wanting to contribute on the quality of the incoming code, try to follow these
suggestions for a happy community :smile:

- Be suggestive, never impose a correction or criticize your peer. Instead of "change this code",
go for something more like "what do you think about implementing it this way?"

- Explain why you suggested such a correction, a change without meaning might not be productive. 
The author has all the right to counterargue a comment if he/she thinks it is for the best of the project.
Provide a clear technical or business justification and even links or references if possible. 
Everybody loves to learn something new about coding to become a better developer.

- Sometimes the literal answer might not be necessary. Instead of pasting the solution _verbatim_,
provide the right direction and let the author figure out the solution.

- Always have a sense of community and try to help others, because they are trying to help us.

- If a debate gets more heated in a review, try to set a call or meeting to clarify points. Letting it get out of hand 
might affect the Ritchie community in general.

- Requesting changes sometimes is perceived as a harsh action for some engineers. Try to do it with parsimony, 
usually when you spot a production-breaking change or to prevent an already approved PR from being merged without
that last important modification to code.

## Style guide

Ritchie does not currently maintain a style guide. Use your best judgement, and
try to follow the example set by the existing files.

## Tests guide

To keep Ritchie easy to maintain is crucial to have tests.
To run:
```
make unit-test:<name-of-test>
make functional-test:<name-of-test>
```

## Hacktoberfest

If you want to contribute with something thay doesn't have any ISSUE yet, you can create a new detailed ISSUE [here](https://github.com/ZupIT/ritchie-cli/issues/new/choose) for our team to evaluate, and they will eventually add a Hackoberfest label to allow you to participate in the event for resolving this ISSUE.
