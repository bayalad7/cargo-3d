import type { ComponentPropsWithRef } from "react";
import { cx } from "../cx";

/** FilterChipProps son las propiedades de `FilterChip`. */
export interface FilterChipProps extends Omit<ComponentPropsWithRef<"button">, "aria-pressed"> {
  /** El filtro está aplicado. */
  active?: boolean;
}

/**
 * FilterChip es un filtro de lista (`.vx-chip`): «Todas», «Cargando», «Con problema».
 *
 * Cómo funciona: es un `<button>` con `aria-pressed`; el activo va en amarillo con texto
 * en tinta. Qué pasa al pulsarlo lo decide la pantalla con `onClick`.
 *
 * Por qué así: `aria-pressed` le dice al lector de pantalla cuál está aplicado, y la
 * clase `vx-chip--active` sale de esa misma propiedad para que nunca se contradigan.
 */
export function FilterChip({ active = false, className, type = "button", ...rest }: FilterChipProps) {
  return (
    <button
      type={type}
      aria-pressed={active}
      className={cx("vx-chip", active && "vx-chip--active", className)}
      {...rest}
    />
  );
}
