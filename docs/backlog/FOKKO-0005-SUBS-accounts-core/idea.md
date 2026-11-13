# FOKKO-0005-SUBS · Núcleo de cuentas (fase 2, parte 4)

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
- Ninguna. La tarea TASK-005 depende de FOKKO-0003-SUBS/TASK-002 y FOKKO-0003-SUBS/TASK-003. Decidido: cualquier dominio de correo; la contraseña de aplicación se valida
  como 16 letras.
- Cambio aplicado a `ARCHITECTURE.md`: en `AccountRepository.Create` y
  `StatementRepository.Create` el comentario dice "asigna ID y PublicID"; el `PublicID`
  se genera en Go antes de guardar (generador inyectado), así que el repositorio solo
  asigna el `ID`.

## Tareas
- [x] [TASK-001 Tipos de `accounts`](TASK-001-ACCOUNT-TYPES.md)
- [x] [TASK-002 `AddAccount` y `ListAccounts`](TASK-002-ADD-AND-LIST-ACCOUNTS.md)
- [x] [TASK-003 `RemoveAccount` y `RecordSync`](TASK-003-REMOVE-ACCOUNT-AND-RECORD-SYNC.md)
- [x] [TASK-004 `Credentials`](TASK-004-CREDENTIALS.md)
- [x] [TASK-005 Puente `AccountDirectory`](TASK-005-ACCOUNT-DIRECTORY-BRIDGE.md)
