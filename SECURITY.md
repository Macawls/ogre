# Security Policy

## Supported Versions

| Version | Supported |
|---------|-----------|
| 0.1.x   | Yes       |

## Reporting a Vulnerability

Report security vulnerabilities via GitHub Security Advisories on this repository.

Do not open a public issue for security vulnerabilities.

## Security Features

Ogre includes built-in protections when used as a server:

- Rate limiting per IP (configurable)
- Render timeout (default 10s)
- Maximum HTML element count (default 1000)
- Maximum request body size (10MB)
- Font size limits (5MB per font, 5 fonts per request)
