import type { ComponentPropsWithRef, ReactNode } from "react";
import type { LucideIcon } from "lucide-react";
import { cx } from "../cx";
import { appIcons } from "../icons";
import { Icon } from "./Icon";

/** BulkActionBarProps son las propiedades de `BulkActionBar`. */
export interface BulkActionBarProps extends ComponentPropsWithRef<"div"> {
  /** Cuántos elementos hay seleccionados. */
  count: number;
  /**
   * Texto del contador. Por omisión «1 seleccionado» / «N seleccionados»; pásalo para
   * nombrar lo que se seleccionó («3 piezas»).
   */
  countLabel?: string;
  /** Las acciones: varios `BulkAction`. */
  children: ReactNode;
  /** Se pulsó «Quitar selección». Sin esta función no aparece el botón. */
  onClearSelection?: () => void;
  /** Texto del botón para quitar la selección. Por omisión «Quitar selección». */
  clearLabel?: string;
  /** Flotar abajo y al centro de la pantalla en lugar de ir en el flujo de la página. */
  floating?: boolean;
  /** Nombre de la barra para lectores de pantalla. Por omisión «Acciones para la selección». */
  label?: string;
}

/**
 * BulkActionBar es la barra de acciones masivas del sistema Voxa (`.vx-bulkbar`).
 * PROPUESTA: todavía no está en el sistema oficial.
 *
 * Cómo funciona: muestra cuántos elementos hay seleccionados y las acciones para todos a
 * la vez («Mover a otro equipo», «Marcar no apilable», «Eliminar»). La pantalla decide
 * cuándo mostrarla (normalmente con `count > 0`) y qué hace cada acción.
 *
 * Por qué así: va en tinta (el 30 % de la paleta) para no competir con el único botón
 * amarillo de la vista. Las acciones destructivas deben abrir un `ConfirmDialog`.
 */
export function BulkActionBar({
  count,
  countLabel,
  children,
  onClearSelection,
  clearLabel = "Quitar selección",
  floating = false,
  label = "Acciones para la selección",
  className,
  ...rest
}: BulkActionBarProps) {
  const text = countLabel ?? (count === 1 ? "1 seleccionado" : `${count} seleccionados`);

  return (
    <div
      role="toolbar"
      aria-label={label}
      className={cx("vx-bulkbar", floating && "vx-bulkbar--floating", className)}
      {...rest}
    >
      <span className="vx-bulkbar__count" aria-live="polite">
        {text}
      </span>
      <div className="vx-bulkbar__actions">{children}</div>
      {onClearSelection && (
        <button type="button" className="vx-bulkbar__action vx-bulkbar__clear" onClick={onClearSelection}>
          <Icon icon={appIcons.close} size="sm" />
          {clearLabel}
        </button>
      )}
    </div>
  );
}

/** BulkActionProps son las propiedades de `BulkAction`. */
export interface BulkActionProps extends ComponentPropsWithRef<"button"> {
  /** Ícono de 14 px antes del texto; de preferencia de `appIcons`. */
  icon?: LucideIcon;
}

/** BulkAction es un botón dentro de `BulkActionBar`; pasa todo lo demás al `<button>`. */
export function BulkAction({ icon, className, type = "button", children, ...rest }: BulkActionProps) {
  return (
    <button type={type} className={cx("vx-bulkbar__action", className)} {...rest}>
      {icon && <Icon icon={icon} size="sm" />}
      {children}
    </button>
  );
}
