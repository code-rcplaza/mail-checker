import { classNames } from "../../../lib/classNames";

const seriesClasses = [
  "bg-series-1",
  "bg-series-2",
  "bg-series-3",
  "bg-series-4",
  "bg-series-5",
  "bg-series-6",
] as const;

const neutralSymbol = "•";

/** The first letter of the name, or a neutral symbol when it has none. */
export function initialOf(name: string): string {
  const letter = name.trim().match(/\p{L}|\p{N}/u)?.[0];
  return letter === undefined ? neutralSymbol : letter.toLocaleUpperCase("es-CL");
}

/** The same name always gets the same series color. */
export function seriesClassOf(name: string): string {
  let hash = 0;
  for (const char of name.trim().toLocaleLowerCase("es-CL")) {
    hash = (hash * 31 + (char.codePointAt(0) ?? 0)) >>> 0;
  }
  return seriesClasses[hash % seriesClasses.length] ?? "bg-series-6";
}

export type ServiceAvatarSize = 32 | 44;

/** No service logos (brand rights): the initial on a stable color. */
export function ServiceAvatar({
  name,
  size = 44,
}: {
  readonly name: string;
  readonly size?: ServiceAvatarSize;
}) {
  return (
    <span
      aria-hidden="true"
      className={classNames(
        "inline-flex flex-none items-center justify-center rounded-full font-semibold text-on-series",
        size === 44 ? "size-11 text-base" : "size-8 text-sm",
        seriesClassOf(name),
      )}
    >
      {initialOf(name)}
    </span>
  );
}
