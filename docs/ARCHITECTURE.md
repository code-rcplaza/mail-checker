# Mail Checker — Arquitectura v1

Detecta suscripciones a partir de (1) avisos de compra que llegan por correo, desde una o más
cuentas Gmail, y (2) cartolas de tarjeta (PDF, Excel, CSV o texto pegado). Soporta **varios
bancos**: BCI es el primero; agregar otro (ej. Falabella) es escribir adaptadores nuevos, sin
modificar los casos de uso.

Uso: personal y local. Sin login ni multiusuario.

## Decisiones

| Tema | Decisión |
|---|---|
| Backend | Go. **Features verticales**, cada una con arquitectura hexagonal: `adapters/ domain/ ports/ usecases/` |
| Bancos | Multi-banco desde el inicio: un parser por banco, registrados en `main.go` |
| Router HTTP | `go-chi/chi` |
| Base de datos | SQLite (`modernc.org/sqlite`, sin cgo) + **sqlc** para las consultas + **goose** para las migraciones |
| Correo | IMAP + contraseña de aplicación de Google |
| Cifrado | AES-256-GCM; la llave vive en el **Keychain de macOS** (`zalando/go-keyring`), con respaldo `MAILCHECKER_KEY` |
| Frontend | React + TypeScript + Vite + Tailwind v4 + **shadcn/ui**, **mobile first** |
| Routing / datos / formularios | TanStack Router + TanStack Query + react-hook-form + zod |
| Fechas | Unix en **segundos** (UTC) en la base y en la API; `time.Time` dentro de Go; la UI muestra hora de `America/Santiago` |
| IDs | `id` autoincremental interno; `public_id` UUID v7 (`google/uuid`, generado en Go) para lo que se expone |
| Fuera de alcance v1 | Suscripciones manuales, OAuth, multiusuario |

---

## Backend

### Features

Cada feature responde a "¿quién eres y qué haces?":

| Feature | Responsabilidad |
|---|---|
| `accounts` | Administro las cuentas de correo y sus credenciales cifradas |
| `emailsync` | Leo los avisos de compra de cada banco y los convierto en cargos |
| `statements` | Importo cartolas, detecto de qué banco son y las convierto en cargos |
| `subscriptions` | Detecto suscripciones a partir de los cargos y calculo el gasto mensual |
| `shared` | Lo mínimo común: el cargo, el dinero, errores y la plataforma (SQLite, HTTP, reloj) |

### Estructura

```
backend/
├── cmd/api/main.go                  composition root: crea adaptadores, registra bancos, arma rutas
└── internal/
    ├── accounts/
    │   ├── adapters/
    │   │   ├── http/                GET/POST/DELETE /api/accounts
    │   │   ├── sqlite/              AccountRepository
    │   │   └── keychain/            SecretCipher (AES-GCM + llave en Keychain)
    │   ├── domain/                  Account
    │   ├── ports/                   AccountRepository, SecretCipher
    │   └── usecases/                add_account, list_accounts, remove_account, credentials, record_sync
    │
    ├── emailsync/
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
    │   │   ├── http/                GET/POST/DELETE /api/statements, GET /api/banks
    │   │   ├── extract/             FileExtractor: pdf / xlsx / csv / texto
    │   │   ├── bci/                 BankStatementParser de BCI
    │   │   └── sqlite/              StatementRepository
    │   ├── domain/                  Statement, StatementLine
    │   ├── ports/                   StatementRepository, FileExtractor, BankStatementParser
    │   └── usecases/                import_statement, list_statements, remove_statement, list_banks
    │
    ├── subscriptions/
    │   ├── adapters/
    │   │   ├── http/                GET /api/summary
    │   │   ├── catalog/             ServiceCatalog (descriptor de cobro → servicio)
    │   │   └── sqlite/              ChargeReader
    │   ├── domain/                  Subscription, Summary, reglas de recurrencia y de duplicados
    │   ├── ports/                   ChargeReader, ServiceCatalog
    │   └── usecases/                get_summary
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

### `shared/domain`

```go
type Money int64               // centésimas: CLP 5.990 => 599000, USD 20,00 => 2000
type Currency string           // "CLP", "USD"
type Bank string               // "bci", "falabella"…
type Card string               // "credito", "debito"
type Source string             // "email", "statement"

