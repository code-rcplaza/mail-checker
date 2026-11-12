# 005 · Núcleo de cuentas (fase 2, parte 4)

## Problema
La sincronización necesita saber qué cuentas de correo leer y con qué contraseña, pero
todavía no hay forma de registrarlas ni de guardar la contraseña sin exponerla. Esta
idea implementa la gestión de cuentas y el puente que le entrega las credenciales a
`emailsync`, con TDD y dobles en memoria, sin base ni Keychain.

## Alcance
- `accounts/domain`:
  - `Account` (`ID`, `PublicID`, correo, contraseña cifrada, creación, última
    sincronización, llave del último error).
  - `ParseEmail`: quita espacios, pasa a minúsculas y valida el formato
    (`algo@dominio`); cualquier dominio (sirve Gmail y Google Workspace).
  - `ParseAppPassword`: quita los espacios y exige 16 letras a-z.
  - Llaves: `accounts.email_invalid`, `accounts.app_password_invalid`,
    `accounts.email_taken`, `accounts.not_found`, `accounts.password_unreadable`.
- `accounts/ports`: `AccountRepository`, `SecretCipher` (dobles en memoria).
- `accounts/usecases`:
  - `AddAccount(email, appPassword)`: valida, cifra y guarda. `PublicID` lo genera un
    generador inyectado (como en `statements`). Correo ya registrado → `ErrConflict`
    con `accounts.email_taken` (lo detecta el repositorio con `UNIQUE`).
  - `ListAccounts` (sin contraseñas fuera de `domain`).
  - `RemoveAccount(publicID)`: `ErrNotFound` con `accounts.not_found`; el repositorio
    borra también sus cobros.
  - `Credentials`: cuentas con la contraseña descifrada. Si una no se puede descifrar
    (llave perdida) → error con `accounts.password_unreadable`; el usuario vuelve a
    ingresar la contraseña.
  - `RecordSync(accountID, at, errKey)`: guarda el resultado de la última
    sincronización.
- `emailsync/adapters/accounts`: puente que implementa `AccountDirectory` llamando a
  `Credentials` y `RecordSync`.

## Fuera
- Cifrado real AES-256-GCM y llave en Keychain (`SecretCipher` real) → fase 4.
- Repositorio SQLite y generador UUID v7 real → fase 4.
- Endpoints `/api/accounts` → fase 5.
- Editar una cuenta o cambiar su contraseña → YAGNI (se borra y se vuelve a agregar).

## Features afectadas
- `internal/accounts`
- `internal/emailsync` (solo el adaptador puente)

## Dudas abiertas
- Ninguna. La historia 05 depende de 003-02 y 003-03. Decidido: cualquier dominio de correo; la contraseña de aplicación se valida
  como 16 letras.
- Cambio aplicado a `ARCHITECTURE.md`: en `AccountRepository.Create` y
  `StatementRepository.Create` el comentario dice "asigna ID y PublicID"; el `PublicID`
  se genera en Go antes de guardar (generador inyectado), así que el repositorio solo
  asigna el `ID`.

## Historias
- [x] [01 Tipos de `accounts`](01-account-types.md)
- [ ] [02 `AddAccount` y `ListAccounts`](02-add-and-list-accounts.md)
- [ ] [03 `RemoveAccount` y `RecordSync`](03-remove-account-and-record-sync.md)
- [ ] [04 `Credentials`](04-credentials.md)
- [ ] [05 Puente `AccountDirectory`](05-account-directory-bridge.md)
