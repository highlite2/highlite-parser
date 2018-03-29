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
curl http://localhost:8888/api/v1/product-attributes/highlite_brand \
   -H "Authorization: Bearer OWM0MjBlYTY0ZDUyOWE4OTg0ODgwY2IyOTRlYTZiYTE1ZjI5OGU5N2ExYmY2NDgyNzY1ZjU1NjcyMDJlZTExNA" \
   -H "Accept: application/json"
```
```json
{
    "_links": {
        "self": {
            "href": "/api/v1/product-attributes/highlite_brand"
        }
    },
    "code": "highlite_brand",
    "configuration": {
        "choices": {
            "b5107c20-3344-11e8-ac1a-1113d1c6a0dd": {
                "en_US": "Brand 1",
                "ru_RU": "\u0411\u0440\u0435\u043d\u0434 1"
            },
            "b5108b8e-3344-11e8-8a3f-17afd4afb247": {
                "en_US": "Brand 2",
                "ru_RU": "\u0411\u0440\u0435\u043d\u0434 2"
            },
            "b51091ce-3344-11e8-866e-2329bf5be702": {
                "en_US": "Brand 3",
                "ru_RU": "\u0411\u0440\u0435\u043d\u0434 3"
            }
        },
        "multiple": false
    },
    "createdAt": "2018-03-29T11:31:22+00:00",
    "id": 6,
    "position": 0,
    "translations": {
        "ru_RU": {
            "id": 6,
            "locale": "ru_RU",
            "name": "Brand"
        }
    },
    "type": "select",
    "updatedAt": "2018-03-29T11:31:22+00:00"
}
```

