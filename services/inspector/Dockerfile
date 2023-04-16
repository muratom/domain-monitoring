# syntax=docker/dockerfile:1

FROM golang:1.20-buster AS builder

ENV USER=appuser
ENV UID=1001

# https://stackoverflow.com/a/55757473
RUN adduser \
    --disabled-password \
    --gecos "" \
    --no-create-home \
    --uid "$UID" \
    "$USER"

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

RUN CGO_ENABLED=0 go build -C cmd/inspector -v -o /usr/local/bin/app
RUN CGO_ENABLED=0 go build -C cmd/healthchecker -v -o /usr/local/bin/healthchecker

FROM gcr.io/distroless/static-debian11:nonroot

WORKDIR /

COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group

COPY --from=builder /usr/local/bin/app /app
COPY --from=builder /usr/local/bin/healthchecker /healthchecker

EXPOSE 8000

# Use an unprivileged user.
USER $USER:$USER

CMD ["/app"]