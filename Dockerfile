FROM golang
ADD . /src
RUN \
cd /src && go-wrapper download

RUN cd /src && go build cli.go && \
go build executor.go && mv cli /cli && \
mv executor /executor && cd / && rm -rf src

ENTRYPOINT ["/cli"]