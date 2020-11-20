# START container setup
ARG ARCH=
FROM --platform=$BUILDPLATFORM golang:1.15.5-alpine3.12 AS build

RUN apk add --update --no-cache python2 npm make cmake && npm install yarn -g

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

FROM alpine
COPY --from=build /dist /dist

EXPOSE 8000

ENTRYPOINT ["/bin/sh", "-c", "/dist/main"]
