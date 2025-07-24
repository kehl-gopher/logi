FROM golang:1.24.4-bullseye AS build

WORKDIR /app/

COPY ../go.mod ../go.sum ./

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go mod download


FROM build AS dev

RUN go install github.com/air-verse/air@latest 

COPY ../ .

CMD ["air", "-c", ".air.toml"]

FROM build AS prod

RUN useradd -u 1001 logi


RUN go build \
    -ldflags="-w -s -linkmode external -extldflags -static" \
    -tags netgo \
    -o main

FROM scratch AS final-prod

ENV GIN_MODE=release

WORKDIR /app/

COPY --from=prod /etc/passwd /etc/passwd
COPY --from=prod /app/main main

USER logi

EXPOSE 8080

CMD ["/main"]