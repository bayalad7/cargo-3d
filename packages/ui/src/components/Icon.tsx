import type { LucideIcon, LucideProps } from "lucide-react";

/** IconSize son los tres tamaños de ícono del sistema. */
export type IconSize = "sm" | "md" | "lg";

/** iconSizes traduce cada tamaño a píxeles: 14 (botones), 16 (texto) y 18 (indicadores). */
export const iconSizes: Readonly<Record<IconSize, number>> = { sm: 14, md: 16, lg: 18 };

/** IconProps son las propiedades de `Icon`. */
export interface IconProps extends Omit<LucideProps, "size" | "ref"> {
  /** El ícono de Lucide; de preferencia uno de `appIcons`. */
  icon: LucideIcon;
  /** Tamaño del sistema. Por omisión `md` (16 px). */
  size?: IconSize;
  /**
   * Texto para lectores de pantalla. Solo cuando el ícono va solo (sin texto al lado);
   * si no se da, el ícono es decorativo y se oculta.
   */
  label?: string;
}

/**
 * Icon dibuja un ícono de Lucide con las reglas de Voxa.
 *
 * Cómo funciona: fija el trazo en 1.5 px y el tamaño en 14, 16 o 18 px; el color es
 * `currentColor`, así que toma el del texto (por regla, `ink`).
 *
 * Por qué así: Lucide trae trazo de 2 px y cualquier tamaño; el sistema pide trazos finos
 * y redondeados en pocos tamaños para que todo se vea de la misma familia.
 */
export function Icon({ icon: Glyph, size = "md", label, strokeWidth = 1.5, ...rest }: IconProps) {
  return (
    <Glyph
      size={iconSizes[size]}
      strokeWidth={strokeWidth}
      focusable="false"
      {...(label ? { role: "img", "aria-label": label } : { "aria-hidden": true })}
      {...rest}
    />
  );
}
