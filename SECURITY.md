# Security Policy

## Supported Versions

We provide security updates for the following versions:

| Version | Supported          |
| ------- | ------------------ |
| 1.x.x   | :white_check_mark: |
| < 1.0   | :x:                |

## Reporting a Vulnerability

We take security seriously. If you discover a security vulnerability, please follow these steps:

### 1. **DO NOT** create a public issue

Security vulnerabilities should be reported privately to protect users.

### 2. Email us directly

Send an email to: [security@tindertrip.com](mailto:security@tindertrip.com)

Include the following information:
- Description of the vulnerability
- Steps to reproduce
- Potential impact
- Suggested fix (if any)
- Your contact information

### 3. Response Timeline

- **Acknowledgment**: Within 48 hours
- **Initial Assessment**: Within 7 days
- **Fix Timeline**: Depends on severity (1-90 days)
- **Public Disclosure**: After fix is deployed

### 4. Vulnerability Severity

We use the following severity levels:

- **Critical**: Remote code execution, authentication bypass
- **High**: Privilege escalation, data exposure
- **Medium**: Information disclosure, denial of service
- **Low**: Minor security improvements

## Security Measures

### Authentication & Authorization

- JWT tokens with secure signing
- Password hashing with bcrypt
- OAuth2 integration (Google)
- Rate limiting on auth endpoints
- Session management

### Data Protection

- Input validation and sanitization
- SQL injection prevention
- XSS protection
- CSRF protection
- Secure headers

### Infrastructure

- HTTPS enforcement
- Secure database connections
- Environment variable protection
- Docker security best practices
- Regular dependency updates

### Monitoring

- API request logging
- Error tracking
- Security event monitoring
- Performance monitoring

## Security Best Practices

### For Developers

1. **Never commit secrets** to version control
2. **Use environment variables** for sensitive data
3. **Validate all inputs** from users
4. **Use prepared statements** for database queries
5. **Keep dependencies updated**
6. **Follow secure coding practices**

### For Users

1. **Use strong passwords**
2. **Enable 2FA** when available
3. **Keep your client updated**
4. **Report suspicious activity**
5. **Use HTTPS** connections

## Security Tools

We use the following security tools:

- **golangci-lint**: Code quality and security linting
- **Trivy**: Vulnerability scanning
- **gosec**: Go security scanner
- **CodeQL**: GitHub security analysis
- **Dependabot**: Dependency updates

## Security Updates

Security updates are released as:
- **Patch releases** for critical vulnerabilities
- **Minor releases** for high/medium vulnerabilities
- **Major releases** for significant security improvements

## Responsible Disclosure

We follow responsible disclosure practices:

1. **Report privately** first
2. **Allow time** for fix development
3. **Coordinate disclosure** with maintainers
4. **Credit researchers** appropriately
5. **Learn and improve** from each incident

## Security Contacts

- **Security Team**: [security@tindertrip.com](mailto:security@tindertrip.com)
- **Maintainer**: [callmetos@github.com](mailto:callmetos@github.com)
- **Emergency**: [emergency@tindertrip.com](mailto:emergency@tindertrip.com)

## Security Changelog

### 2024-09-30
- Initial security policy
- Added vulnerability reporting process
- Implemented security scanning tools

## Acknowledgments

We thank the security researchers and community members who help keep our project secure.

## License

This security policy is part of our project and is subject to the same license terms.
