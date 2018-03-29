# Sylius attributes API
## Create new "select" attribute
```bash
curl http://localhost:8888/api/v1/product-attributes/select \
    -H "Authorization: Bearer OWM0MjBlYTY0ZDUyOWE4OTg0ODgwY2IyOTRlYTZiYTE1ZjI5OGU5N2ExYmY2NDgyNzY1ZjU1NjcyMDJlZTExNA" \
    -H "Content-Type: application/json" \
    -X POST \
    --data '
        {
            "code": "highlite_brand",
            "translations": {  
                "ru_RU": {
                    "name": "Brand"
                }
            },
            "configuration": {
                "choices": [
                    {
                        "ru_RU": "Бренд 1",
                        "en_US": "Brand 1"
                    },
                    {
                        "ru_RU": "Бренд 2",
                        "en_US": "Brand 2"
                    }
                ]
            }
        }
    '
```
## Get single attribute by code
```bash
curl http://localhost:8888/api/v1/product-attributes/highlite_brand_test3 \
   -H "Authorization: Bearer OWM0MjBlYTY0ZDUyOWE4OTg0ODgwY2IyOTRlYTZiYTE1ZjI5OGU5N2ExYmY2NDgyNzY1ZjU1NjcyMDJlZTExNA" \
   -H "Accept: application/json"
```

