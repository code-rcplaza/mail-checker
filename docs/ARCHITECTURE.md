# Mail Checker — Arquitectura v1

Detecta suscripciones a partir de las cartolas de tarjeta (PDF; CSV en el post-MVP). En el post-MVP suma los avisos
de compra que llegan al correo (Gmail). Soporta **varios bancos**: BCI es el primero; agregar otro
(ej. Falabella) es escribir adaptadores nuevos, sin modificar los casos de uso.

Producto multiusuario: cada usuario ve y borra solo lo suyo. Qué entra en el MVP (pantallas,
endpoints, decisiones abiertas) está en `docs/mvp.md`; lo marcado **post-MVP** viene después.

## Contenido

- [Decisiones](#decisiones)
- [Backend](#backend)
  - [Features](#features) · [Estructura del backend](#estructura-del-backend) · [Reglas de dependencias](#reglas-de-dependencias) · [Multiusuario](#multiusuario)
  - Por feature: [`shared/domain`](#shareddomain) · [`auth`](#auth) · [`accounts`](#accounts-post-mvp) · [`emailsync`](#emailsync-post-mvp) · [`statements`](#statements) · [`subscriptions`](#subscriptions)
  - [Cifrado](#cifrado-post-mvp) · [Modelo de datos](#modelo-de-datos-sqlite) · [API REST](#api-rest) · [Agregar un banco](#agregar-un-banco-ej-falabella) · [Tests](#tests)
- [Convenciones Go](#convenciones-go)
- [Frontend](#frontend)
  - [Estructura del frontend](#estructura-del-frontend) · [Mobile first](#mobile-first) · [Desarrollo y build](#desarrollo-y-build)
- [Plan](#plan) · [Post-MVP](#post-mvp)
- [Pendiente](#pendiente)

## Decisiones

| Tema | Decisión |
|---|---|
| Backend | Go. **Features verticales**, cada una con arquitectura hexagonal: `adapters/ domain/ ports/ usecases/` |
| Bancos | Multi-banco desde el inicio: un parser por banco, registrados en `main.go` |
| Router HTTP | `go-chi/chi` |
| Base de datos | SQLite (`modernc.org/sqlite`, sin cgo) + **sqlc** para las consultas + **goose** para las migraciones |
| Fuente de datos | MVP: **cartolas** (PDF). Post-MVP: CSV y correo por IMAP + contraseña de aplicación de Google (solo Gmail) |
| Acceso | Código por correo, sin contraseña; todo dato pertenece a un usuario |
| Cifrado (post-MVP) | AES-256-GCM para las contraseñas de aplicación; dónde vive la llave depende del alojamiento (hoy: Keychain de macOS con `zalando/go-keyring`, respaldo `MAILCHECKER_KEY`) |
| Alojamiento | Abierto (`docs/mvp.md`, decisión 20) |
| Frontend | React + TypeScript + Vite + Tailwind v4, componentes propios, **mobile first** |
| Routing / datos / formularios | TanStack Router + TanStack Query + react-hook-form + zod |
| Fechas | Unix en **segundos** (UTC) en la base y en la API; `time.Time` dentro de Go; la UI muestra hora de `America/Santiago` |
| IDs | `id` autoincremental interno; `public_id` UUID v7 (`google/uuid`, generado en Go) para lo que se expone |
| Fuera del MVP | Lectura de correo, conciliación cartola-correo, OAuth de Gmail, login con Google y Apple, calendario, notificaciones |

---

## Backend

### Features

Cada feature responde a "¿quién eres y qué haces?":

| Feature | Responsabilidad |
|---|---|
| `auth` | Registro usuarios, les doy acceso por código y mantengo su sesión |
| `accounts` (post-MVP) | Administro las cuentas de correo y sus credenciales cifradas |
| `emailsync` (post-MVP) | Leo los avisos de compra de cada banco y los convierto en cargos |
| `statements` | Importo cartolas, detecto de qué banco son y las convierto en cargos |
| `subscriptions` | Detecto suscripciones a partir de los cargos, calculo el gasto mensual, el próximo cobro y las alzas, y guardo lo que el usuario decide |

`accounts` y `emailsync` tienen su núcleo escrito (ideas 003 y 005) pero no entran al MVP.
| `shared` | Lo mínimo común: el cargo, el dinero, errores y la plataforma (SQLite, HTTP, reloj) |

### Estructura del backend

```
backend/
├── cmd/api/
│   ├── main.go                      señales del SO; todo lo demás en run
│   ├── run.go                       config → base → migraciones → rutas → limpieza → servidor
│   ├── wiring.go                    composition root: crea adaptadores y casos de uso, registra bancos, arma rutas
│   ├── cleanup.go                   limpieza periódica de accesos vencidos (cada hora)
│   └── config.go                    variables de entorno, fail fast
└── internal/
    ├── auth/                            registro, código por correo, sesión
    │   ├── adapters/                    http/, sqlite/, mailer/ (envío del código), export/ (puente: copia de los datos)
    │   ├── domain/                      User
    │   ├── ports/                       UserRepository, LoginCodes, CodeSends, Sessions, CodeSender
    │   └── usecases/                    code_delivery, register_user, login, verify_code, open_session,
    │                                    validate_session, close_session, get_me, delete_user
    │
    ├── accounts/                        (post-MVP)
    │   ├── adapters/
    │   │   ├── http/                GET/POST/DELETE /api/accounts
    │   │   ├── sqlite/              AccountRepository
    │   │   └── keychain/            SecretCipher (AES-GCM + llave en Keychain)
    │   ├── domain/                  Account
    │   ├── ports/                   AccountRepository, SecretCipher
    │   └── usecases/                add_account, list_accounts, remove_account, credentials, record_sync
    │
    ├── emailsync/                       (post-MVP)
    │   ├── adapters/
    │   │   ├── http/                POST /api/sync, GET /api/unparsed
    │   │   ├── imap/                MailReader (Gmail)
    │   │   ├── bci/                 BankEmailParser de BCI
    │   │   ├── sqlite/              ChargeWriter, UnparsedRepository
    │   │   └── accounts/            puente: AccountDirectory → casos de uso de accounts
    │   ├── domain/                  SyncResult, UnparsedEmail
    │   ├── ports/                   AccountDirectory, MailReader, BankEmailParser, ChargeWriter, UnparsedRepository
    │   └── usecases/                sync_all, list_unparsed
    │
    ├── statements/
    │   ├── adapters/
    │   │   ├── http/                POST /api/statements, GET /api/banks (GET y DELETE de cartolas: post-MVP)
    │   │   ├── extract/             FileExtractor: pdf (csv: post-MVP)
    │   │   ├── bci/                 BankStatementParser de BCI
    │   │   └── sqlite/              StatementRepository
    │   ├── domain/                  Statement, StatementLine
    │   ├── ports/                   StatementRepository, FileExtractor, BankStatementParser
    │   └── usecases/                import_statement, list_statements, remove_statement, list_banks
    │
    ├── subscriptions/
    │   ├── adapters/
    │   │   ├── http/                summary, decisions, aliases, overrides, manual, price-alerts, cancellations (reconciliation: post-MVP)
    │   │   ├── catalog/             ServiceCatalog (descriptor de cobro → servicio, enlace para cancelar)
    │   │   ├── exchangerates/       ExchangeRates (tipo de cambio; fuente abierta)
    │   │   └── sqlite/              ChargeReader, SubscriptionDecisions, LineAliases, ManualSubscriptions, LineOverrides, PriceAlertDecisions, HiddenLines, Cancellations, ExchangeRates
    │   ├── domain/                  Subscription, Summary, Reconciliation, reglas de recurrencia y de coincidencia
    │   ├── ports/                   ChargeReader, ServiceCatalog, SubscriptionDecisions, LineAliases, ManualSubscriptions, LineOverrides, PriceAlertDecisions, Cancellations, ExchangeRates
    │   └── usecases/                get_summary, decide_possible_subscription, undo_decision, name_line, manuales, overrides, alzas, cancelaciones (reconcile_charges: post-MVP)
    │
    └── shared/
        ├── domain/                  Money, Currency, Bank, Card, Source, Charge, ParsedCharge, errores
        ├── sqlite/                  conexión + migraciones goose (todas las tablas)
        ├── httpserver/              router chi, JSON, mapeo de errores a códigos HTTP
        └── clock/                   Clock real
```

### Reglas de dependencias

1. Dentro de una feature: `adapters → usecases → ports → domain`. **`domain` y `usecases` no importan
   ninguna librería de infraestructura** (chi, sqlite, imap, pdf…) ni nada de `adapters`.
2. **Una feature no importa a otra**, salvo desde un adaptador puente. Ejemplo: `emailsync` necesita
   las credenciales, así que declara el puerto `AccountDirectory`, y `emailsync/adapters/accounts`
   lo implementa llamando a `accounts/usecases`.
3. Todas pueden importar `shared/domain`. Solo los adaptadores y `main.go` importan `shared/sqlite`, `shared/httpserver` y `shared/clock`.
4. `main.go` es el único lugar que conoce todas las piezas concretas.

### Objetos en vez de listas de parámetros

Máximo 2 parámetros por función (sin contar `ctx`); con más, un objeto con el nombre de lo que es:
- **Entrada** de un caso de uso: `XInput` (`DecisionInput`, `LineOverrideInput`,
  `ManualSubscriptionInput`).
- **Dependencias** de un constructor: `XDependencies` (`ImportStatementDependencies`,
  `SyncAllDependencies`) o, si se leen como fuentes, `SummarySources`.
- **Valores con clave** en los puertos: `LineDecision`, `LineAlias`, `SyncRecord`, `ImportedStatement`.
- **Contexto de un cálculo** en `domain`: `MerchantCatalog{Identify, Normalize}`,
  `SubscriptionDetector{LooksLikeSubscription, Now}` (con métodos), `SummaryInput`. El `Summary`
  guarda `AsOf`, así los pasos que lo transforman no vuelven a recibir `now`.
- **HTTP**: `httpserver.Exchange{Writer, Request}`, `Reply{Status, Body}` y
  `Responder{Logger}` con `Respond` y `WriteError`; todos los handlers responden por ahí.

### Multiusuario

- Todo dato (cargos, cartolas, decisiones, alias, manuales, overrides, decisiones de alza,
  cancelaciones) pertenece a un usuario y lleva `user_id`.
- Todo puerto de lectura y escritura filtra por el usuario; los `UNIQUE` pasan a ser por usuario.
  Un usuario nunca ve ni toca datos de otro.
- El `userID` sale de la sesión: el middleware `RequireSession` (auth) lo deja en el contexto de la
  petición (`httpserver.WithUserID`), el handler lo lee con `httpserver.UserIDFrom` (sin sesión es
  401, nunca un usuario cero) y desde ahí viaja como parámetro (`userID int64`) por casos de uso y
  puertos; el contexto solo lo lleva dentro de la capa HTTP.
- Eliminar al usuario borra todo lo suyo en cascada.

### `shared/domain`

```go
type Money int64               // centésimas: CLP 5.990 => 599000, USD 20,00 => 2000
type Currency string           // "CLP", "USD"
type Bank string               // "bci", "falabella"…
type Card string               // "credito", "debito"
type Source string             // "statement" (MVP), "email" (post-MVP)

// ParsedCharge: lo que un parser de banco extrae de un aviso o línea de cartola.
type ParsedCharge struct {
    Amount   Money
    Currency Currency
    Merchant string
    Card     Card
    CardLast4 string   // últimos 4 dígitos; vacío si el aviso o la cartola no los trae
    Date     time.Time // cero si la línea de cartola no trae fecha
}

// Charge: un cargo guardado.
type Charge struct {
    ID          int64
    UserID      int64    // dueño del cargo
    Bank        Bank
    Source      Source
    ExternalID  string  // correo: Message-ID; cartola: "<statementPublicID>:<línea>"
    AccountID   *int64
    StatementID *int64
    OccurredAt  time.Time
    Amount      Money
    Currency    Currency
    Merchant    string
    Card        Card
    CardLast4   string // vacío o exactamente 4 dígitos (common.card_last_digits_invalid)
}

var ErrInvalid, ErrNotFound, ErrConflict, ErrTooMany, ErrUnauthorized, ErrForbidden, ErrUnavailable error // HTTP 400, 404, 409, 429, 401, 403, 503
// Invalid(key, params), NotFound(key), Conflict(key), TooMany(key, params), Unauthorized(key), Forbidden(key):
// envuelven los anteriores y llevan una llave i18n (nunca texto para el usuario). Llaves comunes:
// common.invalid_json, common.not_found, common.conflict, common.route_not_found,
// common.validation_failed, common.internal, common.unauthorized, common.too_many_requests, common.forbidden.
// InvalidFields(fields) es common.validation_failed con la llave de cada campo (`fields`).
```

### `auth`

- **Registro** `{name, email, acceptedTerms}` y **login** `{email}` envían un código al correo;
  **verificar el código** abre la sesión; **cerrar sesión** la termina. No hay contraseñas: el único
  secreto es el código.
- Fail fast en el borde: nombre no vacío y de hasta 60 caracteres, correo con formato válido
  (se guarda en minúsculas y sin espacios) y términos aceptados; si no, `ErrInvalid` con una llave
  `auth.*`.
- **Respuesta neutra**: la respuesta no revela qué correos tienen cuenta. Registrar un correo que ya
  existe no crea otro usuario ni falla: se envía un código a su dueño. Entrar con un correo
  desconocido responde igual que con uno registrado y no envía nada. Verificar con un correo sin
  código responde `auth.code_invalid`.
- Código de 6 dígitos, vigencia de 10 minutos (a los 10 minutos exactos ya venció), 5 intentos por
  código, 3 peticiones por hora por correo y sesión de 30 días (`docs/mvp.md`, decisión 11). El límite
  cuenta las peticiones por correo exista o no la cuenta (si no, un 429 revelaría el correo).
- Un código nuevo reemplaza al anterior. El quinto intento fallido invalida el código. Un código
  correcto se consume.
- Ni el código ni el token de sesión se guardan en claro: solo su hash SHA-256 (`HashSecret`), y el
  código se compara en tiempo constante. Con 6 dígitos el hash no impide un ataque offline: lo que
  protege es el límite de intentos. El token son 32 bytes de `crypto/rand` en base64 URL-safe.
- La fuente aleatoria (`io.Reader`) y el reloj se inyectan.
- Envío: `mailer.ResendSender` hace `POST /emails` a la API de Resend con `net/http`, sin SDK y con
  timeout. Un rechazo se reporta solo con el status, nunca con el cuerpo, el correo ni el código.
  Asunto "Tu código de acceso a Fokko", texto simple en español. Configuración: `MAILCHECKER_ENV`
  (`development` por defecto, o `production`), `MAILCHECKER_RESEND_API_KEY` y `MAILCHECKER_MAIL_FROM`.
  En producción faltar la llave o el remitente detiene el arranque; en desarrollo, sin llave, se usa
  `mailer.ConsoleSender`. La llave es un `Secret` que nunca se imprime ni se loguea.

```go
type User struct {
    ID              int64
    PublicID        string    // UUID v7, generado en el caso de uso
    Name            string
    Email           string    // normalizado
    AcceptedTermsAt time.Time
}

type LoginCode struct { UserID int64; CodeHash string; ExpiresAt time.Time; Attempts int }
type Session   struct { UserID int64; TokenHash string; ExpiresAt time.Time }
type CodeSend  struct { Email string; SentAt time.Time }

type UserRepository interface {
    Create(ctx, user User) (User, error)          // asigna el ID
    FindByEmail(ctx, email string) (User, error)  // ErrNotFound, sin distinguir mayúsculas
    FindByID(ctx, id int64) (User, error)         // ErrNotFound
    Delete(ctx, id int64) error                   // ErrNotFound; una transacción: sus envíos (por correo) y el usuario, el resto por cascada
}
type LoginCodes interface { // un código por usuario
    Save(ctx, code LoginCode) error                       // reemplaza el anterior
    Find(ctx, userID int64) (LoginCode, error)            // ErrNotFound
    AddFailedAttempt(ctx, userID int64) (int, error)      // atómico; devuelve los intentos
    Delete(ctx, userID int64) error                       // idempotente
}
type CodeSends interface {
    ListSince(ctx, email string, since time.Time) ([]time.Time, error) // envíos después de since
    Record(ctx, send CodeSend) error
}
type Sessions interface {
    Create(ctx, session Session) error
    FindByTokenHash(ctx, tokenHash string) (Session, error) // ErrNotFound
    Delete(ctx, tokenHash string) error                     // idempotente
}
type CodeSender interface {
    Send(ctx, email, code string) error // consola en desarrollo, Resend en producción
}
type ExpiredAccess interface { // borra en una transacción
    DeleteExpired(ctx, cutoff CleanupCutoff) error // CleanupCutoffAt(now): vencidos a now; envíos de hace 1 h o más
}
```

- Casos de uso: `CodeDelivery` (componente: admite la petición según el límite y la registra; emite
  el código y lo envía), `RegisterUser`, `Login`, `VerifyCode` (devuelve el usuario),
  `OpenSession` (verifica y abre la sesión; devuelve el token en claro una sola vez),
  `ValidateSession` (token → `userID`; vencida o desconocida → `ErrUnauthorized` con
  `common.unauthorized`, sin distinguir), `CloseSession` (idempotente), `CleanExpiredAccess`
  (borra códigos, sesiones y envíos que ya no sirven; `cmd/api` la corre cada hora y un fallo solo
  se loguea), `DeleteUser` (borra la cuenta y todo lo suyo; sus sesiones caen con ella).
- `DELETE /api/me` no necesita puente: toda tabla con `user_id` tiene `ON DELETE CASCADE`, y
  `code_sends` (por correo) se borra en la misma transacción.
- `GET /api/export` lo arma el puente `auth/adapters/export` (`UserCopies`), el único lugar de auth
  que importa casos de uso de otras features: `GetMe`, `ListStatements`, `ListSavedChoices`
  (subscriptions: cargos, decisiones, alias, manuales y overrides tal como se guardaron) y
  `GetSummary` (12 meses, para el CSV).

### `accounts` (post-MVP)

> Fuera del MVP. El núcleo existe (idea 005) y queda sin uso hasta la lectura de correo por IMAP.

```go
// domain
type Account struct {
    ID          int64      // interno, nunca sale por la API
    PublicID    string     // UUID v7, el que se expone
    Email       string
    PasswordEnc []byte
    CreatedAt   time.Time
    LastSyncAt  *time.Time
    LastErrorKey shared.MessageKey // llave i18n de la última sincronización fallida; vacía si no hubo error
}

// ports
type AccountRepository interface {
    Create(ctx context.Context, a *Account) error             // asigna ID (PublicID viene generado en Go); ErrConflict si el correo ya existe
    List(ctx context.Context) ([]Account, error)
    Delete(ctx context.Context, publicID string) error        // ErrNotFound; borra también sus cargos
    MarkSynced(ctx context.Context, r SyncRecord) error       // SyncRecord{AccountID, At, ErrorKey}; ErrNotFound
}
type SecretCipher interface {
    Encrypt(plain string) ([]byte, error)
    Decrypt(sealed []byte) (string, error)
}

// usecases: AddAccount, ListAccounts, RemoveAccount,
//           Credentials (cuentas con su contraseña descifrada, para emailsync),
//           RecordSync (guarda resultado de la última sincronización)
```

Las cuentas **ya no guardan filtros** (`from`/`subject`): cada banco aporta el suyo.

### `emailsync` (post-MVP)

> Fuera del MVP. El núcleo existe (idea 003) y queda sin uso hasta la lectura de correo por IMAP.

```go
// domain
type SyncResult struct {
    AccountPublicID string
    Email     string
    Found, New, Unparsed int
    ErrorKey  shared.MessageKey // llave i18n, ej. emailsync.mail_auth_failed; vacía si no hubo error
}
type UnparsedEmail struct {
    ExternalID string
    AccountID  int64
    Bank       Bank
    Subject    string
    ReceivedAt time.Time
    Preview    string
}

// ports
type AccountDirectory interface {
    Credentials(ctx context.Context) ([]Credential, error)    // {AccountID, AccountPublicID, Email, Password}
    RecordSync(ctx context.Context, r SyncRecord) error       // SyncRecord{AccountID, At, ErrorKey}, propio de emailsync
}
type MailReader interface {
    Read(ctx context.Context, cred Credential, q MailQuery) ([]RawEmail, error)
}
type BankEmailParser interface {
    Bank() Bank
    Query() MailQuery                         // ej. BCI: from "bci.cl", asunto "uso de tu tarjeta"
    Parse(e RawEmail) (ParsedCharge, bool)
}
type ChargeWriter interface {
    Insert(ctx context.Context, c Charge) (inserted bool, err error) // idempotente por Source+ExternalID
}
type UnparsedRepository interface {
    Save(ctx context.Context, u UnparsedEmail) error
    Forget(ctx context.Context, externalID string) error
    List(ctx context.Context, limit int) ([]UnparsedEmail, error)
}
```

**`SyncAll(months)`:** para cada cuenta y cada banco registrado, busca con el filtro del banco, parsea
con su parser y guarda los cargos. El mismo aviso recibido en dos cuentas se guarda una sola vez
(Message-ID). Un error en una cuenta queda en su `SyncResult` y no detiene a las demás. Hay una sola
sincronización a la vez.

### `statements`

```go
// domain
type Statement struct {
    ID          int64      // interno
    PublicID    string     // UUID v7, el que se expone
    Bank        Bank
    Filename    string
    Kind        FileKind   // pdf (csv: post-MVP)
    BillingDate time.Time
    ImportedAt  time.Time
    Charges     int
    BadLines    []int      // números de línea; el front arma el texto
}

// ports
type FileExtractor interface {
    Extract(filename string, data []byte) (text string, kind FileKind, err error)
}
type BankStatementParser interface {
    Bank() Bank
    Detect(text string) bool                                // ¿esta cartola es de mi banco?
    Parse(text string) (lines []StatementLine, bad []int)    // StatementLine{Number, Charge}; bad: números de línea, nunca el texto
    BillingDate(text string) (time.Time, bool)               // fecha de facturación que trae el archivo
    TransferPair(text string) (TransferPair, bool)           // deuda traspasada en pesos y dólares (tipo de cambio)
}
type ExchangeRateRecorder interface {
    Record(ctx context.Context, userID int64, rate ExchangeRate) error // uno por usuario y fecha
}
type StatementRepository interface {
    Create(ctx context.Context, userID int64, imported ImportedStatement) error // ImportedStatement{Statement *Statement, Charges}; atómico; asigna ID
    List(ctx context.Context, userID int64) ([]Statement, error)
    Delete(ctx context.Context, userID int64, publicID string) error            // ErrNotFound; borra también sus cargos
}
```

**`ImportStatement(file, billingDate?, bank?)`:**
1. Extrae el texto del archivo.
2. Banco: si viene elegido, usa ese. Si no, pregunta a cada parser `Detect`. Si exactamente uno
   lo reconoce, usa ese; si no, responde la llave `statements.bank_not_detected`.
3. Parsea. Las líneas sin fecha toman la fecha de facturación. Cada línea trae su número en el texto,
   que se usa en el external id (`"<publicID>:<n>"`) y en las líneas malas.
4. Guarda la cartola y sus cargos en una transacción.

Reglas:
- Una cartola por usuario + banco + fecha de facturación. Si ya existe, responde `ErrConflict` con
  `statements.already_imported`; para reemplazarla se borra y se vuelve a subir.
- Si ninguna línea queda como cargo, responde `ErrInvalid` con `statements.no_charges`
  (params `badLines`) y no guarda nada.
- El `PublicID` se genera en Go antes de guardar (generador inyectado), porque el external id
  de cada cargo lo usa.
- La fecha de facturación sale del propio archivo (`BillingDate`). Si el parser no la encuentra, se
  usa la `billingDate` que el usuario escribió en la pantalla de subida; si tampoco viene,
  `ErrInvalid` `statements.billing_date_required`.
- **Qué se guarda:** por línea, el descriptor del comercio, el monto, la fecha, la tarjeta y sus
  últimos 4 dígitos. El texto de la cartola se descarta al terminar de parsear; `BadLines` guarda
  números, nunca texto.
- La subida responde archivo y meses leídos. Lo que Análisis muestra como "seguras" y "por
  confirmar" lo lee de `GET /api/summary`, para que `statements` no importe `subscriptions`.

**`RemoveStatement(publicID)`:** (se expone en el post-MVP, pantalla Fuentes) `Delete` devuelve `shared.ErrNotFound` sin llave y el caso de
uso agrega `statements.not_found`; los adaptadores no conocen llaves i18n.

### `subscriptions`

```go
// ports
type ChargeReader interface {
    ListSince(ctx context.Context, userID int64, since time.Time) ([]Charge, error)
}
type ServiceCatalog interface {
    Identify(merchant string) (Service, bool)   // "ANYSPHERE INC" -> Cursor
    Normalize(merchant string) string           // "C. Verde Mall" -> "C VERDE MALL"; "DLOCAL*SPOTIFY" -> "SPOTIFY"
    LooksLikeSubscription(merchant string) bool // "…SUBSCRIPTION", "MEMBERSHIP"
}
// Lo que el usuario decide sobre una línea se guarda como un valor con su clave.
type SubscriptionDecisions interface {
    List(ctx context.Context, userID int64) (map[GroupKey]Decision, error) // Decision: confirmed (+ cadencia) | dismissed
    Save(ctx context.Context, userID int64, d LineDecision) error          // LineDecision{Key, Decision}
    Delete(ctx context.Context, userID int64, key GroupKey) error
}
type LineAliases interface {
    List(ctx context.Context, userID int64) (map[GroupKey]string, error)
    Save(ctx context.Context, userID int64, a LineAlias) error             // LineAlias{Key, Alias}
    Delete(ctx context.Context, userID int64, key GroupKey) error
}
type LineOverrides interface {
    List(ctx context.Context, userID int64) (map[GroupKey]LineOverride, error)
    Save(ctx context.Context, userID int64, o OverriddenLine) error        // OverriddenLine{Key, Override}
    Delete(ctx context.Context, userID int64, key GroupKey) error
}
type PriceAlertDecisions interface {
    List(ctx context.Context, userID int64) (map[GroupKey][]Money, error)  // historial de montos mantenidos, del más antiguo al último
    Save(ctx context.Context, userID int64, k KeptRise) error              // KeptRise{Key, NewAmount}; agrega al historial, sin repetir
}
type HiddenLines interface {
    List(ctx context.Context, userID int64) ([]GroupKey, error)
    Hide(ctx context.Context, userID int64, h HiddenLine) error            // HiddenLine{Key, HiddenAt}; idempotente
}
type Cancellations interface {
    List(ctx context.Context, userID int64) (map[GroupKey]time.Time, error)
    Mark(ctx context.Context, userID int64, m CancellationMark) error      // CancellationMark{Key, MarkedAt}; conserva la primera fecha
}
type ManualSubscriptions interface {
    List(ctx context.Context, userID int64) ([]ManualSubscription, error)
    Create(ctx context.Context, userID int64, m ManualSubscription) error
    Update(ctx context.Context, userID int64, m ManualSubscription) error  // ErrNotFound
    Delete(ctx context.Context, userID int64, publicID string) error       // ErrNotFound
}
type ExchangeRates interface {
    List(ctx context.Context, userID int64) ([]ExchangeRate, error)
}
```

Todos los puertos filtran por usuario. Los casos de uso que leen o verifican líneas componen
`ChargeLines{Charges, Catalog, Now}` (lee, deduplica y agrupa), así la línea que se ve en el
resumen es la misma que se decide, nombra, sobrescribe, oculta o marca.

**`GetSummary(months)`:**
1. Lee los cargos del período.
2. Quita los cargos de cartola que ya llegaron por correo (regla de coincidencia, ver abajo).
   **Post-MVP:** en el MVP no hay correo y este paso no hace nada.
3. Agrupa en **padres** y, dentro de cada padre, en **líneas** (una línea = una suscripción; ver
   "Líneas" abajo):
   - padre: servicio conocido (sin importar moneda ni banco), servicio *umbrella* (Apple, Google
     Play…) o comercio desconocido por nombre normalizado y moneda;
   - los cobros sueltos de un padre (no encadenan con ninguna línea de su origen) van a su grupo
     `<padre>|one-off`, que termina en **otros cobros**.
4. Aplica las reglas de recurrencia del dominio (ver abajo) a cada línea y las decisiones del
   usuario (una decisión solo se aplica mientras la línea es posible). Cada línea lleva su padre y
   el alias que el usuario le dio (`LineAliases`).
5. Arma cuatro secciones:
   - **Suscripciones**: activas + posibles confirmadas (suman al gasto);
   - **Posibles**: servicio conocido con 1 solo cobro, o comercio desconocido que parece
     suscripción; sin decidir (no suman). Solo si siguen vigentes (ventana de **Activa**);
   - **Inactivas**: suscripciones y posibles cuyo último cobro quedó fuera de la ventana;
   - **Otros cobros**: lo que no es suscripción + posibles descartadas.
6. Suma el gasto mensual de **Suscripciones** por moneda (una anual cuenta como monto / 12) y el
   **total en CLP**: los cobros en USD se convierten (ver Pendiente, tipo de cambio).
7. Calcula el **próximo cobro** de cada línea: su último cobro más un intervalo de su cadencia (mes
   o año), en su día típico.
8. Cuenta los **cobros del mes**: cobrados (con un cobro en el mes calendario actual) y restantes
   (con próximo cobro aún dentro del mes).
9. Detecta las **alzas** y descarta las que el usuario ya decidió mantener.

Una **alza** es un cobro mayor que el anterior en la misma línea y moneda; cualquier aumento
cuenta. Trae el monto de antes y de después, la fecha del cobro con el monto nuevo ("desde") y el
impacto: anual = diferencia × 12 en una línea mensual (la diferencia misma si es anual) y mensual =
anual / 12. Si el salto supera 25 % se muestra además como cambio de plan. Sin correos de aviso, el
alza se detecta al llegar el cobro, es decir, con la cartola siguiente.

Las **manuales** y los **overrides** entran antes del cálculo: una manual es una línea sin cobros
que suma al gasto con su monto y su día de cobro; un override reemplaza el monto o el día
detectado. Una manual tiene como próximo cobro la siguiente vez que ocurre su día de cobro; un
override cambia lo que se muestra, el total y el próximo cobro, pero no los cobros reales.
"Dejar de seguir" oculta la línea (no suma ni aparece) y la marca se guarda por clave de línea
(`docs/mvp.md`, decisión 10).

En el MVP el handler HTTP devuelve **Suscripciones**, **Posibles** y las alzas; **Inactivas** y
**Otros cobros** se calculan (la detección los necesita) pero no se exponen.

**`DecidePossibleSubscription(groupKey, confirm(monthly | yearly) | dismiss)`:** guarda la decisión del
usuario sobre una posible. Al confirmar, **el usuario elige la cadencia** (ej. Netflix USD 120 → anual):
la app la respeta y no la vuelve a inferir. La decisión se guarda por **clave de línea**, estable
entre recálculos: `<padre>|<últimos 4 dígitos>|d<día típico>` (ej. `service:Netflix|4821|d01`;
sin dígitos: `umbrella:App Store||d09`). El padre es `service:<nombre>`, `umbrella:<nombre>` o
`merchant:<nombre normalizado>:<moneda>`. El día típico es la mediana de los días del mes de la
línea (el menor de los dos del medio si son pares). Si dos líneas chocan, la segunda (por su primer
cobro) lleva `#2`, la tercera `#3`. Sin ids ni correos en la clave (sale por la API). Si la línea
cambia de día de verdad, la clave cambia y la decisión se pierde (aceptado). Si la
llave no es una posible actual → `ErrNotFound`. La posible se busca en la ventana máxima (24 meses),
con el mismo paso de lectura y agrupación de `GetSummary`; se puede volver a decidir.

**`UndoDecision(groupKey)`:** borra la decisión guardada para esa clave de línea; la línea vuelve a
**Posibles** en el siguiente `GetSummary`. Es idempotente: si no había decisión, no falla.

**`NameLine(lineKey, alias)`:** el usuario nombra una línea (Personal, Padres…) para distinguir
suscripciones paralelas de un servicio. Se recorta; vacío borra el alias; más de 40 caracteres
(runas) → `ErrInvalid` `subscriptions.alias_too_long`. La línea (activa, posible o inactiva, nunca
el grupo `|one-off`) se busca en la ventana máxima; si no existe → `ErrNotFound`
`subscriptions.line_not_found`.

**Otras acciones del usuario (MVP)**, cada una con su caso de uso y su validación en el borde:
- **Manual** `{name, amount, billingDay}`: crear, editar y quitar una suscripción agregada a mano.
- **Override** `{lineKey, amount?, billingDay?}`: monto y día que reemplazan lo detectado.
- **Decisión de alza** `{lineKey}`: solo `keep` (cierra la tarjeta del alza). `later` (posponer o recordar) es post-MVP: sin correo no hay aviso previo.
- **Por cancelar** `{lineKey}`: marca la línea; no cancela nada ni lee correos. El enlace al sitio
  del servicio sale del catálogo.
- **Dejar de seguir** una línea: la oculta y guarda la marca por clave (`docs/mvp.md`, decisión 10).

**Acciones pendientes** del Resumen: primero las alzas no mantenidas, de mayor a menor impacto
anual **en pesos** (las de USD se convierten con el mismo tipo de cambio del total, que el
`Summary` guarda en `ExchangeRate`; sin tipo de cambio, las alzas en USD van al final), y después
las posibles sin decidir.

**Líneas** (constantes con nombre en `subscriptions/domain`):
1. **Origen**: dentro del padre, los cobros se separan por cuenta de correo y por últimos 4 dígitos (en el MVP no hay cuentas de correo: separan solo los dígitos)
   (un cobro sin dígitos no separa). Un cobro de cartola sin cuenta se une al origen de cuenta con
   el que encadena; si no encadena con ninguno, va al origen de cartola de su tarjeta.
2. **Banda de monto**: dentro del origen, los cobros se separan en bandas de ±25 % (misma moneda).
3. **Cadena**: en cada banda, un cobro continúa una línea si está a un intervalo válido de su
   último cobro (mensual k × 24–38 o anual 330–400 días); con varias, la de monto más parecido y,
   si empatan, la de fecha esperada más cercana. Si ninguna calza, abre una línea nueva.
4. **Cambio de plan**: en el mismo origen, una línea que empieza justo cuando otra termina (salto de
   monto > 25 %) se une a ella. En otro origen es otra suscripción.
5. **Cobro suelto**: una línea de un solo cobro cuyo origen volvió a cobrar más de 38 días después
   no siguió: va a `<padre>|one-off`. Hasta entonces se ve como posible.

**Reglas de recurrencia** (constantes con nombre en `subscriptions/domain`):
- **Intervalo**: mensual si los cobros están a k × 24–38 días con k ≤ 3 (tolera hasta 2 meses sin
  aviso); anual si están a 330–400 días.
- **Servicio conocido**: con **≥ 2 cobros** que cumplan el intervalo es suscripción. Con **1 solo
  cobro** es **posible** y el usuario decide (con 1 cobro no se sabe si es mensual o anual).
- **Comercio desconocido**: nunca es seguro; es **posible** si tiene ≥ 2 cobros con cadencia y monto
  consistentes, o 1 cobro con descriptor de suscripción (`LooksLikeSubscription`). El usuario la
  confirma o la descarta. Si no, va a otros cobros.
- **Activa**: último cobro hace ≤ 45 días (≤ 400 si es anual: un año + margen por una renovación
  atrasada), contando desde **hoy** (reloj inyectado). Una posible sin cadencia (1 solo cobro) usa
  ≤ 365 días. El objetivo es ver lo que se cobra hoy: una posible fuera de su ventana es inactiva y
  no pide decisión.
- **Monto, servicio conocido**: el ±25 % **no separa**; el monto actual es el del último cobro y un
  salto > 25 % se muestra como **cambio de plan** (ej. Claude USD 20 → USD 100). Los cobros del
  servicio que no siguen el patrón (ej. créditos de API) van a **otros cobros del servicio**, fuera
  del gasto mensual. Un cobro es del patrón si otro lo **respalda** (misma moneda, monto a ±25 % e
  intervalo mensual o anual), o si está a un intervalo de la cadencia del patrón del cobro respaldado
  más cercano (así se mantiene un cambio de plan reciente). Si ningún cobro tiene respaldo, un cobro
  a menos de 24 días del anterior es otro cobro. Un cambio de moneda es cambio de plan.
- **Monto, comercio desconocido**: cada cobro se compara con el anterior. Con a lo sumo un
  salto > 25 % el monto es consistente y ese salto es el cambio de plan; con más de uno, no es suscripción.
- **Umbrella** (Apple, Google Play…): sin regla propia; sus suscripciones salen de las líneas.

**`ReconcileCharges(months)`** (post-MVP: necesita correo y cartola; vista diferenciada, para ver qué cobros no se avisan):
1. Lee los cargos del período.
2. Cruza cartola con correo usando la misma regla de coincidencia.
3. Devuelve tres grupos: **avisados** (en ambos), **sin aviso** (solo en cartola: el hallazgo
   importante) y **solo correo** (aún no aparece en una cartola importada).

**Regla de coincidencia** (una sola función en `domain`, usada por ambos casos de uso): mismo banco,
mismo servicio, mismo monto y ±35 días; si ambos traen los últimos 4 dígitos, deben ser iguales.

El catálogo es el mismo para todos los bancos, porque el comercio ("NETFLIX.COM") no depende del banco.

### Cifrado (post-MVP)

> Solo aplica cuando se guarden contraseñas de aplicación (lectura de correo). El Keychain sirve
> si la app corre en un Mac; con alojamiento, la llave va a otro lado (decisión abierta).

- AES-256-GCM, nonce aleatorio de 96 bits antepuesto al texto cifrado.
- Llave de 32 bytes creada la primera vez y guardada en el **Keychain** (servicio `mail-checker`,
  cuenta `encryption-key`). Nunca se escribe junto a la base.
- Si existe `MAILCHECKER_KEY` (base64), se usa esa (Docker/Linux).
- Si se pierde la llave, hay que volver a ingresar las contraseñas de aplicación; los cargos no se pierden.

### Modelo de datos (SQLite)

| Tabla | Feature dueña | Columnas clave | Notas |
|---|---|---|---|
| `users` | auth | `id`, `public_id` UNIQUE, `name`, `email` UNIQUE, `accepted_terms_at`, `created_at` | |
| `login_codes` | auth | `id`, `user_id` UNIQUE, `code_hash`, `expires_at`, `attempts` | un código por usuario; `ON DELETE CASCADE` |
| `code_sends` | auth | `id`, `email`, `sent_at` | por correo, exista o no la cuenta |
| `sessions` | auth | `id`, `user_id`, `token_hash` UNIQUE, `expires_at` | `ON DELETE CASCADE` |
| `accounts` (post-MVP) | accounts | `id`, `public_id` UNIQUE, `user_id`, `email`, `password_enc`, `last_sync_at`, `last_error_key` | `UNIQUE(user_id, email)` |
| `statements` | statements | `id`, `public_id` UNIQUE, `user_id`, `bank`, `filename`, `kind`, `billing_date`, `imported_at`, `bad_lines` (JSON) | `UNIQUE(user_id, bank, billing_date)` |
| `charges` | shared | `id`, `user_id`, `bank`, `source`, `external_id`, `account_id`?, `statement_id`?, `occurred_at`, `amount`, `currency`, `merchant`, `card`, `card_last4` (vacío o 4 dígitos) | `UNIQUE(user_id, source, external_id)`; `ON DELETE CASCADE` desde la cuenta y desde la cartola; `account_id` llega con la tabla `accounts` (post-MVP) |
| `subscription_decisions` | subscriptions | `id`, `user_id`, `group_key` (clave de línea; `UNIQUE(user_id, group_key)`), `decision` (`confirmed`/`dismissed`), `cadence` (`monthly`/`yearly`, solo si confirmada), `decided_at` | |
| `line_aliases` | subscriptions | `id`, `user_id`, `line_key`, `alias` (≤ 40 caracteres) | `UNIQUE(user_id, line_key)` |
| `unparsed_emails` (post-MVP) | emailsync | `id`, `external_id` UNIQUE, `account_id`, `bank`, `subject`, `received_at`, `preview` | |
| `manual_subscriptions` | subscriptions | `id`, `public_id` UNIQUE, `user_id`, `name`, `amount`, `currency`, `billing_day` | `CHECK(billing_day BETWEEN 1 AND 31)` |
| `line_overrides` | subscriptions | `id`, `user_id`, `line_key`, `amount`?, `billing_day`? | `UNIQUE(user_id, line_key)` |
| `price_alert_decisions` | subscriptions | `id`, `user_id`, `line_key`, `amount` (el monto nuevo) | `UNIQUE(user_id, line_key, amount)`: historial de alzas mantenidas; no vuelve a avisar un monto ya mantenido. El resumen muestra siempre el último monto cobrado |
| `exchange_rates` | statements escribe, subscriptions lee | `id`, `user_id`, `pesos_per_dollar`, `date` | `UNIQUE(user_id, date)`; `pesos_per_dollar > 0`; cada feature con su adaptador, como `charges` |
| `hidden_lines` | subscriptions | `id`, `user_id`, `line_key`, `hidden_at` | `UNIQUE(user_id, line_key)`; la usa "dejar de seguir" |
| `cancellations` | subscriptions | `id`, `user_id`, `line_key`, `marked_at` | `UNIQUE(user_id, line_key)` |

Toda tabla de datos lleva `user_id` con `ON DELETE CASCADE` desde `users`.

- `id`: `INTEGER PRIMARY KEY AUTOINCREMENT` en todas las tablas; solo para uso interno (joins, FKs,
  `external_id`), nunca sale por la API.
- `public_id`: `TEXT UNIQUE NOT NULL`, UUID v7 generado en Go por el caso de uso (`NewPublicID`). Solo en las
  tablas que la API referencia (`users`, `accounts`, `statements`, `manual_subscriptions`).
- Fechas: `INTEGER` unix en segundos, UTC. Una fecha sin hora (`billing_date`) se guarda como las
  00:00 de ese día en `America/Santiago`. La conversión `time.Time` ↔ unix vive solo en los
  adaptadores sqlite y http.
- Migraciones: centralizadas en `shared/sqlite/migrations`, porque es una sola base.
- Consultas: cada feature tiene las suyas en `adapters/sqlite/queries.sql`, y sqlc genera el código Go
  en `adapters/sqlite/queries` (paquete `queries`). `backend/sqlc.yaml` toma el esquema de las
  migraciones y lleva una entrada por feature (`omit_unused_structs`: cada paquete solo trae los
  modelos que usa). sqlc v1.31.1 se corre con `go run`, sin instalarlo; lo generado se commitea y
  `sqlc diff` comprueba que está al día.
- Un `UNIQUE` violado se traduce a `ErrConflict` con `sqlite.IsUniqueViolation`. Los tests de repos
  usan una base temporal real ya migrada (`sqlitetest.OpenMigrated`).
- Conexión: `PRAGMA journal_mode=WAL`, `foreign_keys=ON`, `busy_timeout=5000`, `synchronous=FULL`;
  un solo escritor. Si algo falla al abrir la base o migrar, la app no arranca.
- Transacciones: viven dentro de los métodos de los repos (adaptador). Los casos de uso no las
  manejan. La sincronización usa una transacción por cuenta.
- Integridad en la base: `NOT NULL` y `CHECK` en montos (`amount >= 0`) y monedas (`currency IN ('CLP','USD')`).

### API REST

Toda ruta exige sesión salvo `/api/health` y `/api/auth/*`; una ruta que no existe responde 404 sin
pedir sesión.

- Cada petición pasa por: log (método, ruta, status y duración; nunca query ni cuerpo; 5xx como
  error), recuperación de panics (500 `common.internal`, detalle solo al log) y comprobación de
  origen.
- Sesión: cookie `session` `HttpOnly`, `SameSite=Lax`, `Path=/`, 30 días; `Secure` en producción.
- CSRF: `POST`, `PUT`, `PATCH` y `DELETE` con un `Origin` distinto del `Host` de la petición
  responden 403 `common.forbidden`; sin `Origin` pasan (los navegadores lo mandan siempre en un POST
  y la cookie `Lax` no viaja desde otro sitio).
- Sin CORS: el front se sirve desde el mismo origen; no se manda `Access-Control-Allow-Origin`.
- Servidor con tiempos máximos de lectura de cabeceras, lectura, escritura e inactividad.
- Errores de validación por campo: `InvalidFields` → `{error:{key:"common.validation_failed", fields:{campo: llave}}}`.
- Rutas de una línea: la clave (`lineKey` o `groupKey`) ausente o en blanco es 400 con
  `fields:{lineKey|groupKey: "subscriptions.line_key_required"}`; una clave que no es línea del
  usuario es 404 `subscriptions.line_not_found`.

| Método | Ruta | Feature | Descripción |
|---|---|---|---|
| GET | `/api/health` | — | `{status:"ok"}` si la base responde en 1 s; si no, 503 `common.unavailable` (la causa solo al log) |
| POST | `/api/auth/register` | auth | `{name, email, acceptedTerms}`; envía el código; 204 |
| POST | `/api/auth/login` | auth | `{email}`; envía el código; 204 exista o no la cuenta |
| POST | `/api/auth/code/verify` | auth | `{email, code}`; abre la sesión y fija la cookie; 204. Un código que no son 6 dígitos es 400 `auth.code_format_invalid` y no gasta intento |
| POST | `/api/auth/code/resend` | auth | `{email}`; igual que login (el código nuevo reemplaza al anterior); 204 |
| POST | `/api/auth/logout` | auth | cierra la sesión y borra la cookie; 204 siempre |
| GET | `/api/me` | auth | `{id, name, email, initials}` (`id` = `public_id`; iniciales de las dos primeras palabras) |
| DELETE | `/api/me` | auth | elimina al usuario y todo lo suyo en una transacción; borra la cookie; 204 |
| GET | `/api/export?format=csv` | auth + puente | copia de los datos (mvp.md decisión 13), como adjunto `fokko-copia.csv` o `.json`. `csv` (por defecto): las suscripciones del resumen, columnas `name, alias, amount` (2 decimales)`, currency, cadence, nextCharge` (AAAA-MM-DD); un campo que empieza con `=`, `+`, `-`, `@`, tab o CR lleva `'` delante. `json`: `{profile, statements, charges, decisions, aliases, manuals, overrides}` con ids públicos, fechas unix s y montos `{amount, currency}`. Otro formato: 400 `fields:{format:"auth.export_format_unknown"}` |
| GET | `/api/banks` | statements | `[{id}]` con los bancos registrados; el nombre visible lo pone el front con i18n |
| POST | `/api/statements` | statements | multipart: `file` (uno), `billingDate` opcional (unix s), `bank` opcional; 201 `{id, filename, months:["AAAA-MM"], charges, badLines}` (`badLines` es la cantidad). Hasta 10 MB (`statements.file_too_large`) y 60 s para leer e importar; no multipart es `common.invalid_request` |
| GET | `/api/summary?months=12` | subscriptions | `months` 1 a 24 (12 si falta). `{subscriptions, possible, total, charged, remaining, rises, pendingActions, upcoming}`: montos `{amount, currency}` en centésimas, fechas unix s (`null` si no hay), listas vacías como `[]`. Sin inactivas, otros cobros ni ids internos |
| POST | `/api/subscriptions/decisions` | subscriptions | `{groupKey, decision, cadence?}` (`confirm` con `monthly`/`yearly`, o `dismiss`): confirmar o descartar una posible; 204. Una clave que no es posible hoy es 404 `subscriptions.possible_not_found` |
| DELETE | `/api/subscriptions/decisions?groupKey=` | subscriptions | deshace la decisión; la línea vuelve a posibles; 204 aunque no hubiera decisión |
| PUT | `/api/subscriptions/aliases` | subscriptions | `{lineKey, alias}`: nombra una línea; alias vacío lo borra; 204 |
| PUT | `/api/subscriptions/overrides` | subscriptions | `{lineKey, amount?, billingDay?}` (monto en centésimas, en la moneda de la línea): reemplaza lo detectado; sin ninguno lo borra; 204 |
| POST | `/api/subscriptions/manual` | subscriptions | `{name, amount, currency, billingDay}` (monto en centésimas); 201 `{id}` |
| PUT, DELETE | `/api/subscriptions/manual/{publicId}` | subscriptions | edita (mismo cuerpo que crear) o quita una manual; 204. Ajena, inexistente o id mal formado: 404 |
| DELETE | `/api/subscriptions/line?lineKey=` | subscriptions | "dejar de seguir": oculta la línea (mvp.md decisión 10); 204, también si ya estaba oculta |
| POST | `/api/subscriptions/price-alerts/decision` | subscriptions | `{lineKey}`: el usuario mantiene la suscripción tras el alza; 204. Sin alza vigente: 404 `subscriptions.rise_not_found` |
| GET | `/api/subscriptions/cancel-guide?lineKey=` | subscriptions | `{url}` del sitio del servicio (catálogo); `url` vacío si el comercio no tiene enlace |
| POST | `/api/subscriptions/cancellations` | subscriptions | `{lineKey}`: marca la línea como "por cancelar"; 204 |

**Post-MVP:**

| Método | Ruta | Feature | Descripción |
|---|---|---|---|
| GET, POST | `/api/accounts` | accounts | lista (sin contraseñas) / `{email, appPassword}` |
| DELETE | `/api/accounts/{publicId}` | accounts | borra la cuenta y sus cargos |
| POST | `/api/sync` | emailsync | `{months}` → resultado por cuenta |
| GET | `/api/unparsed` | emailsync | correos no reconocidos (revisar contra "solo servicio, monto y fecha") |
| GET | `/api/statements` | statements | |
| DELETE | `/api/statements/{publicId}` | statements | borra la cartola y sus cargos |
| GET | `/api/reconciliation?months=12` | subscriptions | avisados, sin aviso, solo correo |

Las claves de línea llevan `|`, `:`, `#` y espacios, así que nunca van en el path: viajan en el
cuerpo o, en el `DELETE`, como parámetro de consulta codificado.

En el JSON, `id` es siempre el `public_id` (UUID v7) y las fechas son unix en segundos.

**El back no devuelve textos para el usuario: devuelve llaves i18n** y el front las traduce.

```json
{ "error": { "key": "statements.bank_not_detected", "params": { "count": 2 } } }
{ "error": { "key": "common.validation_failed", "fields": { "email": "accounts.email_invalid" } } }
```

- Llaves `<feature>.<motivo>` en inglés y snake_case; las genéricas en `common.*`. Cada llave es una
  constante en el `domain` de su feature (las comunes en `shared/domain`).
- Status: 400 inválido, 404 no existe, 409 conflicto o sincronización en curso, 500 interno
  (`common.internal`, sin detalles; el error real va al log).
- Toda falla que ve el usuario viaja como llave, también dentro de respuestas exitosas
  (ej. `SyncResult.errorKey`).
- Los logs y los errores de arranque sí van en texto humano: no los lee el front.

### Agregar un banco (ej. Falabella)

1. `statements/adapters/falabella/`: implementa `BankStatementParser` (`Detect`, `BillingDate` y parser de su cartola).
2. Tests con cartolas reales (con nombres anonimizados).
3. Registrarlo en `main.go`.
4. **Post-MVP:** `emailsync/adapters/falabella/` con `BankEmailParser` (filtro y parser de sus avisos).

No se modifica ningún caso de uso, ni el dominio, ni la base, ni el frontend.

### Tests

- `domain` y `usecases` de cada feature: con dobles en memoria de sus puertos (sin red ni base).
- Parsers de banco, catálogo y extractor: con textos reales anonimizados.
- `adapters/sqlite`: contra una base temporal.
- `adapters/http`: con `httptest`.

---

## Convenciones Go

| Qué | Convención | Ejemplo |
|---|---|---|
| Paquetes y carpetas | minúscula, una palabra, sin `_` ni `-` | `emailsync`, `usecases` |
| Archivos | minúscula, snake_case | `sync_all.go`, `sync_all_test.go` |
| Funciones, tipos, variables | camelCase; mayúscula inicial si es exportado | `SyncAll`, `BankEmailParser`, `parseAmount` |
| Constantes | igual que variables | `maxSyncMonths` |
| Errores | variables `ErrXxx`; mensajes en minúscula | `ErrNotFound` |

---

## Frontend

### Estructura del frontend

Clean architecture por feature, igual que el back: las dependencias apuntan hacia adentro
(`presentation → application → domain`; `infrastructure → application`).

```
frontend/src/
├── main.tsx                   composición: crea los gateways reales y monta `GatewaysProvider`
├── routes/                    TanStack Router (file-based); cada ruta monta un contenedor de `screens`
├── components/
│   ├── ui/                    componentes de base presentacionales, por Atomic Design (lista en `docs/mvp.md`):
│   │   ├── atoms/             una sola pieza: Button, Input, Checkbox, Badge, ServiceAvatar, Skeleton, AnimatedNumber
│   │   ├── molecules/         atoms juntos con un fin: FormField, CodeInput, ResendTimer, ListRow, SegmentedControl, Toast, EmptyState, Card
│   │   └── organisms/         bloques completos y genéricos: SegmentedRing
│   └── layout/                templates: AppShell, BrandHeader, BottomNav, SideNav (presentacionales)
├── styles/                    index.css (base) y tokens.css (tokens raw y semantic)
├── features/<feature>/        auth · statements · subscriptions · profile
│   ├── domain/                tipos `readonly`, funciones puras y schemas zod; sin React ni red
│   ├── application/           hooks de caso de uso (`useSummary`, `useRegister`) y puertos (gateways)
│   ├── infrastructure/        implementa los puertos con `lib/api`; mappers DTO → dominio
│   └── presentation/
│       ├── components/        presentacionales
│       └── screens/           contenedores: llaman hooks y pasan props
└── lib/
    ├── api/                   client.ts, types.ts
    ├── gateways/              `GatewaysProvider`, `useGateways`, `renderWithGateways` (pruebas)
    ├── format.ts              montos y fechas (es-CL)
    ├── i18n/                  es.ts (catálogo tipado de llaves) + t(key, params)
    └── theme/                 preferencia claro/oscuro/sistema (localStorage; si falla, sistema)
```

- **Atomic Design:** atoms → molecules → organisms en `components/ui`; los templates son
  `components/layout`; los organismos propios de una feature viven en su
  `presentation/components` y las páginas son los contenedores de `screens`. Una pieza solo
  importa piezas de su nivel o de niveles menores.

- **Componente presentacional:** recibe datos y callbacks por props; puede usar `t()` y `react`, nada más.
  No usa `useQuery`, router, `fetch` ni importa `application` o `infrastructure`. Con las mismas props
  pinta lo mismo, así que se prueba sin servidor y se extrae copiando su carpeta.
- **Contenedor (`screens`):** único lugar que llama hooks de `application`; sin lógica inline en el JSX.
- **Puertos:** se entregan por Context (`useGateways`). Aunque cada puerto tenga una sola implementación
  real, las pruebas usan un doble en memoria, así que se permiten (como en el back).
- **Guardia:** Biome con `noRestrictedImports` por carpeta rompe el lint al cruzar un límite
  (ideas `FOKKO-0085C-SUBS` y `FOKKO-0085D-SUBS`).

### Mobile first

- Las clases base de Tailwind son para ~390 px; desde `lg:` (1024 px) se ve el diseño de escritorio.
- **Celular:** el diseño del handoff (`docs/mvp.md`): header de marca, navegación inferior fija
  (Resumen · Subs), hoja con esquinas redondeadas y objetivos táctiles de al menos 44 px.
- **Escritorio (`lg:`, 1024 px):** barra lateral con Resumen y Subs y contenido en dos columnas (detalle en
  `docs/mvp.md`). La tablet (768–1023 px) no está diseñada: usa el layout de escritorio con barra inferior.
- Cada cargo y suscripción muestra su banco y tarjeta (ej. "BCI · crédito").
- Tema claro/oscuro con tokens propios (`src/styles/tokens.css`), sin librería de UI:
  - **Raw:** la paleta aprobada de Fokko (Púrpura, Coral, Emerald, Gold, Onyx; 12 pasos OKLCH por
    esquema: 1–2 fondos, 3–5 componentes, 6–8 bordes, 9–10 sólidos, 11–12 texto) y el crema de la
    página clara (`sand`), que no está en la paleta. Son variables CSS sin clases de Tailwind.
  - **Semantic:** para qué es el color (`page`, `ink`, `primary`, `alert`, `danger`…). Son los únicos
    colores que conoce Tailwind (`@theme` con `--color-*: initial`) y los únicos que usan los
    componentes. Cada esquema los apunta a un paso raw; salen del `tokens.json` del handoff, llevado
    al paso más cercano de la paleta.
  - Un test falla si fuera de `tokens.css` aparece `--raw-*`, un hex, `rgb()`/`hsl()`/`oklch()` o una
    clase de la paleta por defecto de Tailwind.
  - Oscuro: `data-theme="dark"` en `<html>`, o el sistema cuando la preferencia es "sistema".
  - Radios `sm` 6, `md` 10, `lg` 14, `sheet` 28 px; fuentes Be Vietnam Pro (cuerpo) y Prompt 600
    (wordmark), servidas desde el propio binario (`@fontsource`), sin pedir nada a terceros.

Los componentes del MVP están listados en `docs/mvp.md`. Se escriben cuando la pantalla que los
usa los pida (YAGNI), no todos de golpe. Si hace falta un overlay (diálogo, hoja), hay que cuidar
foco, Esc y scroll lock a mano; si eso resulta frágil, se propone agregar Radix Primitives como
dependencia aprobada.

### Desarrollo y build

- Desarrollo: Vite en `:5173` con proxy `/api` → Go en `:8080` (el puerto por defecto de `MAILCHECKER_PORT`).
- Build: `pnpm build` deja el front en `backend/internal/shared/httpserver/web/dist` y Go lo
  embebe, así queda **un solo binario**.

---

## Plan

El plan vive en `docs/backlog/README.md`: ideas pequeñas `FOKKO-XXXX-SUBS`, cada una con tareas
unitarias `TASK-XXX`, agrupadas en bloques (catálogo y cartolas, detección, acceso, persistencia,
HTTP, front, cierre). Cubre el 100 % de `docs/mvp.md`. Se trabaja con `/idea` → `/stories` →
`/story`; al terminar cada tarea se detiene para revisión.

**Alojamiento:** todo va dockerizado (Dockerfile multi-stage, un solo binario con el front
embebido, `docker-compose` local) para que el despliegue no dependa de la plataforma; el proveedor
queda abierto (`docs/mvp.md`, decisión 20). El correo de los códigos de acceso se envía con Resend
por HTTP (`net/http` con timeout, sin SDK), detrás de un puerto; en desarrollo, un doble de consola.

### Post-MVP

- Lectura de correo por IMAP: `accounts`, `emailsync`, cifrado, sincronización, parser de avisos
  BCI, conciliación cartola-correo y las pantallas Cuentas de correo, Añadir cuenta y Fuentes.
- Evaluar OAuth de Gmail con usuarios reales (exige verificación y auditoría de Google).
- Login con Google y Apple, calendario, notificaciones, recuperar acceso, uso sin conexión.

## Pendiente

Ajustes al código ya escrito, por las decisiones del MVP:

- **`UserID`** en `Charge` y en todos los puertos de datos (ver "Multiusuario").
- **`BankStatementParser.BillingDate`**: sacar la fecha de facturación del archivo.

Otros:

- **Deshacer una decisión** (`UndoDecision` + `SubscriptionDecisions.Delete`): está en el diseño
  pero no en el código; el núcleo de suscripciones (idea 002) se cerró sin él. Está en el backlog (`FOKKO-0022-SUBS`).
- **Cartola BCI en PDF**: ya hay dos reales (sep y oct 2026) para escribir el parser (`FOKKO-0011A-SUBS` en adelante); se
  anonimizan antes de crear los fixtures. Cada una trae un estado nacional y uno internacional,
  que se leen como un solo flujo de cobros (`docs/mvp.md`, decisión 6). Falta un ejemplo de CSV
  (post-MVP) y uno con un cobro en euros u otra moneda.
- **Tipo de cambio** (`docs/mvp.md`, decisión 6): el diseño con el puerto `ExchangeRates` y su
  "fuente abierta" quedó viejo. El parser de cada banco debe devolver el monto de origen y el tipo
  de cambio que derive de la cartola (en BCI, el par de traspasos), y se guarda el último derivado
  con su fecha: falta una tabla en el modelo de datos. Se rediseña al abrir la idea de lectura de
  cartolas.
- **Apagado con conexiones sin usar**: un cliente (p. ej. un navegador que preconecta) puede dejar
  una conexión abierta que nunca manda petición; `Shutdown` la espera y a los 5 s `run` termina con
  `shutdown: context deadline exceeded` (código 1). Lo común es forzar `server.Close()` al vencer el
  plazo y no tratarlo como falla. Lo detectaron las pruebas de humo (`cmd/api/smoke_test.go`).
- **Sin diseño en el handoff:** tablet (768–1023 px) y las versiones de escritorio de Revisión vacía y Resumen vacío.
- **Post-MVP:** ocultar la contraseña de `Credential` en logs (`LogValue()` que la reemplace), para que un
  `slog.Any("cred", cred)` accidental no la exponga.
- **Post-MVP:** registrar el detalle de un error de lectura de correo en `SyncAll` (hoy solo queda la
  llave en el `SyncResult`): con `slog`, sin correo ni contraseña, cuando exista el lector IMAP real.
