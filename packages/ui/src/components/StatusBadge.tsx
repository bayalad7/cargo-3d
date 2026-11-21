import type { ComponentPropsWithRef } from "react";
import { cx } from "../cx";
import { statusLabels, type Status } from "../status";

/** StatusBadgeProps son las propiedades de `StatusBadge`. */
export interface StatusBadgeProps extends ComponentPropsWithRef<"span"> {
  /** El estado que se muestra. */
  status: Status;
  /**
   * El estado sigue cambiando en tiempo real (por ejemplo, el patio está cargando):
   * el punto late. Se apaga solo si el usuario pidió menos movimiento.
   */
  live?: boolean;
  /**
   * Anunciar a los lectores de pantalla cuando el estado cambie (`role="status"`).
   * Úsalo en el estado principal de la vista, no en cada fila de una lista.
   */
  announce?: boolean;
}

/**
 * StatusBadge es la píldora de estado del sistema Voxa (`.vx-pill`).
 *
 * Cómo funciona: elige la clase `vx-pill--<estado>`; el texto es el de `statusLabels`
 * salvo que la pantalla ponga uno más preciso como hijo («Excede VGM»). Con `live`
 * el punto late; con `announce` el cambio se anuncia en voz.
 *
 * Por qué así: el estado siempre lleva texto y punto, nunca solo color, y ninguno usa el
 * amarillo de la marca. Este componente no sabe de dónde viene el estado (API, WebSocket):
 * solo lo dibuja.
 */
export function StatusBadge({ status, live = false, announce = false, className, children, ...rest }: StatusBadgeProps) {
  return (
    <span
      className={cx("vx-pill", `vx-pill--${status}`, live && "vx-pill--live", className)}
      {...(announce ? { role: "status", "aria-live": "polite" as const } : {})}
      {...rest}
    >
      {children ?? statusLabels[status]}
    </span>
  );
}
