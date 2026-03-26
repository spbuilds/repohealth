# Security Policy

## Reporting a Vulnerability

If you discover a security vulnerability in RepoHealth, please report it by
opening a GitHub Issue with the `[SECURITY]` prefix in the title.

RepoHealth is a local analysis tool that does not handle sensitive data,
authentication, or network services. However, we take all security reports
seriously.

## Scope

RepoHealth runs locally on your machine and:

- Reads files from the repository directory
- Executes `git` commands for activity metrics
- Does not send data to external services
- Does not modify any files in the analyzed repository
