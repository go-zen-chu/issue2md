FROM gcr.io/distroless/static-debian11:latest

COPY issue2md /bin/issue2md
ENTRYPOINT ["/bin/issue2md"]
