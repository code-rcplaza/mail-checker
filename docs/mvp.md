# MVP — pantallas, endpoints y componentes

Alcance del MVP según el handoff de diseño Fokko Subs: entran solo las pantallas marcadas con ★.
El handoff es referencia visual y de flujo; el front se construye con componentes propios en
Tailwind (sin librería de UI).

Este documento define lo que entra en la primera versión. Lo que no aparece aquí queda fuera.

**Fuente de datos del MVP: solo cartolas en PDF.** El CSV y la lectura de correo van en el post-MVP (ver al final).

## Pantallas

| Entra al MVP (★) | Opcional (☆) | Después |
|---|---|---|
| Registro, Login, Código, Bienvenida | Revisión vacía | Calendario |
| Elegir fuente, Cargando, Análisis, Error | Resumen vacío | Notificaciones, Permiso de notificaciones |
| Revisión, Resumen, Alerta, Cómo cancelar | | Cuentas de correo, Añadir cuenta |
| Subs, Detalle, Editar, Agregar | | Fuentes, Recuperar acceso |
| Perfil, Eliminar todo | | Sin conexión |

Flujos:

- **Acceso:** Registro ↔ Login → Código → Bienvenida.
- **Primera vez:** Bienvenida → Elegir fuente (solo subir cartola) → Cargando → Análisis → Revisión →
  Resumen (Error si el archivo no se puede leer).
- **Ciclo continuo:** Resumen → Alerta (alza) → Cómo cancelar → Resumen.
- **Suscripción:** Subs → Detalle → Editar; Subs → Agregar; Revisión → Editar.
- **Cuenta:** avatar → Perfil → Apariencia / Tus datos (Eliminar todo) / Cerrar sesión.

Cambios al diseño por las decisiones:

- Registro y Login sin los botones de Google y Apple.
- Cómo cancelar sin el paso de reenviar el correo de cancelación: el usuario marca "por cancelar".
- Elegir fuente sin la opción de reenviar correos: solo subir cartola.
- Resumen sin "último correo recibido".
- Inactivas y Otros cobros no se muestran.

## Escritorio

El handoff trae referencias de escritorio (`desk_*`) para las 18 pantallas ★. Solo faltan las
versiones de escritorio de Revisión vacía y Resumen vacío (☆). Las de Permiso de notificaciones y
Recuperar acceso existen pero quedan fuera del MVP.

- **Corte en 1024 px (`lg:`).** Desde ahí se ve el diseño de escritorio; por debajo, el móvil con
  barra inferior. No es el móvil estirado: es una composición propia.
- **Barra lateral (≥1024 px):** 256 px, fija, borde derecho. Marca arriba; navegación Resumen y Subs
  (Calendario y Fuentes son de después); pie con avatar, nombre y correo que lleva a Perfil.
- **Contenido con sesión:** máximo 1200 px, dos columnas con 48 px de separación. Izquierda: hero en
  tarjeta con el anillo (240 px; Bienvenida 260 px). Derecha: el contenido de la pantalla.
- **Acceso (Registro, Login, Código):** sin barra lateral ni inferior; encabezado con marca y botón
  de tema, y la misma composición de dos columnas.
- **Tablet (768–1023 px):** no está diseñada; usa el layout de escritorio con barra inferior.
- Revisar a 390, 768 y 1280 px antes de dar una pantalla por terminada.

## Endpoints

### Acceso (Registro, Login, Código, Perfil)

| Endpoint | Para qué |
|---|---|
| `POST /api/auth/register` | `{name, email, acceptedTerms}`; envía el código |
| `POST /api/auth/login` | `{email}`; envía el código |
| `POST /api/auth/code/verify` | `{email, code}`; abre la sesión |
| `POST /api/auth/code/resend` | "Reenviar código" |
| `POST /api/auth/logout` | Perfil: "Cerrar sesión" |
| `GET /api/me` | nombre, correo, iniciales del avatar |

Bienvenida no llama a nada.

### Fuente, Cargando, Análisis, Error

| Endpoint | Para qué |
|---|---|
| `POST /api/statements` | subir PDF (`file`, `billingDate` opcional, `bank` opcional); responde archivo y meses leídos |
| `GET /api/banks` | selector de banco (el diseño no lo tiene) |

- Cargando no necesita endpoint si la subida responde de una vez: es el estado pendiente de la
  petición.
