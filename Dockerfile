FROM ubuntu:latest

COPY issue2md /bin/issue2md
CMD ["/bin/issue2md"]
