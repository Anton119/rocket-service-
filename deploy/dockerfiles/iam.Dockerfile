# iam.Dockerfile — gRPC-сервис (IAMService).
# Multi-stage сборка: тяжёлый Go-компилятор только на этапе build,
# в финальный образ попадает только скомпилированный бинарник.

FROM golang:1.26-alpine AS builder

WORKDIR /app

# Копируем файлы Go workspace и все go.mod/go.sum.
# Отдельным слоем — Docker кеширует его, пока зависимости не изменятся.
COPY go.work go.work.sum ./
COPY platform/go.mod platform/go.sum ./platform/
COPY shared/go.mod shared/go.sum ./shared/
COPY order/go.mod order/go.sum ./order/
COPY inventory/go.mod inventory/go.sum ./inventory/
COPY payment/go.mod payment/go.sum ./payment/
COPY iam/go.mod iam/go.sum ./iam/
COPY assembly/go.mod assembly/go.sum ./assembly/

RUN go work sync

# Копируем исходный код.
COPY . .

# Собираем бинарник. CGO_ENABLED=0 — статическая линковка,
# не нужны libc и другие системные библиотеки в runtime-образе.
RUN CGO_ENABLED=0 go build -o /bin/service ./iam/cmd/main.go

# --- Runtime ---
# alpine — минимальный образ (~5 MB), достаточный для запуска Go-бинарника.
FROM alpine:3.23

COPY --from=builder /bin/service /bin/service

ENTRYPOINT ["/bin/service"]
