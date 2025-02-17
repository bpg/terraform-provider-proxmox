# Security Policy

## Supported Versions

We release patches for security vulnerabilities for the following versions:

| Version                           | Supported          |
|-----------------------------------|--------------------|
| Latest release                    | :white_check_mark: |
| Previous minor release            | :white_check_mark: |
| Older than previous minor release | :x:                |

## Reporting a Vulnerability

We take the security of the Terraform Provider for Proxmox and its users seriously. If you believe you have found a security vulnerability, please report it to us privately.

**Please do not report security vulnerabilities through public GitHub issues, discussions, or pull requests.**

Instead, please send an email to [bpg.github.com.tn75g@passmail.net](mailto:bpg.github.com.tn75g@passmail.net) or report it through GitHub's Security Advisory feature:

1. Go to <https://github.com/bpg/terraform-provider-proxmox/security/advisories/new>
2. Provide a descriptive title
3. Fill in a detailed description of the issue
4. Click "Submit report"

Please include the following information in your report:

- Description of the vulnerability
- Steps to reproduce the issue
- Potential impact of the vulnerability
- Suggested fix if possible
- Your name/handle for credit (optional)

## What to Expect

When you submit a vulnerability report, you can expect:

- Acknowledgment of your report within 48 hours
- Regular updates about our progress
- Credit for discovering the vulnerability (if desired)

## Disclosure Policy

When we receive a security bug report, we will:

1. Confirm the problem and determine the affected versions
2. Audit code to find any potential similar problems
3. Prepare fixes for all supported versions
4. Release new versions and notify users

## Security-Related Configuration

When using this provider, please follow these security best practices:

1. **API Credentials**:
   - Use environment variables or encrypted credential files to store sensitive information
   - Never commit API tokens or credentials to version control
   - Use the most restrictive permissions possible for API users

2. **Network Security**:
   - Use HTTPS/TLS for all API connections
   - Configure appropriate firewall rules
   - Use private networks where possible

3. **State File Security**:
   - Encrypt your Terraform state files
   - Use remote state with appropriate access controls
   - Be cautious with state file contents as they may contain sensitive information
