FROM golang
ENV CGO_ENABLED=0
ENV GOOS=linux
WORKDIR /app
COPY go.* ./
RUN go mod download
RUN go mod verify
COPY main.go ./
COPY internal/ internal/
COPY pkg/ pkg/
RUN go mod tidy
RUN go build -o /main ./

FROM alpine
WORKDIR /app
COPY --from=0 /main ./
RUN apk update && apk add ffmpeg
CMD [ "/app/main" ]