- Error es el 400 de `POST /api/statements` con su llave i18n.
- Análisis lee "seguras" y "por confirmar" de `GET /api/summary`: así `statements` no importa `subscriptions`.
- La fecha de facturación sale del archivo; si el parser no la encuentra, se pide en un campo de la pantalla de subida.

### Revisión

| Endpoint | Para qué |
|---|---|
| `GET /api/summary` | posibles con su evidencia ("Lo vimos en tu archivo: …") |
| `POST /api/subscriptions/decisions` | `{groupKey, decision, cadence?}`: "Sí, es una suscripción" (con cadencia) / "No lo es" |
| `DELETE /api/subscriptions/decisions?groupKey=` | "Deshacer" |

### Resumen, Subs, Detalle

| Endpoint | Para qué |
|---|---|
| `GET /api/summary` | total del mes, reparto por servicio, próximo cobro, próximos cobros, acciones pendientes y, por línea, el detalle (ciclo, últimos cobros, alza) |
| `DELETE /api/subscriptions/line?lineKey=` | "Dejar de seguir esta suscripción" |

El resumen debe entregar:

- **Total del mes en CLP**, con los cobros en USD convertidos.
- **Fecha del próximo cobro por línea.**
- **Cobrados y restantes del mes** ("Llevas 3 cobros este mes y te quedan 5").
- **Alzas pendientes**, con monto de antes y después, desde cuándo e impacto mensual y anual.

Las claves de línea llevan `|`, `:`, `#` y espacios, así que nunca van en el path: viajan en el
cuerpo o como parámetro de consulta codificado.

### Editar y Agregar

| Endpoint | Para qué |
|---|---|
| `PUT /api/subscriptions/aliases` | Editar: nombre (`{lineKey, alias}`; alias vacío lo borra) |
| `PUT /api/subscriptions/overrides` | Editar: monto y día de cobro que reemplazan lo detectado |
| `POST /api/subscriptions/manual` | Agregar: `{name, amount, billingDay}` |
| `PUT` y `DELETE /api/subscriptions/manual/{id}` | editar o quitar una manual |

### Alerta y Cómo cancelar

| Endpoint | Para qué |
|---|---|
| `POST /api/subscriptions/price-alerts/decision` | `{lineKey}`: mantener la suscripción tras el alza (`later` es post-MVP) |
| `GET /api/subscriptions/cancel-guide?lineKey=` | enlace al sitio del servicio (dato del catálogo); los 3 pasos son texto fijo del front |
| `POST /api/subscriptions/cancellations` | `{lineKey}`: el usuario marca la suscripción como "por cancelar" |

### Perfil y Eliminar todo

| Endpoint | Para qué |
|---|---|
| `GET /api/export?format=csv\|json` | "Descargar copia primero"; CSV por defecto |
| `DELETE /api/me` | "Sí, eliminar todo" |

- Apariencia (claro, oscuro, sistema) vive en el navegador, sin endpoint.
- Las filas "Cuentas de correo" y "Notificaciones" de Perfil llevan a pantallas de después del MVP.

## Componentes

### `components/layout`

- `AppShell`: zona superior con brillo + hoja con esquinas de 28 px superpuesta.
- `BrandHeader`: wordmark y pastilla del producto; avatar con sesión, botón de tema sin sesión.
- `BottomNav` (<1024 px) y `SideNav` (≥1024 px): en el MVP solo Resumen y Subs (Calendario y
  Fuentes son de después). `AppShell` elige cuál se ve según el ancho.
- `BackLink`, `StepIndicator` ("Paso 2 de 3 · Análisis").

### `components/ui`

- `Button`: principal, secundario, texto, destructivo.
- `Input`, `FormField` (etiqueta, ayuda, error), `Checkbox` (términos).
- `CodeInput`: 6 dígitos, con `ResendTimer`.
- `Badge`, `Card` (solo hero, alerta y tip), `ListRow`, `ServiceAvatar`.
- `SegmentedRing`: variantes normal, vacía (punteada) y de carga; tamaños 224/176/144/120/104;
  respeta `prefers-reduced-motion`.
- `AnimatedNumber` (de 0 al valor en 0,9 s), `Skeleton`, `Toast`, `SegmentedControl`.
- `EmptyState`: para Revisión vacía y Resumen vacío.

### `features/auth`

