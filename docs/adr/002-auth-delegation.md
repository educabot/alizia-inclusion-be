# ADR 002: Delegación de Autenticación a auth-service

**Date**: 2026-05-22
**Status**: accepted

## Context

El sistema original tenía autenticación local con `password_hash` almacenado en la tabla
`users`. Esto generaba duplicación de lógica de credenciales con otros servicios del
ecosistema Educabot y dificultaba implementar SSO entre plataformas.

## Decision

Delegamos autenticación completamente al auth-service externo de Educabot:

- El backend **no gestiona passwords** ni sesiones propias
- El middleware valida JWTs firmados con RS256 usando la clave pública del auth-service
- `org_uuid` y `user_id` se extraen del payload del token y se propagan via `context.Context`
- La migración `000011` hizo `password_hash` nullable para mantener compatibilidad de schema
- El endpoint `POST /auth/login` es un proxy al auth-service (no valida credenciales localmente)

Flujo:

```
Cliente → JWT en Authorization header → JWTMiddleware (valida RS256) → Handler (lee ctx)
```

## Consequences

**Positivas**:
- SSO unificado con otros servicios Educabot sin fricción para el usuario
- No se almacenan ni gestionan passwords localmente (menor superficie de ataque)
- Tokens RS256 verificables offline — no requiere llamada al auth-service por request
- Responsabilidad de seguridad de credenciales centralizada en un solo servicio

**Negativas**:
- Dependencia del auth-service para el flujo de login (si cae, los usuarios no pueden iniciar sesión)
- Se requiere sincronización de usuarios entre auth-service y la tabla local `users`
- Rotación de clave pública RS256 requiere actualización coordinada de config
