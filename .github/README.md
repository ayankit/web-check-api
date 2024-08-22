<h1 align="center">Web Check API</h1>
<p align="center">
  <i>A light-weight Go API for discovering website data</i><br />
</p>

## Usage

### Deploying using Docker

```
docker run -p 8080:8080 ayankit/web-check-api
```
### Example
Health Check
```
http://localhost:8080/health
```
Get IP of `google.com`
```
http://localhost:8080/api/get-ip?url=google.com
```

## Available Endpoints
```
/health - Health check

/api/get-ip
/api/firewall
/api/tech-stack
/api/headers
/api/http-security
/api/ports
/api/tls
/api/dns
/api/dns-server
/api/dnssec
/api/linked-pages
/api/block-lists
/api/carbon
/api/cookies
/api/hsts
/api/legacy-rank
/api/quality
/api/rank
/api/redirects
/api/social-tags
/api/trace-route
```
