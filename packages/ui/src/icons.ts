import {
  Anchor,
  Bell,
  Box,
  Boxes,
  Building,
  CircleCheck,
  Clock,
  Container,
  Copy,
  Crosshair,
  Download,
  FileSpreadsheet,
  Forklift,
  Gauge,
  GlassWater,
  Info,
  Layers,
  LayoutGrid,
  ListOrdered,
  Loader,
  MapPin,
  Moon,
  Move3d,
  OctagonExclamationPoint,
  Package,
  Pencil,
  Plus,
  Printer,
  Rotate3d,
  RotateCw,
  Ruler,
  Search,
  Send,
  Settings,
  Sun,
  Trash,
  TriangleExclamationPoint,
  Truck,
  Undo2,
  Upload,
  Users,
  Weight,
  X,
  Zap,
  type LucideIcon,
} from "lucide-react";
import type { Status } from "./status";

/**
 * appIcons es el vocabulario de íconos de Voxa: un ícono de Lucide por cada concepto del
 * dominio, para que todas las apps (web, escritorio y móvil) usen el mismo dibujo para lo mismo.
 *
 * Cómo funciona: se usa con el componente `Icon`, por ejemplo
 * `<Icon icon={appIcons.weight} />`.
 *
 * Por qué así: si cada pantalla escoge su propio ícono, «peso» termina con tres dibujos
 * distintos. Un concepto nuevo se agrega aquí, no en la pantalla.
 */
export const appIcons = {
  // Navegación
  plans: LayoutGrid, // Planes de carga
  settings: Settings,
  notifications: Bell,
  search: Search,
  customers: Users,
  company: Building,

  // Equipos y lugares
  container: Container, // Contenedor marítimo
  dryVan: Truck, // Caja seca
  yard: Forklift, // Patio
  port: Anchor, // Puerto
  destination: MapPin,

  // Carga y restricciones
  piece: Package,
  pieces: Boxes,
  weight: Weight, // Peso, VGM
  axleLoad: Gauge, // Peso por eje
  dimensions: Ruler,
  volume: Box,
  centerOfGravity: Crosshair,
  stacking: Layers, // Apilabilidad
  unloadOrder: ListOrdered, // Orden de descarga (LIFO por destino)
  fragile: GlassWater,

  // Visor 3D
  view3d: Rotate3d,
  move: Move3d,
  rotate: RotateCw,

  // Acciones
  optimize: Zap, // «Optimizar carga»
  releasePlan: Send, // «Liberar plan»
  add: Plus,
  edit: Pencil,
  duplicate: Copy,
  delete: Trash,
  undo: Undo2,
  importFile: FileSpreadsheet,
  upload: Upload,
  download: Download,
  print: Printer,
  close: X,

  // Tema
  themeLight: Sun,
  themeDark: Moon,
} as const satisfies Record<string, LucideIcon>;

/** statusIcons es el ícono de cada estado, para toasts, modales y listas. */
export const statusIcons: Readonly<Record<Status, LucideIcon>> = {
  success: CircleCheck,
  progress: Loader,
  warning: TriangleExclamationPoint,
  danger: OctagonExclamationPoint,
  delayed: Clock,
};

/** infoIcon es el ícono de un aviso neutro (sin estado). */
export const infoIcon: LucideIcon = Info;

export type { LucideIcon };
