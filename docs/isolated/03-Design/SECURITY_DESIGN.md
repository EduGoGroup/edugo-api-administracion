# Seguridad - spec-03

## STRIDE

### Spoofing
**Mitigación:** JWT requerido en todos los endpoints

### Tampering
**Mitigación:** Trigger SQL previene ciclos en jerarquía

### Repudiation
**Mitigación:** Logs con user_id en todas las acciones

### Information Disclosure
**Mitigación:** Usuarios solo ven unidades de su escuela

### Denial of Service
**Mitigación:** Rate limiting 100 req/min por usuario

### Elevation of Privilege
**Mitigación:** Solo owners pueden modificar unidades

## Permisos Jerárquicos

| Acción | Requiere |
|--------|----------|
| Crear unidad hija | Owner de unidad padre |
| Modificar unidad | Owner de esa unidad |
| Agregar miembro | Owner de la unidad |
| Ver unidad | Miembro o owner |
