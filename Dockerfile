FROM golang:alpine
ADD . /src
RUN apk add --no-cache git && \
cd /src && go-wrapper download

RUN cd /src && go build cli.go && \
go build executor.go && mv cli /cli && \
mv executor /executor && cd / && \
rm -rf /src /go
WORKDIR "/"
ENTRYPOINT ["/cli"]