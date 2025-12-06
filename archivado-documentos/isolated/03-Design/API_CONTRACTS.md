# Contratos de API - spec-03

## Endpoints

### POST /v1/schools
```json
Request:
{
  "name": "Colegio San José",
  "code": "CSJ"
}

Response 201:
{
  "school_id": "uuid",
  "name": "Colegio San José",
  "code": "CSJ"
}
```

### POST /v1/units
```json
Request:
{
  "school_id": "uuid",
  "parent_unit_id": "uuid|null",
  "unit_type": "grade|section|club",
  "display_name": "5.º Año",
  "code": "5TO"
}

Response 201:
{
  "unit_id": "uuid",
  "display_name": "5.º Año"
}
```

### GET /v1/units/:id/tree
```json
Response 200:
{
  "id": "uuid",
  "display_name": "5.º Año",
  "children": [
    {
      "id": "uuid",
      "display_name": "Sección A",
      "children": []
    }
  ]
}
```

### POST /v1/units/:id/members
```json
Request:
{
  "user_id": "uuid",
  "role": "student|teacher|owner"
}

Response 201: {"membership_id": "uuid"}
```
