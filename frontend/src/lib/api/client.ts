/**
 * The only door to the Go API. It returns the JSON as `unknown`: each gateway
 * checks the shape it expects. Errors come in two kinds: the API answered
 * with an i18n key (`ApiError`), or it could not be reached (`NetworkError`).
 */

export type ErrorParams = Readonly<Record<string, unknown>>;
export type FieldErrors = Readonly<Record<string, string>>;

const internalKey = "common.internal";

export class ApiError extends Error {
  readonly status: number;
  readonly key: string;
  readonly params: ErrorParams;
  readonly fields: FieldErrors;

  constructor(reply: { status: number; key: string; params?: ErrorParams; fields?: FieldErrors }) {
    super(`api ${reply.status}: ${reply.key}`);
    this.name = "ApiError";
    this.status = reply.status;
    this.key = reply.key;
    this.params = reply.params ?? {};
    this.fields = reply.fields ?? {};
  }

  get isUnauthorized(): boolean {
    return this.status === 401;
  }
}

/** The request never got an answer: offline, DNS, server down. Worth one retry. */
export class NetworkError extends Error {
  readonly isRetryable = true;

  constructor(cause: unknown) {
    super("network request failed", { cause });
    this.name = "NetworkError";
  }
}

export type WriteMethod = "POST" | "PUT" | "DELETE";

export interface WriteRequest {
  readonly method: WriteMethod;
  readonly path: string;
  readonly body?: unknown;
}

export interface ApiClient {
  get(path: string): Promise<unknown>;
  send(request: WriteRequest): Promise<unknown>;
}

type FetchFunction = (input: string, init: RequestInit) => Promise<Response>;

function isRecord(value: unknown): value is Record<string, unknown> {
  return typeof value === "object" && value !== null && !Array.isArray(value);
}

function stringsOnly(value: unknown): FieldErrors {
  if (!isRecord(value)) {
    return {};
  }
  return Object.fromEntries(
    Object.entries(value).filter((entry): entry is [string, string] => typeof entry[1] === "string"),
  );
}

/** Reads `{error:{key,params,fields}}`; anything else is common.internal. */
function errorFrom(status: number, body: unknown): ApiError {
  const error = isRecord(body) ? body.error : undefined;
  if (!isRecord(error) || typeof error.key !== "string") {
    return new ApiError({ status, key: internalKey });
  }
  return new ApiError({
    status,
    key: error.key,
    params: isRecord(error.params) ? error.params : {},
    fields: stringsOnly(error.fields),
  });
}

async function readJSON(response: Response): Promise<{ ok: true; body: unknown } | { ok: false }> {
  try {
    const body: unknown = await response.json();
    return { ok: true, body };
  } catch {
    return { ok: false };
  }
}

async function answerOf(response: Response): Promise<unknown> {
  if (response.status === 204) {
    return undefined;
  }
  const read = await readJSON(response);
  if (!read.ok) {
    throw new ApiError({ status: response.status, key: internalKey });
  }
  if (!response.ok) {
    throw errorFrom(response.status, read.body);
  }
  return read.body;
}

/** JSON by default; a FormData goes as is, so the browser writes the multipart boundary. */
function bodyOf(body: unknown): RequestInit {
  if (body === undefined) {
    return {};
  }
  if (body instanceof FormData) {
    return { body };
  }
  return { body: JSON.stringify(body), headers: { "Content-Type": "application/json" } };
}

export function createApiClient(fetchFunction: FetchFunction): ApiClient {
  async function request(path: string, init: RequestInit): Promise<unknown> {
    let response: Response;
    try {
      // Same origin: the session cookie goes along and the browser sets its own Origin.
      response = await fetchFunction(path, { credentials: "same-origin", mode: "same-origin", ...init });
    } catch (cause) {
      throw new NetworkError(cause);
    }
    return answerOf(response);
  }

  return {
    get: (path) => request(path, { method: "GET" }),
    send: ({ method, path, body }) => request(path, { method, ...bodyOf(body) }),
  };
}
