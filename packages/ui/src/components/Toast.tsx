import type { ComponentPropsWithRef, ReactNode } from "react";
import type { LucideIcon } from "lucide-react";
import { cx } from "../cx";
import { appIcons, infoIcon, statusIcons } from "../icons";
import type { Status } from "../status";
import { Icon } from "./Icon";

/** ToastTone es el tono de una notificación: un estado o `info` (neutro). */
export type ToastTone = Status | "info";

/** ToastProps son las propiedades de `Toast`. */
export interface ToastProps extends Omit<ComponentPropsWithRef<"div">, "title"> {
  /** Por omisión `info`. */
  tone?: ToastTone;
  /** Qué pasó, corto: «Plan liberado». */
  title: string;
  /** Detalle opcional: «MZO-0142 · 312 cajas · 24 310 kg». */
  message?: ReactNode;
  /** Ícono; por omisión el del tono. */
  icon?: LucideIcon;
  /** Botones de la notificación (por ejemplo «Deshacer»), de tamaño `sm`. */
  actions?: ReactNode;
  /** Se pulsó cerrar. Sin esta función no hay botón de cerrar. */
  onDismiss?: () => void;
  /** Texto del botón de cerrar para lectores de pantalla. Por omisión «Cerrar aviso». */
  dismissLabel?: string;
}

/**
 * Toast es una notificación breve del sistema Voxa (`.vx-toast`).
 *
 * Cómo funciona: tarjeta con borde izquierdo e ícono del color del estado; el texto va
 * siempre en tinta. Solo se dibuja: cuándo aparece, cuánto dura y la cola de avisos las
 * maneja la app. Los de tono `danger` llevan `role="alert"` para anunciarse de inmediato.
 *
 * Por qué así: el color del estado va en el borde y el ícono, no en el texto, para que
 * todo se lea igual en claro y oscuro.
 */
export function Toast({
  tone = "info",
  title,
  message,
  icon,
  actions,
  onDismiss,
  dismissLabel = "Cerrar aviso",
  className,
  ...rest
}: ToastProps) {
  const glyph = icon ?? (tone === "info" ? infoIcon : statusIcons[tone]);

  return (
    <div
      className={cx("vx-toast", tone !== "info" && `vx-toast--${tone}`, className)}
      {...(tone === "danger" ? { role: "alert" } : {})}
      {...rest}
    >
      <span className="vx-toast__icon">
        <Icon icon={glyph} size="lg" />
      </span>
      <div className="vx-toast__content">
        <p className="vx-toast__title">{title}</p>
        {message && <div className="vx-toast__message">{message}</div>}
        {actions && <div className="vx-toast__actions">{actions}</div>}
      </div>
      {onDismiss && (
        <button type="button" className="vx-toast__close" onClick={onDismiss} aria-label={dismissLabel}>
          <Icon icon={appIcons.close} size="md" />
        </button>
      )}
    </div>
  );
}

/** ToasterProps son las propiedades de `Toaster`. */
export interface ToasterProps extends ComponentPropsWithRef<"section"> {
  /** Nombre de la región. Por omisión «Notificaciones». */
  label?: string;
}

/**
 * Toaster es la región donde se apilan los `Toast`: abajo a la derecha en escritorio y a
 * lo ancho (con 16 px de margen) en teléfono.
 *
 * Cómo funciona: es una región `aria-live="polite"`, así que cada aviso nuevo se lee en
 * voz sin interrumpir. Debe estar en la página desde el inicio (aunque vacía) para que
 * los lectores de pantalla la detecten.
 */
export function Toaster({ label = "Notificaciones", className, children, ...rest }: ToasterProps) {
  return (
    <section
      className={cx("vx-toaster", className)}
      aria-label={label}
      aria-live="polite"
      aria-relevant="additions"
      {...rest}
    >
      {children}
    </section>
  );
}
