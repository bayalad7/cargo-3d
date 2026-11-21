/**
 * cx une nombres de clase CSS e ignora los vacíos (`false`, `null`, `undefined`, `""`).
 *
 * @example cx("vx-btn", primary && "vx-btn--primary", className)
 */
export function cx(...classes: Array<string | false | null | undefined>): string {
  return classes.filter(Boolean).join(" ");
}