- `RegisterForm`, `LoginForm`, `CodeForm`, `WelcomeScreen`.
- Hooks: `useSession`, `useRegister`, `useLogin`, `useVerifyCode`, `useLogout`.
- Protección de rutas con sesión.

### `features/statements`

- `UploadDropzone`, `BankSelect`.
- `AnalysisLoading`, `AnalysisResult`, `UploadError`.
- Hooks: `useUploadStatement`, `useBanks`.

### `features/subscriptions`

- **Resumen:** `SpendHero`, `ServiceLegend`, `NextChargeRow`, `PendingActionsList`,
  `UpcomingChargesList`, `CycleTip`.
- **Revisión:** `ReviewCard` (pendiente, agregada, descartada) y `ReviewProgressRing`. Hay que
  sumar el selector de cadencia: se exige al confirmar y el diseño no la tiene.
- **Subs y Detalle:** `SubscriptionList`, `SubscriptionRow`, `SubscriptionDetail`,
  `PriceChangeBlock`, `ChargeHistory`.
- **Editar y Agregar:** un solo `SubscriptionForm` (nombre, monto, día) con su esquema zod.
- **Alerta y Cancelar:** `PriceAlertCard`, `AlertDecisionActions`, `CancelSteps`,
  `CancellationStatus`.
- Hooks: `useSummary`, `useLine`, `useDecideSubscription`, `useUndoDecision`, `useSaveAlias`,
  `useSaveOverride`, `useAddManualSubscription`, `useDecidePriceAlert`, `useStartCancellation`.

### `features/profile`

- `ProfileHeader`, `SettingsRow`, `AppearanceSelector`, `DeleteAllConfirm`.
- Hooks: `useMe`, `useExportData`, `useDeleteAll`.

### `lib`

- `api/client.ts`, `api/types.ts`.
- `format.ts`: montos CLP y fechas es-CL.
- `i18n/es.ts` + `t()`.
- Tokens de tema en `:root`, tomados del `tokens.json` del handoff.

## Decisiones

Cerradas:

1. **Fuente de datos.** El MVP usa solo cartolas en PDF: de cada línea se guarda servicio,
   monto y fecha, y el texto de la cartola se descarta. El correo va en el post-MVP, por IMAP con
   contraseña de aplicación (solo Gmail). El reenvío a un buzón de la app se descartó por ahora:
   Gmail solo permite crear el filtro desde escritorio. OAuth de Gmail se evalúa con usuarios reales
   (exige verificación y auditoría de Google).
2. **Acceso.** Solo código por correo. Google y Apple quedan para después.
3. **Cancelación.** El usuario marca "por cancelar" a mano. Sin parser de correos de cancelación.
4. **Desconectar una cuenta** (cuando existan cuentas conectadas, después del MVP). Lo ya detectado
   se conserva hasta que el usuario lo elimine.
5. **Total en CLP.** El total del mes se muestra en CLP: los cobros en USD se convierten.
6. **Tipo de cambio** (según cartolas reales de BCI; cada banco lo resuelve en su parser).
   - La cartola internacional trae, por cobro, el monto de origen y el monto en US$ (Visa o
     Mastercard convierte a dólares); no trae el monto en pesos ni el código de la moneda de origen.
   - Si el monto de origen es muchísimo mayor que el de US$ (razón de cientos, constante con
     nombre en `domain`), es un cobro en pesos y se usa el origen tal cual, sin convertir. Así la
     línea coincide con la del estado nacional, donde el mismo servicio cambia de sección entre meses.
   - En cualquier otro caso se usa el monto en US$ (cubre USD, EUR, GBP) convertido a pesos con el
     último tipo de cambio derivado del traspaso: `TRASPASO DEUDA INTERNACIONAL` (CLP, estado
     nacional) ÷ `TRASPASO DEUDA NACIONAL` (US$, estado internacional). Es el tipo de cambio
     vendedor del banco, con margen, o sea lo que el usuario paga de verdad. Se marca "≈".
   - No todos los meses trae traspaso (en uno era US$ 0,00): se guarda el último derivado, con su
     fecha, y se usa el más reciente. Si no existe ninguno, se muestra solo el monto en US$ y el
     total avisa que está incompleto.
   - No hay API externa en el MVP. El dólar observado del Banco Central (sin margen, subestima el
     costo) queda como respaldo para un banco que no entregue nada de lo anterior, cuando llegue.
   - Un cobro en euros u otra moneda aún no se ha visto en una cartola real: su test es sintético
     y se cambia por uno real cuando aparezca.
   - Que se muestre el monto original junto al convertido está cerrado (decisión 19).
