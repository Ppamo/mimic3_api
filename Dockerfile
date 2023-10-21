FROM golang:1.20-buster AS builder
ARG APP=mimic3-api
ENV APP=${APP}
WORKDIR /go/src/ppamo/${APP}
COPY src /go/src/ppamo/${APP}
RUN go mod tidy &&  CGO_ENABLED=0 go build -v -ldflags="-s -w -extldflags '-static'" -a -o /bin/${APP} main.go

FROM mycroftai/mimic3:0.2.4
ARG APP=mimic3-api
ENV APP=${APP}
USER 0
RUN mkdir -p /var/cache/apt/amd64/archives/partial /var/cache/apt/arm64/archives/partial
RUN apt update && apt -y upgrade && apt install -y ffmpeg
RUN mkdir -p /opt/mimic3-server/voices && /home/mimic3/app/.venv/bin/mimic3-download --output-dir "/opt/mimic3-server/voices"  "en_US/*" "es_ES/*"
COPY --from=builder /bin/${APP} /opt/mimic3-server/${APP}
COPY ./effects/* /opt/mimic3-server/effects/
COPY ./src/config.json /opt/mimic3-server/config.json
ENV CONFIG_PATH=/opt/mimic3-server/config.json
ENTRYPOINT /opt/mimic3-server/$APP
