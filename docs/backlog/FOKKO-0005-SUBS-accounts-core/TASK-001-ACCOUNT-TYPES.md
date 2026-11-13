# FOKKO-0005-SUBS/TASK-001 · Tipos de `accounts`
Feature: `internal/accounts/domain`

Como usuario quiero que la app rechace un correo o una contraseña de aplicación mal
escritos al agregar la cuenta, y no recién al sincronizar.

## Criterios
- Dado `ParseEmail("  User@Example.com ")`, cuando se llama, entonces devuelve
  `user@example.com` sin error.
- Dado un correo de otro dominio (ej. `user@company.cl`), cuando se llama `ParseEmail`,
  entonces lo acepta.
- Dado `""`, `"user"`, `"user@"`, `"@example.com"`, `"a@b@example.com"` o
  `"us er@example.com"`, cuando se llama `ParseEmail`, entonces `ErrInvalid` con
  `accounts.email_invalid`.
- Dado `"abcd efgh ijkl mnop"`, cuando se llama `ParseAppPassword`, entonces devuelve
  `abcdefghijklmnop` sin error.
- Dado un texto con menos o más de 16 letras, o con dígitos o símbolos, cuando se llama
  `ParseAppPassword`, entonces `ErrInvalid` con `accounts.app_password_invalid`.
- Las llaves `accounts.email_invalid`, `accounts.app_password_invalid`,
  `accounts.email_taken`, `accounts.not_found` y `accounts.password_unreadable` existen
  como constantes.

## Fuera
- Uso en `AddAccount` → FOKKO-0005-SUBS/TASK-002.

## Notas
- Correos sintéticos (`example.com`); contraseñas de ejemplo inventadas.
