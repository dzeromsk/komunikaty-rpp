FROM pdf2htmlex/pdf2htmlex:0.18.8.rc1-master-20200630-alpine-3.12.0-x86_64 as pdf2htmlex

FROM golang:alpine as golang

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./
ENV CGO_ENABLED=0
ENV GOBIN=/src
RUN go install github.com/suntong/html2md
RUN go build addindex.go
RUN go build diff.go
RUN go build download.go
RUN go build publish.go

FROM alpine as hugo

WORKDIR /src
RUN apk add --no-cache curl tar
RUN curl -Lo hugo.tar.gz https://github.com/gohugoio/hugo/releases/download/v0.88.1/hugo_extended_0.88.1_Linux-64bit.tar.gz
RUN tar -xvmf hugo.tar.gz hugo

FROM alpine

WORKDIR /github/workspace/

RUN apk add --no-cache \
      tar              \
      libstdc++        \
      libgcc           \
      gnu-libiconv     \
      gettext          \
      glib             \
      freetype         \
      fontconfig       \
      cairo            \
      libpng           \
      libjpeg-turbo    \
      libxml2          \
      make             \
      libc6-compat     \
      git

COPY --from=pdf2htmlex usr/local/bin/pdf2htmlEX /usr/local/bin/
COPY --from=pdf2htmlex usr/local/share/pdf2htmlEX /usr/local/share/pdf2htmlEX
COPY --from=golang src/html2md /usr/local/bin/
COPY --from=golang src/addindex /usr/local/bin/
COPY --from=golang src/diff /usr/local/bin/
COPY --from=golang src/download /usr/local/bin/
COPY --from=golang src/publish /usr/local/bin/
COPY --from=hugo src/hugo /usr/local/bin/

ENV PDF2HTMLEX=/usr/local/bin/pdf2htmlEX \
    HTML2MD=/usr/local/bin/html2md \
    ADDINDEX=/usr/local/bin/addindex \
    DIFF=/usr/local/bin/diff \
    DOWNLOAD=/usr/local/bin/download \
    HUGO=/usr/local/bin/hugo \
    PUBLISH=/usr/local/bin/publish

