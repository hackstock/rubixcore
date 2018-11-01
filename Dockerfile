FROM golang:1.11
ENV NAME=rubixcore 
ENV APP_DIR=/${NAME}
ENV GOOS=linux
ENV GO_LINKER_FLAGS=-ldflags="-s -w"
COPY . ${APP_DIR}
WORKDIR ${APP_DIR}
RUN go build -o ${NAME}  -ldflags="-s -w" ${APP_DIR}/cmd/rubixcore/main.go 


ENTRYPOINT [ "./rubixcore" ]

