FROM golang:1.20.8-bullseye as builder
COPY . .
RUN bash build.sh websocket /websocket

FROM golang:1.20.8-bullseye
COPY --from=builder /websocket /websocket
COPY websocket.toml /websocket.toml
ENV GODEBUG cgocheck=0
COPY swagger /swagger
WORKDIR /
#HEALTHCHECK --interval=30s --timeout=15s \
#    CMD curl --fail http://localhost:80/health || exit 1
ENTRYPOINT [ "/websocket" ]
CMD ["run"]