// ParsedCharge: lo que un parser de banco extrae de un aviso o línea de cartola.
type ParsedCharge struct {
    Amount   Money
    Currency Currency
    Merchant string
    Card     Card
    Date     time.Time // cero si la línea de cartola no trae fecha
}

// Charge: un cargo guardado.
type Charge struct {
    ID          int64
    Bank        Bank
    Source      Source
    ExternalID  string  // correo: Message-ID; cartola: "<statementID>:<línea>"
    AccountID   *int64
    StatementID *int64
    OccurredAt  time.Time
    Amount      Money
    Currency    Currency
    Merchant    string
    Card        Card
}

var ErrInvalid, ErrNotFound, ErrConflict error // + Invalid(msg) con mensaje para el usuario
```

### `accounts`

```go
// domain
type Account struct {
    ID          int64      // interno, nunca sale por la API
    PublicID    string     // UUID v7, el que se expone
    Email       string
    PasswordEnc []byte
    CreatedAt   time.Time
    LastSyncAt  *time.Time
    LastError   string
}

// ports
type AccountRepository interface {
    Create(ctx context.Context, a *Account) error             // asigna ID y PublicID; ErrConflict si el correo ya existe
    List(ctx context.Context) ([]Account, error)
    Delete(ctx context.Context, publicID string) error        // ErrNotFound; borra también sus cargos
    MarkSynced(ctx context.Context, id int64, at time.Time, syncErr string) error
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

### `emailsync`

```go
// domain
type SyncResult struct {
    AccountPublicID string
    Email     string
    Found, New, Unparsed int
    Error     string
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
    RecordSync(ctx context.Context, accountID int64, at time.Time, syncErr string) error
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
    Kind        FileKind   // pdf, xlsx, csv, text
    BillingDate time.Time
    ImportedAt  time.Time
    Charges     int
    BadLines    []string
}

// ports
type FileExtractor interface {
    Extract(filename string, data []byte) (text string, kind FileKind, err error)
}
type BankStatementParser interface {
    Bank() Bank
    Detect(text string) bool                                // ¿esta cartola es de mi banco?
    Parse(text string) (lines []ParsedCharge, bad []string)
}
type StatementRepository interface {
    Create(ctx context.Context, s *Statement, charges []Charge) error   // atómico; asigna ID y PublicID
    List(ctx context.Context) ([]Statement, error)
    Delete(ctx context.Context, publicID string) error                  // ErrNotFound; borra también sus cargos
}
```

**`ImportStatement(file, billingDate, bank?)`:**
1. Extrae el texto del archivo.
2. Banco: si viene elegido, usa ese. Si no, pregunta a cada parser `Detect`. Si exactamente uno
   lo reconoce, usa ese; si no, responde *"no se pudo detectar el banco, elígelo"*.
3. Parsea. Las líneas sin fecha toman la fecha de facturación.
4. Guarda la cartola y sus cargos en una transacción.

### `subscriptions`

```go
// ports
type ChargeReader interface {
    ListSince(ctx context.Context, since time.Time) ([]Charge, error)
}
type ServiceCatalog interface {
    Identify(merchant string) (Service, bool)   // "ANYSPHERE INC" -> Cursor
    Normalize(merchant string) string           // "C. VERDE MALL" -> "C VERDE"
    LooksLikeSubscription(merchant string) bool // "…SUBSCRIPTION", "MEMBERSHIP"
}
```

**`GetSummary(months)`:**
1. Lee los cargos del período.
2. Quita los cargos de cartola que ya llegaron por correo: **mismo banco**, mismo servicio y monto,
   ±35 días.
3. Agrupa:
   - servicio conocido → un grupo por servicio (sin importar moneda ni banco);
   - servicio *umbrella* (Apple, Google Play…) → un grupo por monto;
   - comercio desconocido → por nombre normalizado y moneda.
4. Aplica las reglas de recurrencia del dominio: mensual o anual, tolerancia de ±25 %, meses
   faltantes, activa si tuvo cargo hace menos de 45 días (400 días si es anual).
5. Suma el gasto mensual de las activas, por moneda.

El catálogo es el mismo para todos los bancos, porque el comercio ("NETFLIX.COM") no depende del banco.

### Cifrado

- AES-256-GCM, nonce aleatorio de 96 bits antepuesto al texto cifrado.
- Llave de 32 bytes creada la primera vez y guardada en el **Keychain** (servicio `mail-checker`,
  cuenta `encryption-key`). Nunca se escribe junto a la base.
- Si existe `MAILCHECKER_KEY` (base64), se usa esa (Docker/Linux).
- Si se pierde la llave, hay que volver a ingresar las contraseñas de aplicación; los cargos no se pierden.

### Modelo de datos (SQLite)

| Tabla | Feature dueña | Columnas clave | Notas |
|---|---|---|---|
| `accounts` | accounts | `id`, `public_id` UNIQUE, `email` UNIQUE, `password_enc`, `last_sync_at`, `last_error` | |
| `statements` | statements | `id`, `public_id` UNIQUE, `bank`, `filename`, `kind`, `billing_date`, `bad_lines` (JSON) | |
| `charges` | shared | `id`, `bank`, `source`, `external_id`, `account_id`?, `statement_id`?, `occurred_at`, `amount`, `currency`, `merchant`, `card` | `UNIQUE(source, external_id)`; `ON DELETE CASCADE` desde la cuenta y desde la cartola |
| `unparsed_emails` | emailsync | `id`, `external_id` UNIQUE, `account_id`, `bank`, `subject`, `received_at`, `preview` | |

- `id`: `INTEGER PRIMARY KEY AUTOINCREMENT` en todas las tablas; solo para uso interno (joins, FKs,
  `external_id`), nunca sale por la API.
- `public_id`: `TEXT UNIQUE NOT NULL`, UUID v7 generado en Go por el repo en `Create`. Solo en las
  tablas que la API referencia (`accounts`, `statements`).
- Fechas: `INTEGER` unix en segundos, UTC. Una fecha sin hora (`billing_date`) se guarda como las
  00:00 de ese día en `America/Santiago`. La conversión `time.Time` ↔ unix vive solo en los
  adaptadores sqlite y http.
- Migraciones: centralizadas en `shared/sqlite/migrations`, porque es una sola base.
- Consultas: cada feature tiene las suyas en `adapters/sqlite/queries.sql`, y sqlc genera el código Go.
- Conexión: `PRAGMA journal_mode=WAL`, `foreign_keys=ON`, `busy_timeout=5000`, `synchronous=FULL`;
  un solo escritor. Si algo falla al abrir la base o migrar, la app no arranca.
- Transacciones: viven dentro de los métodos de los repos (adaptador). Los casos de uso no las
  manejan. La sincronización usa una transacción por cuenta.
- Integridad en la base: `NOT NULL` y `CHECK` en montos (`amount >= 0`) y monedas (`currency IN ('CLP','USD')`).

### API REST

| Método | Ruta | Feature | Descripción |
|---|---|---|---|
| GET | `/api/health` | — | |
| GET | `/api/accounts` | accounts | lista (sin contraseñas) |
| POST | `/api/accounts` | accounts | `{email, appPassword}` |
| DELETE | `/api/accounts/{publicId}` | accounts | borra la cuenta y sus cargos |
| POST | `/api/sync` | emailsync | `{months}` → resultado por cuenta |
| GET | `/api/unparsed` | emailsync | correos no reconocidos |
| GET | `/api/banks` | statements | bancos soportados (para el selector) |
| GET | `/api/statements` | statements | |
| POST | `/api/statements` | statements | multipart: `file` o `text`, + `billingDate`, + `bank` opcional |
| DELETE | `/api/statements/{publicId}` | statements | borra la cartola y sus cargos |
| GET | `/api/summary?months=12` | subscriptions | suscripciones, otros cargos, totales |

En el JSON, `id` es siempre el `public_id` (UUID v7) y las fechas son unix en segundos.

Errores: `{ "error": "mensaje" }` → 400 inválido, 404 no existe, 409 conflicto o sincronización en
curso, 500 interno.

### Agregar un banco (ej. Falabella)

1. `emailsync/adapters/falabella/`: implementa `BankEmailParser` (filtro y parser de sus avisos).
2. `statements/adapters/falabella/`: implementa `BankStatementParser` (`Detect` y parser de su cartola).
3. Tests de ambos con correos y cartolas reales (con nombres anonimizados).
4. Registrarlos en `main.go`.

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

### Estructura (refleja las features del backend)

```
frontend/src/
├── main.tsx
├── routes/                    TanStack Router (file-based)
│   ├── __root.tsx             AppShell
│   ├── index.tsx              Resumen
│   ├── cuentas.tsx
│   └── cartolas.tsx
├── components/
│   ├── ui/                    generados por shadcn
│   └── layout/                AppShell, BottomNav, Sidebar
├── features/
│   ├── subscriptions/         SubscriptionCard, SubscriptionTable, TotalsTiles, useSummary
│   ├── accounts/              AccountList, AddAccountDrawer, useAccounts
│   ├── emailsync/             SyncButton, SyncResults, UnparsedList, useSync
│   └── statements/            UploadDropzone, BankSelect, PasteForm, StatementList, useStatements
└── lib/
    ├── api/                   client.ts, types.ts
    ├── format.ts              montos y fechas (es-CL)
    └── utils.ts               cn()
```

### Mobile first

- Las clases base de Tailwind son para ~375 px; `md:` y `lg:` agregan el layout de escritorio.
- **Celular:**
  - header compacto y navegación inferior fija (Resumen · Cuentas · Cartolas);
  - suscripciones como `Card` apiladas;
  - formularios en `Drawer` desde abajo;
  - objetivos táctiles de al menos 44 px.
- **md+:** barra lateral, `Table` de shadcn, formularios en `Dialog`/`Sheet`.
- Cada cargo y suscripción muestra su banco y tarjeta (ej. "BCI · crédito").
- Tema claro/oscuro con las variables CSS de shadcn.

Componentes shadcn previstos: `button`, `card`, `table`, `badge`, `drawer`, `dialog`, `sheet`, `form`,
`input`, `select`, `tabs`, `skeleton`, `sonner` (toasts), `collapsible`.

### Desarrollo y build

- Desarrollo: Vite en `:5173` con proxy `/api` → Go en `:8765`.
- Build: `pnpm build` deja el front en `backend/internal/shared/httpserver/web/dist` y Go lo
  embebe, así queda **un solo binario**.

---

## Plan por fases

Al terminar cada fase **paramos** y revisamos antes de seguir.

1. **Reorganizar el núcleo en features**: `shared/domain` y, por cada feature, `domain` + `ports` +
   `usecases` con tests usando dobles. Reutiliza la lógica y los tests ya escritos e incorpora el
   esquema multi-banco.
2. **Adaptadores de banco y lectura**: `bci` (avisos y cartola), `extract`, `catalog`, reutilizando
   los parsers y sus tests.
3. **Persistencia y cifrado**: `shared/sqlite` (goose), repositorios con sqlc y `keychain`.
4. **IMAP + HTTP + `main.go`**: prueba real con tu Gmail.
5. **Front, base**: Vite + Tailwind + shadcn + TanStack Router + AppShell mobile first.
6. **Pantalla Resumen.**
7. **Pantalla Cuentas** (agregar, quitar, sincronizar).
8. **Pantalla Cartolas** (subir archivo, pegar texto, elegir banco).
9. **Build único**, Makefile, README.

## Pendiente

- **PDF real de cartola BCI** para ajustar el parser (fase 2 u 8).
