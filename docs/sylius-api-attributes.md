# Sylius attributes API
## Create new "select" attribute
```bash
curl http://localhost:8888/api/v1/product-attributes/select \
    -H "Authorization: Bearer YmUxYmFmMTcxNTVmNDRiMWZlMmMzNzY2ZTQyNGQ4MDU1MmJiYjc0ZmUxOTU1YWVjMWM2MDJhNGYyNjM5N2EwOQ" \
    -H "Content-Type: application/json" \
    -X POST \
    --data '
        {
            "code": "highlite_brand",
            "translations": {  
                "ru_RU": { "name": "Бренд" },
                "en_US": { "name": "Brand" }
            },
            "configuration": {
                "choices": [
                    {
                        "ru_RU": "Бренд 1",
                        "en_US": "Brand 1"
                    }
                ]
            }
        }
    '
```
## Update attribute
```bash
curl http://localhost:8888/api/v1/product-attributes/highlite_brand \
    -H "Authorization: Bearer YmUxYmFmMTcxNTVmNDRiMWZlMmMzNzY2ZTQyNGQ4MDU1MmJiYjc0ZmUxOTU1YWVjMWM2MDJhNGYyNjM5N2EwOQ" \
    -H "Content-Type: application/json" \
    -i \
    -X PATCH \
    --data '
        {
            "configuration": {
                "choices": {
                    "highlite_brand_option_1": {
                        "ru_RU": "Бренд 1",
                        "en_US": "Brand 1"
                    },
                    "highlite_brand_option_2": {
                        "ru_RU": "Бренд 2",
                        "en_US": "Brand 2"
                    },
                    "highlite_brand_option_3": {
                        "ru_RU": "Бренд 3",
                        "en_US": "Brand 3"
                    }
                }
            }
        }
    '
```

## Get single attribute by code
```bash
curl http://localhost:8888/api/v1/product-attributes/highlite_brand \
   -H "Authorization: Bearer YmUxYmFmMTcxNTVmNDRiMWZlMmMzNzY2ZTQyNGQ4MDU1MmJiYjc0ZmUxOTU1YWVjMWM2MDJhNGYyNjM5N2EwOQ" \
   -H "Accept: application/json" | python -m json.tool
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
            "highlite_brand_option_1": {
                "en_US": "Brand 1",
                "ru_RU": "\u0411\u0440\u0435\u043d\u0434 1"
            },
            "highlite_brand_option_2": {
                "en_US": "Brand 2",
                "ru_RU": "\u0411\u0440\u0435\u043d\u0434 2"
            },
            "highlite_brand_option_3": {
                "en_US": "Brand 3",
                "ru_RU": "\u0411\u0440\u0435\u043d\u0434 3"
            }
        },
        "multiple": false
    },
    "createdAt": "2018-03-30T10:02:48+00:00",
    "id": 2,
    "position": 0,
    "translations": {
        "en_US": {
            "id": 3,
            "locale": "en_US",
            "name": "Brand"
        },
        "ru_RU": {
            "id": 2,
            "locale": "ru_RU",
            "name": "\u0411\u0440\u0435\u043d\u0434"
        }
    },
    "type": "select",
    "updatedAt": "2018-03-30T10:04:13+00:00"
}
```

