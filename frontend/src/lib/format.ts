/** Amounts travel in cents (like `Money` in the back) with their currency. */
export interface MoneyAmount {
  readonly amount: number;
  readonly currency: string;
}

/** "$84.060" for pesos, "US$19,99" for dollars: the way people in Chile read them. */
export function formatMoney({ amount, currency }: MoneyAmount): string {
  return new Intl.NumberFormat("es-CL", { style: "currency", currency }).format(amount / 100);
}

const monthPattern = /^(\d{4})-(\d{2})$/;
const monthFormat = new Intl.DateTimeFormat("es-CL", { month: "long", year: "numeric", timeZone: "UTC" });

/** "2026-09" → "septiembre de 2026"; anything else stays as it came. */
export function formatMonth(month: string): string {
  const [, year, number] = monthPattern.exec(month) ?? [];
  if (year === undefined || number === undefined) {
    return month;
  }
  return monthFormat.format(new Date(Date.UTC(Number(year), Number(number) - 1, 1)));
}

const listFormat = new Intl.ListFormat("es", { style: "long", type: "conjunction" });

export function formatMonths(months: readonly string[]): string {
  return listFormat.format(months.map(formatMonth));
}
