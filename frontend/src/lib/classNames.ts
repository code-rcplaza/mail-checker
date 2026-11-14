/** Joins the class names that apply; false and undefined drop out. */
export function classNames(...names: ReadonlyArray<string | false | undefined>): string {
  return names.filter(Boolean).join(" ");
}
