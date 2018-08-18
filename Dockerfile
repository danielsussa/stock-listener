FROM golang
RUN go get github.com/githubnemo/CompileDaemon
ENV PATH /scripts:$PATH