7. **Recepción de correos.** No hay buzón en el MVP: se descarta el diseño de recepción (proveedor,
   remitente, límites).
8. **Detección de alza.** Sin correos de aviso, el alza se detecta al llegar el cobro mayor, con la
   cartola siguiente.
9. **Recordatorio "2 días antes del cobro".** Pasa al post-MVP: exige tarea programada y envío de
    correo, y es una predicción desde el próximo cobro.
10. **Manuales y overrides.** Entran al MVP. Una manual es una línea sin cobros que suma al total y
    cuyo próximo cobro es la siguiente vez que ocurre su día de cobro. Un override reemplaza el
    monto o el día detectado en pantalla, total y próximo cobro; los cobros reales no se tocan.
    "Dejar de seguir" oculta la línea (no suma ni aparece) y la marca se guarda por clave de línea
    para que no vuelva al recalcular; para recuperarla se agrega de nuevo.
11. **Sesión y código.** Código de 6 dígitos, vigencia 10 minutos, 5 intentos por código, 3 envíos
    por hora por correo, sesión de 30 días.
13. **Formato de la copia** (`GET /api/export?format=csv|json`). El usuario elige; por defecto CSV.
15. **Inactivas y Otros cobros.** Se calculan, pero no se muestran en el MVP.
16. **Fecha de facturación.** Sale del archivo (`FECHA ESTADO DE CUENTA` y `PERÍODO FACTURADO` en
    BCI; el día de cierre cambia de un mes a otro); si el parser no la encuentra, se pide en un
    campo de la pantalla de subida.
17. **Qué es un alza.** Cualquier aumento entre dos cobros consecutivos de la misma línea y moneda.
    Sin umbral: los Seguros CLC suben ≈ 0,5 % al mes por reajuste y alertarán cada mes; se acepta
    y se revisa con usuarios reales.
18. **Aviso de alza.** Entra al MVP. Sin acceso al correo no se puede avisar antes del cobro: solo se avisa que
    ya subió, en la app, cuando el usuario sube la cartola donde aparece el cobro mayor. Un aviso
    previo (el servicio anuncia el precio nuevo) solo es posible con lectura de correo, en el post-MVP.
19. **Moneda extranjera.** Un cobro en USD se muestra con su monto original y el equivalente en CLP
    (ej. "USD 20 · ≈ $19.000"); el total del mes va en CLP.
20. **Alojamiento.** Todo va dockerizado para que el despliegue sea independiente de la plataforma; el proveedor donde corre queda abierto. Los códigos de acceso se envían con Resend.

Abiertas (se deciden al avanzar; los números se conservan):

12. **Textos legales y qué se guarda.** Términos y política de privacidad (el registro los exige), y
    que lo guardado coincida con lo que promete Análisis: solo servicio, monto y fecha.
14. **Base de datos y secretos** en el contenedor (volumen para SQLite, secretos por variables de entorno); proveedor abierto (decisión 20).

Pendientes que el propio handoff deja abiertos y tocan al MVP:

- Validar el formulario de Editar/Agregar (categoría, moneda distinta de CLP).
- El texto de privacidad ("El resto lo descartamos y no se guarda") solo vale si el proceso
  realmente descarta el resto.
- El tip "Optimización de ciclo" menciona "tu cuenta corriente" y contradice "sin acceso a tu banco".
- Detalle está diseñado con un solo ejemplo (un alza); falta el caso sin alza.

## Post-MVP

- **Cartolas en CSV.** Falta un CSV real de ejemplo (el usuario debe ver cómo obtenerlo del banco);
  sin él no se escribe el parser.
- **Lectura de correo por IMAP** (contraseña de aplicación, solo Gmail): cuentas de correo,
  sincronización, avisos de compra de BCI y conciliación con la cartola. Trae las pantallas Cuentas
  de correo, Añadir cuenta y Fuentes.
- **OAuth de Gmail**, si con usuarios reales el paso de la contraseña de aplicación resulta ser el
  obstáculo (exige verificación y auditoría de Google).
- **Recordatorio "2 días antes del cobro"** por correo (tarea programada).
- Login con Google y Apple, Calendario, Notificaciones y su permiso, Recuperar acceso, Sin conexión.
