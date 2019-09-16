<!-- 🎉🎉🎉 Thank you for the PR!!! 🎉🎉🎉 -->


Relates to _in case of new feature, this should point to issue/(s) which describes the feature_

Fixes _in case of a bug fix, this should point to a bug and any other related issue(s)_

Should merge before : _list any dependent PRs_

# Description

<!-- Describe your changes here- ideally you can get that description straight from
your descriptive commit message(s)! -->

# Output Changes
Write n/a if not output or log lines added

## Before
Paste skaffold output before your change

## After
Paste skaffold output after your change

## Next PRs.

In this section describe a list of follow up PRs if the current PR is a part of big feature change.

See example #2811

Write n/a if not applicable.


# Submitter Checklist

These are the criteria that every PR should meet, please check them off as you
review them:

- [ ] Includes [unit tests](../DEVELOPMENT.md#creating-a-pr)
- [ ] Mentions any output changes. 

_See [the contribution guide](../CONTRIBUTING.md) for more details._

Double check this list of stuff that's easy to miss:

- If you are adding [a example to the `examples` dir](../examples), please copy them to [`integration/examples`](../integration/examples)
- Every new example added in [`integration/examples`](../integration/examples) dir should be tested in [integration test](../integration)

## Reviewer Notes

- [ ] The code flow looks good. 
- [ ] Unit test added.
- [ ] User facing document changes look good.


# Release Notes

Describe any user facing changes here so maintainer can include it in the release notes, or delete this block.

```
Examples of user facing changes:
- Skaffold config changes like
  e.g. "Add buildArgs to `Kustomize` deployer skaffold config."
- Bug fixes 
  e.g. "Improve skaffold init behaviour when tags are used in manifests"
- Any changes in skaffold behavior
  e.g. "Artiface cachine is turned on by default."

```