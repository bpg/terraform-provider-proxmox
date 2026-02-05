### What does this PR do?
<!--- Clearly state the problem and how this PR solves it. -->

### Contributor's Note
<!---
Please mark the following items with an [x] if they apply to your PR.
Leave the [ ] if they are not applicable, or if you have not completed the item.
--->
- [ ] I have run `make lint` and fixed any issues.
- [ ] I have updated documentation (FWK: schema descriptions + `make docs`; SDK: manual `/docs/` edits).
- [ ] I have added / updated acceptance tests (**required** for new resources and bug fixes — see [ADR-006](docs/adr/006-testing-requirements.md)).
- [ ] I have considered backward compatibility (no breaking schema changes without `!` in PR title).
- [ ] For new resources: I followed the [reference examples](docs/adr/reference-examples.md).
- [ ] I have run `make example` to verify the change works (mainly for SDK / provider config changes).

<!---
You can find more information about coding conventions and local testing in the [CONTRIBUTING.md](https://github.com/bpg/terraform-provider-proxmox/blob/main/CONTRIBUTING.md) file.

If you are unsure how to run `make example`, see [Deploying the example resources](https://github.com/bpg/terraform-provider-proxmox?tab=readme-ov-file#deploying-the-example-resources) section in README.

**PR Title:** Must follow [Conventional Commits](https://www.conventionalcommits.org/) format (e.g., `feat(vm): add clone support`). This becomes the squash commit message.
--->

<!--
*IF* your code contains breaking changes, add `!` before the colon in the PR title, e.g.:
```
feat(vm)!: add support for new feature
```
Also, uncomment the section just below and add a description of the breaking change.
--->

<!---
#### ⚠ BREAKING CHANGES

>>> Put your description here <<<
--->

### Proof of Work
<!---
REQUIRED for code changes. Include at minimum:
- Acceptance test output (`./testacc TestAccYourResource`)
- For bug fixes: test output showing the fix works
- For API changes: either mitmproxy logs showing correct API calls,
  or terraform/tofu output showing successful resource creation/update
  together with the test resource configuration used
Empty proof of work will delay review.
--->

<!--- Please keep this note for the community --->
### Community Note

- Please vote on this pull request by adding a 👍 [reaction](https://blog.github.com/2016-03-10-add-reactions-to-pull-requests-issues-and-comments/) to the original pull request comment to help the community and maintainers prioritize this request
- Please do not leave "+1" or other comments that do not add relevant new information or questions, they generate extra noise for pull request followers and do not help prioritize the request
<!--- Thank you for keeping this note for the community --->

<!--- If your PR fully resolves and should automatically close the linked issue, use Closes. Otherwise, use Relates --->
Closes #0000 | Relates #0000

<!--- Release note for [CHANGELOG](https://github.com/bpg/terraform-provider-proxmox/blob/main/CHANGELOG.md) will be created automatically using the PR's title, update it accordingly. --->
