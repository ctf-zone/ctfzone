FROM node:latest AS ui
WORKDIR /usr/src/
COPY web .
ENV NODE_ENV production

WORKDIR /usr/src/public
RUN npm --production=false install
RUN npm run build

WORKDIR /usr/src/admin
RUN npm --production=false install
RUN npm run build


FROM golang:latest AS api
WORKDIR /go/src/github.com/ctf-zone/ctfzone
COPY . .
ENV GO111MODULE on
ENV GOOS linux
ENV CGO_ENABLED 0
RUN make


FROM alpine:latest
RUN addgroup -S ctfzone && adduser -S -u 1001 -G ctfzone ctfzone
RUN apk add --no-cache ca-certificates
USER ctfzone
WORKDIR /home/ctfzone
RUN mkdir -p files static/public static/admin
COPY templates ./templates
COPY --from=ui --chown=ctfzone:ctfzone /usr/src/public/dist static/public/
COPY --from=ui --chown=ctfzone:ctfzone /usr/src/admin/dist static/admin/
COPY --from=api --chown=ctfzone:ctfzone /go/src/github.com/ctf-zone/ctfzone/ctfzone .
EXPOSE 8080 8443
CMD ["./ctfzone"]
