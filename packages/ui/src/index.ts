/**
 * @cargo-3d/ui — Sistema de diseño Voxa (ADR-0014).
 *
 * Componentes de presentación: reciben datos por propiedades y avisan con funciones
 * (`onConfirm`, `onDismiss`…). No guardan estado de negocio ni llaman a la API.
 * Los estilos se importan aparte: `@import "@cargo-3d/ui/styles.css";`.
 *
 * Documentación: docs/src/ux/sistema-diseno.md
 */
export { Logo, type LogoProps, type LogoVariant } from "./brand/Logo";
export { Icon, iconSizes, type IconProps, type IconSize } from "./components/Icon";
export { Button, type ButtonProps, type ButtonVariant } from "./components/Button";
export { StatusBadge, type StatusBadgeProps } from "./components/StatusBadge";
export { FilterChip, type FilterChipProps } from "./components/FilterChip";
export { ConfirmDialog, type ConfirmDialogProps, type ConfirmTone } from "./components/ConfirmDialog";
export { Toast, Toaster, type ToastProps, type ToasterProps, type ToastTone } from "./components/Toast";
export {
  BulkAction,
  BulkActionBar,
  type BulkActionBarProps,
  type BulkActionProps,
} from "./components/BulkActionBar";
export { appIcons, infoIcon, statusIcons, type LucideIcon } from "./icons";
export { statusLabels, type Status } from "./status";
export { cx } from "./cx";
