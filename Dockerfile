# START container setup

FROM golang:latest

RUN apt-get update
RUN apt-get install -y python3 git curl

RUN curl -sL https://deb.nodesource.com/setup_15.x | bash -
RUN apt-get install -y nodejs

RUN npm install yarn -g

# END container setup

ADD . /build
RUN mkdir dist

WORKDIR /build
RUN go mod download && go build -o ../dist/main main.go

WORKDIR /build/web/
RUN mv v1 /dist/web/

WORKDIR /build/web/v2
RUN mkdir /dist/web/v2
RUN yarn install && yarn build && mv build /dist/web/dist

WORKDIR /dist
RUN chmod +x main

EXPOSE 8000

CMD [ "/dist/main" ]




