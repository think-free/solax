FROM golang:alpine AS go-env
ADD . /src
RUN apk add git
WORKDIR /src/
RUN go build .

FROM alpine
WORKDIR /app
COPY --from=go-env /src/solax /app/
EXPOSE 80
ENTRYPOINT ./solax
