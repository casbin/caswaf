# Cookie Security Feature

## Overview

CasWAF now supports automatic addition of security flags to cookies set by backend servers through the reverse proxy. This feature helps protect against session hijacking and other cookie-related attacks.

## Configuration

The following configuration options are available per site in the Site configuration:

- `enableCookieSecure`: Add the `Secure` flag to cookies (requires HTTPS)
- `enableCookieHttpOnly`: Add the `HttpOnly` flag to cookies
- `enableCookieSameSite`: Add the `SameSite=Lax` flag to cookies

## Security Flags Explained

### Secure Flag
- **Purpose**: Ensures cookies are only sent over HTTPS connections
- **Protection**: Prevents cookie theft through man-in-the-middle (MITM) attacks on unencrypted HTTP connections
- **Important**: This feature should only be enabled when the reverse proxy is accessed exclusively via HTTPS

### HttpOnly Flag
- **Purpose**: Prevents JavaScript from accessing the cookie via `document.cookie`
- **Protection**: Mitigates cross-site scripting (XSS) attacks that attempt to steal cookies

### SameSite Flag
- **Purpose**: Controls whether cookies are sent with cross-site requests
- **Protection**: Helps prevent cross-site request forgery (CSRF) attacks
- **Value**: Set to `Lax` by default, which provides a balance between security and usability

## Example Configuration

When creating or updating a site via the API, include these fields:

```json
{
  "owner": "admin",
  "name": "my-site",
  "domain": "example.com",
  "host": "http://backend:8080",
  "enableCookieSecure": true,
  "enableCookieHttpOnly": true,
  "enableCookieSameSite": true
}
```

## Behavior

When these flags are enabled:

1. CasWAF intercepts responses from the backend server
2. For each `Set-Cookie` header, it checks if the security flags are already present
3. If a flag is missing and the corresponding feature is enabled, it adds the flag
4. The modified cookie is sent to the client

### Example

**Backend Response:**
```
Set-Cookie: PHPSESSID=abc123; Path=/
Set-Cookie: loginToken=xyz789; Path=/; Domain=example.com
```

**With Cookie Security Enabled:**
```
Set-Cookie: PHPSESSID=abc123; Path=/; Secure; HttpOnly; SameSite=Lax
Set-Cookie: loginToken=xyz789; Path=/; Domain=example.com; Secure; HttpOnly; SameSite=Lax
```

## Important Notes

1. **HTTPS Requirement**: The `Secure` flag should only be enabled when your site is accessed exclusively via HTTPS. If users access the site via HTTP, cookies with the Secure flag will not be sent to the server.

2. **Backward Compatibility**: These features are opt-in and disabled by default. Existing sites will continue to work without any changes.

3. **No Duplication**: If a backend already sets security flags on cookies, CasWAF will not add duplicate flags.

4. **Smart Detection**: CasWAF properly parses cookie attributes to avoid false positives. For example, a cookie value containing the word "secure" will not prevent the Secure flag from being added.

## Use Cases

This feature is particularly useful when:

- You have legacy backend applications that don't set proper cookie security flags
- You want to enforce consistent cookie security policies across multiple backend services
- You need to add security hardening without modifying backend application code
- You're running a multi-tenant platform where different backends have varying levels of security implementation

## Testing

The feature includes comprehensive unit and integration tests. Run tests with:

```bash
go test ./service -v
```
