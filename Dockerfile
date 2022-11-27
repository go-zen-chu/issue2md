FROM golang:1.19.3-buster AS gobuilder

# use static link build
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

WORKDIR /usr/local/src/repo
COPY . /usr/local/src/repo
RUN go build ./cmd/issue2md

FROM gcr.io/distroless/static-debian11:latest

COPY --from=gobuilder /usr/local/src/repo/issue2md /bin/issue2md
ENTRYPOINT ["/bin/issue2md"]
