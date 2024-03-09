FROM golang:1.21.6-alpine AS builder

COPY . /github.com/Tim-Sa/JobResume/source/
WORKDIR /github.com/Tim-Sa/JobResume/source/

RUN go mod download
RUN go build -o ./bin/resume cmd/main.go

FROM alpine:3.19.1

WORKDIR /root/

COPY --from=builder /github.com/Tim-Sa/JobResume/source/bin/resume .
COPY --from=builder /github.com/Tim-Sa/JobResume/source/templates ./templates
COPY --from=builder /github.com/Tim-Sa/JobResume/source/assets ./assets
COPY --from=builder /github.com/Tim-Sa/JobResume/source/content_config.yaml .

CMD [ "./resume" ]