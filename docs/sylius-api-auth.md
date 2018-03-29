# Sylius API Auth
## Getting token
```bash
curl http://localhost:8888/api/oauth/v2/token \
        -d "client_id"=5lhgthz7ezgg44gcw00swsccgc8oosc84sgok8g04osgokkko8 \
        -d "client_secret"=7pgnroliqog80wwwcswogkwwkgowogw84gsc0skog4gc0sg88 \
        -d "grant_type"=password \
        -d "username"=test@test.com \
        -d "password"=123123
```

