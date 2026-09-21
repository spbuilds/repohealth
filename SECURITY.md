# Security Policy

## Supported Versions

Security fixes are made on the `main` branch and shipped in the next release.
Only the [latest release](https://github.com/spbuilds/repohealth/releases/latest)
is supported; older releases do not receive fixes.

## Reporting a Vulnerability

Please do not report security vulnerabilities through public GitHub issues.

Use GitHub's private vulnerability reporting instead: open the repository's
[Security tab](https://github.com/spbuilds/repohealth/security) and choose
**Report a vulnerability**, or go directly to
<https://github.com/spbuilds/repohealth/security/advisories/new>.
The report is visible only to the maintainer until a fix is published.

Please include a description of the issue, the RepoHealth version
(`repohealth --version`) and steps to reproduce; a minimal repository layout
that triggers the problem is ideal.

We aim to acknowledge reports within seven days. Confirmed issues are fixed on
`main` and released as soon as practical, with credit to the reporter unless
they prefer otherwise.

## Scope

RepoHealth is a local, offline analysis tool. It:

- Reads files from the repository directory
- Runs `git` commands for activity metrics
- Does not send data to external services
- Does not modify any files in the analyzed repository

Issues such as unsafe file handling, unintended command execution or content
injection in generated HTML reports belong in the private process above.
Questions about the accuracy of a check's result are welcome as ordinary
GitHub issues.
