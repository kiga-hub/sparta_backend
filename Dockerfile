FROM golang:1.20.8-bullseye as builder
COPY . .
RUN bash build.sh sparta_backend /sparta_backend

FROM golang:1.20.8-bullseye
COPY --from=builder /sparta_backend /sparta_backend
COPY sparta_backend.toml /sparta_backend.toml
ENV GODEBUG cgocheck=0
COPY swagger /swagger
WORKDIR /
#HEALTHCHECK --interval=30s --timeout=15s \
#    CMD curl --fail http://localhost:80/health || exit 1
ENTRYPOINT [ "/sparta_backend" ]
CMD ["run"]
