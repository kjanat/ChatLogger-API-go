# Security Policy

This document outlines security procedures and policies for the ChatLogger API (Go) project.

## Supported Versions

We maintain security updates for the following versions:

| Version | Supported          | Notes                  |
| ------- | ------------------ | ---------------------- |
| 0.4.x   | :white_check_mark: | Current stable release |
| 0.3.x   | :x:                | Not maintained         |
| 0.2.x   | :x:                | Not maintained         |
| 0.1.x   | :x:                | No longer supported    |

## Reporting a Vulnerability

The ChatLogger API team takes security vulnerabilities seriously. We appreciate your efforts to responsibly disclose your findings.

### How to Report a Security Vulnerability

1. **DO NOT** create public GitHub issues for security vulnerabilities
2. Email security concerns to `security-chatlogger@kjanat.com` with the following details:
   - Description of the vulnerability
   - Steps to reproduce
   - Potential impact
   - Suggested mitigation or fix if available
3. Expect an initial response within 48 hours acknowledging your report
4. Our team will work with you to understand and validate the issue

### What to Expect

- Acknowledgment of your report within 48 hours
- Regular updates on our progress (at least every 72 hours)
- An assessment of the vulnerability and its impact
- A plan for mitigation and release of a fix

### Disclosure Policy

- Vulnerabilities will remain private until a fix is ready
- Once fixed, we'll acknowledge your contribution (unless you prefer to remain anonymous)
- We'll credit you in the release notes and CHANGELOG.md (if you consent)
- Public disclosure will occur after our users have had sufficient time to update

## Security Controls and Best Practices

### Authentication & Authorization

This API implements a dual authentication strategy:

#### 1. JWT-Based Authentication (Dashboard Users)

- JWT tokens are delivered via secure HTTP-only cookies
- Tokens have appropriate expiration settings
- Users have role-based access controls (superadmin, admin, user, viewer)

#### 2. API Key Authentication (Chat Plugins)

- API keys are stored hashed, never in plaintext
- Each API key is scoped to a specific organization
- Keys can be revoked by organization admins

### Password Security

- All user passwords are hashed using bcrypt with appropriate cost factors
- Password requirements enforce minimum complexity standards
- Rate limiting is applied to login attempts

### Container Security

- Docker images run as non-root user (`chatlogger`)
- Container images are signed using Sigstore cosign
- CI/CD pipeline verifies image signatures

### Dependency Management

- Dependencies are regularly updated
- Vulnerability scanning is part of the CI/CD pipeline

## Security Recommendations for Deployment

When deploying this API in production, we recommend:

1. **Use HTTPS/TLS**: Always deploy behind TLS with modern cipher configurations
2. **Environment Security**:
   - Use secret management solutions instead of environment variables where possible
   - Never commit `.env` files with secrets to version control
3. **Database Security**:
   - Use least-privilege database accounts
   - Enable TLS for database connections
   - Consider data encryption at rest for sensitive fields
4. **Network Security**:
   - Use a reverse proxy or API gateway with rate limiting
   - Configure appropriate CORS settings
   - Consider a WAF for additional protection
5. **API Key Management**:
   - Rotate API keys periodically
   - Use unique API keys per integration
6. **JWT Configuration**:
   - Use strong, unique JWT secrets
   - Configure reasonable expiration times
7. **Image Verification**:
   - Verify container image signatures as part of your deployment process

## For Developers

If you're contributing to the project, please adhere to these security practices:

1. **Validate all inputs**: Always validate user input at the handler level
2. **Parameterize queries**: Never use string concatenation for SQL queries
3. **Follow least privilege**: Only request the permissions you need
4. **Secrets handling**: Never hardcode secrets or store them in version control
5. **Error handling**: Don't expose sensitive information in error messages
6. **Dependencies**: Keep dependencies updated and minimize third-party packages
7. **Code reviews**: All PRs should undergo security-focused code review

## Security-Related Features Planned

As listed in our future enhancements:

- [ ] Multi-factor Authentication (2FA)
- [ ] Enhanced audit logging
- [ ] Advanced rate limiting
- [ ] IP-based access controls
- [ ] Enhanced monitoring and alerting

## Compliance

While the ChatLogger API itself doesn't claim specific compliance certifications, it's designed with security best practices that can help as part of a compliant environment:

- Password hashing and security practices align with OWASP recommendations
- Role-based access control can help with regulatory requirements
- Audit logging capabilities can assist with compliance monitoring
- API key management follows security best practices

## Security Contacts

- Primary Contact: [`security-chatlogger@kjanat.com`](mailto:security-chatlogger@kjanat.com)
- Secondary Contact: [`admin-chatlogger@kjanat.com`](mailto:admin-chatlogger@kjanat.com)

---

This security policy was last updated on May 11, 2025.
