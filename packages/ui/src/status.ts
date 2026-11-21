/**
 * Status es el estado de una carga, un plan o un equipo, tal como lo define el sistema Voxa.
 *
 * Cómo funciona: cada valor corresponde a un grupo de tokens (`--status-<estado>`,
 * `-soft` e `-ink`) y a una clase `vx-pill--<estado>`.
 *
 * Por qué así: los estados nunca usan el amarillo de la marca, y siempre se muestran con
 * texto, punto o ícono además del color, para que se distingan sin ver colores.
 */
export type Status = "success" | "progress" | "warning" | "danger" | "delayed";

/**
 * statusLabels es el texto por omisión de cada estado, en español de México.
 * Las pantallas pueden dar un texto más preciso («Excede VGM», «CG fuera de rango»).
 */
export const statusLabels: Readonly<Record<Status, string>> = {
  success: "Listo",
  progress: "Cargando",
  warning: "Cerca del límite",
  danger: "Excede",
  delayed: "Demorado",
};
