# Этап 1: Сборка агента
FROM golang:latest AS builder

WORKDIR /app
COPY . /app
COPY ./examples/hacker/.env /app/.env

RUN go build -o agent examples/hacker/main.go

# Этап 2: Финальный образ
FROM kalilinux/kali-rolling

# Устанавливаем зависимости (если нужно)
RUN --mount=target=/var/lib/apt/lists,type=cache,sharing=locked \
    --mount=target=/var/cache/apt,type=cache,sharing=locked \
    rm -f /etc/apt/apt.conf.d/docker-clean \
    && apt-get update \
    && apt-get -y --no-install-recommends install \
          git curl wget metasploit-framework \
          nmap net-tools 

RUN ifconfig

# Копируем исполняемый файл из этапа сборки
COPY --from=builder /app/agent /app/agent
COPY --from=builder /app/.env /app/.env
WORKDIR /app
CMD ["./agent"]
