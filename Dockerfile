# Global
ARG port=8000
ARG binary=editfolio-api-v1

# Builder
FROM golang:1.17.1-alpine3.13 as builder

ARG binary

LABEL PO=HellP
LABEL PM=HellP
LABEL Maintainers=HellP
LABEL Maintainers_Mail=jlee@stockfolio.ai

WORKDIR /app

RUN apk update && apk upgrade && \
    apk --update add git make gcc g++

COPY . .

ENV BINARY=$binary

RUN make get-dependencies
# RUN make -C app/ unittest
RUN make build

# Distribution
FROM alpine:3.13 as distribution

ARG port
ARG binary

WORKDIR /app

RUN apk update && apk upgrade && \
    apk --update add bash ca-certificates

ENV PORT=$port
ENV BINARY=$binary

COPY --from=builder /core/app/${BINARY} /app

EXPOSE ${PORT}

CMD /app/${BINARY}