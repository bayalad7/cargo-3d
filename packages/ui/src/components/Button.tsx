import type { ComponentPropsWithRef } from "react";
import type { LucideIcon } from "lucide-react";
import { cx } from "../cx";
import { Icon } from "./Icon";

/**
 * ButtonVariant es el tipo de botón.
 *
 * - `primary`: amarillo con texto en tinta. **Uno solo por vista** («Optimizar carga»).
 * - `secondary`: tinta con texto `on-ink`. Acción fuerte pero no principal («Liberar plan»).
 * - `ghost`: tarjeta con borde. Acciones de apoyo («Ver en 3D», «Cancelar»).
 * - `danger` (propuesta): rojo, solo para confirmar algo destructivo dentro de un modal.
 */
export type ButtonVariant = "primary" | "secondary" | "ghost" | "danger";

/** ButtonProps son las propiedades de `Button`: las de un `<button>` más las del sistema. */
export interface ButtonProps extends ComponentPropsWithRef<"button"> {
  /** Por omisión `ghost`: lo seguro es no competir con el botón principal. */
  variant?: ButtonVariant;
  /** `sm` (28 px) dentro de tarjetas y tablas; `md` (36 px) en el resto. */
  size?: "md" | "sm";
  /** Ícono de 14 px antes del texto. */
  icon?: LucideIcon;
}

/**
 * Button es el botón de acción del sistema Voxa (`.vx-btn`).
 *
 * Cómo funciona: solo arma las clases; todo lo demás (`onClick`, `disabled`, `type`,
 * `ref`…) pasa tal cual al `<button>`. Es `type="button"` por omisión para no enviar
 * formularios por accidente.
 *
 * Por qué así: el texto lo pone la pantalla, como verbo + objeto en español
 * («Optimizar carga»). Deshabilitado se marca con `disabled`, nunca solo con color.
 */
export function Button({
  variant = "ghost",
  size = "md",
  icon,
  className,
  type = "button",
  children,
  ...rest
}: ButtonProps) {
  return (
    <button
      type={type}
      className={cx("vx-btn", `vx-btn--${variant}`, size === "sm" && "vx-btn--sm", className)}
      {...rest}
    >
      {icon && <Icon icon={icon} size="sm" />}
      {children}
    </button>
  );
}
