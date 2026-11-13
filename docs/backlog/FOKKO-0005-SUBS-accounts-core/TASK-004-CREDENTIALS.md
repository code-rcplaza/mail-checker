# FOKKO-0005-SUBS/TASK-004 · `Credentials`
Feature: `internal/accounts/usecases`

Como desarrollador quiero las cuentas con su contraseña descifrada para que la
sincronización pueda leer el correo.

## Criterios
- Dadas cuentas guardadas, cuando se llama `Credentials`, entonces devuelve una
  credencial por cuenta, en el orden del repositorio, con id, `PublicID`, correo y
  contraseña descifrada.
- Dada una cuenta cuya contraseña no se puede descifrar, cuando se llama `Credentials`,
  entonces `ErrConflict` con `accounts.password_unreadable` y no devuelve credenciales
  parciales.
- Dado un error del repositorio, cuando se llama `Credentials`, entonces se devuelve
  envuelto con contexto.

## Fuera
- Puente hacia `emailsync` → FOKKO-0005-SUBS/TASK-005.

## Notas
- La contraseña descifrada nunca se loguea ni va en el texto del error.
