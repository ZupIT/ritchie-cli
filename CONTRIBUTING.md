<!-- Contributing from template (https://github.com/docker/docker.github.io/blob/master/CONTRIBUTING.md) -->

# **Contributing Guide**

This is Ritchie's contributing guide. Please read the following sections to learn how to ask questions and how to work on something. We value your contributions, and we want to make it as easy as possible for you to work in this repository.

## **Table of contents**

### 1. [**Before you contribute**](#Before-you-contribute)
> #### i. [**Code of Conduct**](#Code-of-Conduct)
> #### ii. [**Legal**](#Legal)
### 2. [**Prerequisites**](#Prerequisites)
> #### i. [**Developer Certificate of Origin - DCO**](#Developer-Certificate-of-Origin-DCO)
> #### ii.  [**Choose a branch**](#Choose-a-branch)
> #### iii. [**Check out the guidelines**](#Check-out-the-guidelines)
### 3. [**How to contribute?**](#How-to-contribute?)
> #### i. [**Add new feature, bugfixing or improvement**](#Add-new-feature-bugfixing-or-improvement)
> #### ii. [**Pull request guidelines**](#Pull-request-guidelines)
> #### iii. [**Pull request review guidelines**](#Pull-request-review-guidelines)
> #### iv. [**Collaborate on a pull request**](#Collaborate-on-a-pull-request)
> #### v. [**Tests guide**](#Tests-guide)
> #### vi. [**Opening a new issue**](#Opening-a-new-issue)
### 4. [**Community**](#Community)



## **Before you contribute**

### **Code of Conduct**
Please follow the [**Code of Conduct**](https://github.com/ZupIT/ritchie-cli/blob/main/CODE_OF_CONDUCT.md) in all your interactions with our project.

### **Legal**
- Ritchie is licensed over [**ASF - Apache License**](https://github.com/ZupIT/ritchie-cli/blob/main/LICENSE), version 2, so new files must have the ASL version 2 header. For more information, please check out [**Apache license**](https://www.apache.org/licenses/LICENSE-2.0).

- All contributions are subject to the [**Developer Certificate of Origin (DCO)**](https://developercertificate.org). 
When you commit, use the ```**-s -S** ``` option to include the Signed-off-by line at the end of the commit log message.

## **Prerequisites**
Check out the requisites before contributing to Ritchie:

### **Developer Certificate of Origin - DCO**

 This is a security layer for the project and for the developers. It is mandatory.
 
 Follow one of these two methods to add DCO to your commits:
 
**1. Command line**
 Follow the steps: 
 **Step 1:** Configure your local git environment adding the same name and e-mail configured at your GitHub account. It helps to sign commits manually during reviews and suggestions.

 ```
git config --global user.name “Name”
git config --global user.email “email@domain.com.br”
```
**Step 2:** Add the Signed-off-by line with the `'-s -S'` flag in the git commit command:

```
$ git commit -s -S -m "This is my commit message"
```
**2. GitHub website**
You can also manually sign your commits during GitHub reviews and suggestions, follow the steps below: 

**Step 1:** When the commit changes box opens, manually type or paste your signature in the comment box, see the example:

```
$ git commit -m “My signed commit” Signed-off-by: username <email address>
```
For this method, your name and e-mail must be the same registered to your GitHub account.

### **Choose a branch**
One of the first things to decide is which branch to base your work on. If you get confused, just ask and we will help you. If a reviewer realizes you have based your work on the wrong branch, we'll let you know so that you can rebase it.

### **Check out the guidelines**
If you want to contribute code to Ritchie projects, check out:
- [**Ritchie Community Guidelines**](https://docs.ritchiecli.io/faq#community)
- [**Open source contribution guidelines**](https://opensource.guide/how-to-contribute/)

## **How to contribute?** 
See the guidelines to submit your changes: 

### **Add new feature, bugfixing or improvement**
If you want to add an improvement, a new feature or bugfix, follow the steps to contribute: 

**Step 1:** **Fork the project**. 

**Step 2:** Always base your work on the project's **`main`** branch, naming your new branch according to the following guide:

<img class="special-img-class" src="/docs/img/git-branchs.png" /> 

**Examples : `feature/name` or `fix/name`**

**Step 3:** Make your changes and open a GitHub pull request;

### **Pull request guidelines**

Follow these guidelines to help us review your PRs quicker: 

- Try not to touch a large number of files in a single PR if possible.

- Provide a clear description of the motivation of your PR, this is a large
  project, so context is important.

- Don't change whitespace or line wrapping in parts of a file you are not
  editing for other reasons/unknowingly. Make sure your text editor is not configured to automatically reformat the whole file when saving.

- Reviewers will check the staging site and contact you to fix any problems.

- If you agree with the suggested comment, just resolve or react with an emoji to it. You don't need to write a confirmation. The code owner's mailbox will thank you later. :smiley:

>**Note**: If you know other ways we could streamline the review process, [**let us
know**](https://forum.zup.com.br/c/en/9).

### **Pull request review guidelines**

If you want to contribute on the quality of the incoming code, we would appreciate if you follow these suggestions: :smile:

- Be suggestive, never impose a correction or criticize your peer. Instead of "change this code", go for "what do you think about implementing it this way?"

- Explain why you suggested such a correction, a change without meaning might not be productive. 
The author has all the right to counterargue a comment if they think it is the best for the project. Provide a clear technical or business justification and even links or references if possible. Everybody loves to learn something new about coding.

- Sometimes the literal answer might not be necessary. Instead of pasting the solution _verbatim_, provide the right direction and let the author figure out the solution.

- Always have a sense of community and try to help others, because they are trying to help us.

- If a debate gets more heated in a review, try to set a call or meeting to clarify points. Letting it get out of hand might affect the Ritchie community in general.

- Requesting changes sometimes is perceived as a harsh action for some engineers. Try to do it with parsimony, usually when you spot a production-breaking change or to prevent an already approved PR from being merged without that last important modification to code.

### **Collaborate on a pull request**

Unless the PR author specifically disables it, you can push commits into another
contributor's PR. See how you can do it:
- From the command line by adding and fetching their remote;
- Check out their branch, and add commits to it; 
- You can add commits from the Github web UI, by clicking the pencil icon for a
given file in the **Files** view.

If a PR consists of multiple small addendum commits on top of a more significant
one, the commit will usually be "squash-merged", so that only one commit is
merged in. Unless the new addendum commit is a significant one, the commit will usually be "squash-merged", so that only one commit is merged in.
On occasion, this is not appropriate and all commits will be kept separate when merging.


## **Tests guide**

To keep Ritchie easy to maintain we need to have tests. Use the following command to run:
```
make unit-test:<name-of-test>
make functional-test:<name-of-test>
```

## **Opening a new issue**

If you want to contribute with something that doesn't have any **issue** yet, you can:
- Create a new detailed **issue** [**in the repository**](https://github.com/ZupIT/ritchie-cli/issues/new/choose). 
- Choose one template to fill in:
  - Bug Report
  - Feature request
  - Improvement
  - Support request
  - Report a security vulnerability 
Then, you are able to solve it.

- Our team will evaluate and add a **Hackoberfest** label, this will allow you to participate in the event to solve your own **issue**.

## **Community**

- Do you have any questions about Ritchie? Let's chat in our [**forum**](https://forum.zup.com.br/c/en/9).

Thank you for your contribution!

**Ritchie team** 
