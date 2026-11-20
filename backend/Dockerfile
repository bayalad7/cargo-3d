# Imagen del backend (API + motor) en dos etapas.
# Ver docs/src/adr/0010-docker.md
#
# Etapa 1 — "build": una imagen con Go completo para compilar.
# Etapa 2 — "runtime": solo el binario compilado, sin Go, sin shell, sin paquetes.
# Resultado: una imagen de pocos megabytes con muy poca superficie de ataque.

# ---------- Etapa 1: compilar ----------
# Fijamos la versión de Go para que la compilación sea reproducible.
FROM golang:1.27 AS build

# Carpeta de trabajo dentro del contenedor.
WORKDIR /src

# Copiamos primero solo go.mod (y go.sum cuando exista) para que Docker guarde en
# caché la descarga de dependencias: si solo cambia el código, no se vuelve a descargar.
COPY go.mod ./
RUN go mod download

# Ahora sí, todo el código del backend.
COPY . .

# Compilamos un binario estático:
#   CGO_ENABLED=0  → sin dependencias de C, corre en una imagen vacía.
#   GOWORK=off     → dentro del contenedor solo existe este módulo, no el workspace.
#   -trimpath      → quita rutas locales del binario.
#   -ldflags="-s -w" → quita información de depuración (binario más chico).
RUN CGO_ENABLED=0 GOWORK=off go build -trimpath -ldflags="-s -w" -o /out/api ./cmd/api

# ---------- Etapa 2: ejecutar ----------
# distroless/static: imagen mínima de Google, sin shell ni gestor de paquetes.
# La variante "nonroot" corre con un usuario sin privilegios.
FROM gcr.io/distroless/static-debian12:nonroot AS runtime

# Copiamos solo el binario desde la etapa de compilación.
COPY --from=build /out/api /api

# Puerto en el que escucha la API (se puede cambiar con la variable PORT).
ENV PORT=8080
EXPOSE 8080

# Usuario sin privilegios (definido por la imagen distroless).
USER nonroot:nonroot

# Comando que arranca el contenedor.
ENTRYPOINT ["/api"]
