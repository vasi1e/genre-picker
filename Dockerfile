FROM golang:1.19

WORKDIR /genre-picker

COPY ./ ./

RUN CGO_ENABLED=0 go build -o genre-picker

CMD ["./genre-picker"]