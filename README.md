## Highlite2 import
This is a product import tool written in Golang for Highlite SPb project. It takes products from csv export
files and uploads it to [highlite2-sylius](https://github.com/oliosinter/highlite2-sylius) component through
[Sylius REST API](http://docs.sylius.com/en/latest/api/index.html).

## Requirements
In order to use highlite2-import you need to setup [highlite2-sylius](https://github.com/oliosinter/highlite2-sylius)
component first: it is supposed, that sylius has a default channel with EUR currency by default and two enabled
locales: en_US and ru_RU. 


## Sylius "how to" notes

### How to promote user API access?
User has to have the **ROLE_API_ACCESS** role in order to access /api resources. Use folloeing cmd to promote user:
`php bin/console sylius:user:promote`.

### How to create OAuth client?
Use following cmd:
```bash
php bin/console sylius:oauth-server:create-client \
    --grant-type="password" \
    --grant-type="refresh_token" \
    --grant-type="token"
```

