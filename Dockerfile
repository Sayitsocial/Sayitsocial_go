# START container setup

FROM golang:1.15.5-alpine3.12

RUN apk add --update npm && npm install yarn -g

# END container setup
WORKDIR /build
ADD go.mod go.sum ./
RUN mkdir ../dist
ADD pkg/database/migrations /dist/pkg/database/migrations
# RUN mv /build/web/v1 /dist/web/

WORKDIR /tmp/node
ADD web/v2/package.json web/v2/yarn.lock ./
RUN yarn install --frozen-lockfile --ignore-platform --ignore-engines --quiet

WORKDIR /build/web
ADD web/ .
RUN mv /tmp/node/* ./v2
RUN cd v2 && yarn build

WORKDIR /build
ADD pkg ./pkg
ADD swagger /dist/swagger
ADD main.go .
RUN go build -o ../dist/main main.go

WORKDIR /dist
RUN mkdir -p web/v2 && mv ../build/web/v2/build web/v2/dist
RUN chmod +x main

EXPOSE 8000

ADD https://github.com/ufoscout/docker-compose-wait/releases/download/2.7.3/wait /wait
RUN chmod +x /wait

# RUN apk add --update bash && rm -rf /var/cache/apk/*
ENTRYPOINT ["/bin/sh", "-c", "/wait && /dist/main"]
