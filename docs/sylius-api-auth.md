# Sylius API Auth

## Promoting API access
User has to have the **ROLE_API_ACCESS** role in order to access /api resources. Use following cmd to promote a user:

```bash
php bin/console sylius:user:promote
```

## Creating client
```bash
php bin/console sylius:oauth-server:create-client \
    --grant-type="password" \
    --grant-type="refresh_token" \
    --grant-type="token"
```

## Getting token
```bash
curl http://localhost:8888/api/oauth/v2/token \
        -d "client_id"=jvvi9vwn3y84k88gw44c80sogk0ssgokww0o4so0k8wckowg0 \
        -d "client_secret"=uju2e8fgec080w8g4cc0s8sw0swgko8g084wgkgsckkkwwkss \
        -d "grant_type"=password \
        -d "username"=test@test.com \
        -d "password"=123123
```

