FROM golang:1.22.3-alpine AS compiler
WORKDIR /app
COPY . ./
RUN CGO_ENABLED=0 go build -o /alertmanager-msteams

FROM alpine
LABEL maintaier="karoon.siri@gmail.com"
LABEL description="A Go Web Server that accepts a message from Alertmanager and forwards it to Microsoft Teams Channels using an incoming webhook url."
WORKDIR /app
RUN addgroup -S -g 1000 msteamsforwarder \
    && adduser -S -u 1000 msteamsforwarder -G msteamsforwarder
COPY --from=compiler /alertmanager-msteams ./
COPY . ./
USER msteamsforwarder:msteamsforwarder
EXPOSE 8080
ENTRYPOINT ["/app/alertmanager-msteams"]
CMD []
