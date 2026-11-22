import { useEffect, useId, useRef, type ReactNode } from "react";
import type { LucideIcon } from "lucide-react";
import { cx } from "../cx";
import { statusIcons } from "../icons";
import { Button } from "./Button";
import { Icon } from "./Icon";

/** ConfirmTone es qué tan delicada es la acción que se confirma. */
export type ConfirmTone = "default" | "warning" | "danger";

/** ConfirmDialogProps son las propiedades de `ConfirmDialog`. */
export interface ConfirmDialogProps {
  /** El modal está abierto. Lo controla la pantalla. */
  open: boolean;
  /** Pregunta corta: «¿Eliminar el plan MZO-0142?». */
  title: string;
  /** Qué pasará y si se puede deshacer. */
  description?: ReactNode;
  /** Verbo + objeto: «Eliminar plan». */
  confirmLabel: string;
  /** Por omisión «Cancelar». */
  cancelLabel?: string;
  /** `danger` pinta el botón en rojo y pone el foco en «Cancelar». Por omisión `default`. */
  tone?: ConfirmTone;
  /** Ícono del encabezado; por omisión el del tono (ninguno en `default`). */
  icon?: LucideIcon;
  /** La acción está en curso: el botón de confirmar se deshabilita. */
  busy?: boolean;
  /** Se pulsó confirmar. */
  onConfirm: () => void;
  /** Se pulsó cancelar, Esc o se cerró el modal. */
  onCancel: () => void;
}

/**
 * ConfirmDialog es el modal de confirmación del sistema Voxa (`.vx-dialog`).
 *
 * Cómo funciona: usa el `<dialog>` nativo del navegador, que ya atrapa el foco, cierra con
 * Esc y dibuja el velo. El componente solo abre o cierra el `<dialog>` según `open` y avisa
 * con `onConfirm` u `onCancel`; qué se hace al confirmar lo decide la pantalla.
 *
 * Por qué así: un modal hecho a mano suele olvidar el foco o el Esc. En acciones
 * destructivas el foco empieza en «Cancelar», para que un Enter distraído no borre nada.
 */
export function ConfirmDialog({
  open,
  title,
  description,
  confirmLabel,
  cancelLabel = "Cancelar",
  tone = "default",
  icon,
  busy = false,
  onConfirm,
  onCancel,
}: ConfirmDialogProps) {
  const ref = useRef<HTMLDialogElement>(null);
  const titleId = useId();
  const bodyId = useId();
  const headIcon = icon ?? (tone === "default" ? undefined : statusIcons[tone]);

  // Sincroniza el <dialog> nativo con la propiedad `open`.
  useEffect(() => {
    const dialog = ref.current;
    if (!dialog) return;
    if (open && !dialog.open) dialog.showModal();
    if (!open && dialog.open) dialog.close();
  }, [open]);

  return (
    <dialog
      ref={ref}
      className="vx-dialog"
      aria-labelledby={titleId}
      {...(description ? { "aria-describedby": bodyId } : {})}
      onCancel={(event) => {
        // Esc: la pantalla decide si se cierra, para que `open` siga siendo la verdad.
        event.preventDefault();
        onCancel();
      }}
    >
      <div className="vx-dialog__head">
        {headIcon && (
          <span className={cx("vx-dialog__icon", tone !== "default" && `vx-dialog__icon--${tone}`)}>
            <Icon icon={headIcon} size="lg" />
          </span>
        )}
        <div>
          <h2 id={titleId} className="vx-dialog__title">
            {title}
          </h2>
          {description && (
            <div id={bodyId} className="vx-dialog__body">
              {description}
            </div>
          )}
        </div>
      </div>
      <div className="vx-dialog__actions">
        <Button variant="ghost" onClick={onCancel} autoFocus={tone === "danger"}>
          {cancelLabel}
        </Button>
        <Button
          variant={tone === "danger" ? "danger" : "primary"}
          onClick={onConfirm}
          disabled={busy}
          aria-busy={busy}
          autoFocus={tone !== "danger"}
        >
          {confirmLabel}
        </Button>
      </div>
    </dialog>
  );
}
