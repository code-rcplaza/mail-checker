# 003-02 · Tipos de `emailsync`
Feature: `internal/emailsync/domain`

Como desarrollador quiero los tipos de la sincronización y las conversiones desde un
correo para que el caso de uso solo orqueste pasos.

## Criterios
- Dado un `ParsedCharge` sin fecha y un `RawEmail` con Message-ID y fecha, cuando se
  llama `chargeFromEmail(parsed, email, accountID, bank)`, entonces el `Charge` tiene
  fuente `email`, external id = Message-ID, la cuenta, el banco y la fecha del correo.
- Dado un `ParsedCharge` con fecha, cuando se llama `chargeFromEmail`, entonces el
  `Charge` conserva la fecha del aviso.
- Dado un `RawEmail`, cuando se llama `unparsedFromEmail(email, accountID, bank)`,
  entonces el `UnparsedEmail` tiene external id = Message-ID, cuenta, banco, asunto y
  fecha del correo.
- Dado un cuerpo de más de 200 caracteres (con tildes o ñ), cuando se llama
  `unparsedFromEmail`, entonces la vista previa tiene 200 caracteres y no parte ninguno.
- Dado un cuerpo corto, cuando se llama `unparsedFromEmail`, entonces la vista previa
  es el cuerpo completo.
- Las llaves `emailsync.mail_auth_failed`, `emailsync.mail_read_failed`,
  `emailsync.sync_in_progress`, `emailsync.months_out_of_range` y
  `emailsync.limit_out_of_range` existen como constantes; `ErrMailAuthRejected` existe
  para que el lector de correo avise credenciales rechazadas.

## Fuera
- Uso de estas conversiones en `SyncAll` → 003-03 y 003-04.

## Notas
- Tipos: `RawEmail` (Message-ID, asunto, remitente, fecha, cuerpo), `MailQuery`
  (remitente, asunto, desde), `Credential`, `SyncResult`, `UnparsedEmail`.
- Largo de la vista previa como constante con nombre.
- Depende de 003-01.
