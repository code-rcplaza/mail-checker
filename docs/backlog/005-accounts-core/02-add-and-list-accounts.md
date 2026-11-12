# 005-02 · `AddAccount` y `ListAccounts`
Feature: `internal/accounts/usecases`

Como usuario quiero agregar mis cuentas de correo y verlas listadas.

## Criterios
- Dados un correo y una contraseña de aplicación válidos, cuando se llama `AddAccount`,
  entonces se cifra la contraseña ya normalizada y se llama `Create` una vez con el
  correo normalizado, la contraseña cifrada, el `PublicID` del generador y `CreatedAt`
  del reloj.
- Dada una cuenta agregada, entonces el texto plano de la contraseña no aparece en la
  `Account` guardada ni en la devuelta.
- Dado un correo o una contraseña inválidos, cuando se llama `AddAccount`, entonces el
  `ErrInvalid` de `ParseEmail`/`ParseAppPassword` y no se cifra ni se guarda nada.
- Dado un `Create` que devuelve `ErrConflict` con `accounts.email_taken`, cuando se
  agrega, entonces el error se devuelve conservando su tipo y llave.
- Dado un error del cifrador o del repositorio, cuando se agrega, entonces se devuelve
  envuelto con contexto.
- Dadas cuentas guardadas, cuando se llama `ListAccounts`, entonces devuelve las del
  repositorio; un error del repositorio se devuelve envuelto.

## Fuera
- Cifrado real y repositorio SQLite → fase 4.

## Notas
- Puertos `AccountRepository` y `SecretCipher` en `ports`; dobles en memoria en los
  tests (el cifrador doble produce bytes distintos del texto plano).
