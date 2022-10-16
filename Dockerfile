FROM scratch

COPY issue2md /bin/issue2md
ENTRYPOINT ["/bin/issue2md"]
