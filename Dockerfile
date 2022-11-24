FROM golang:1.19 as build_stage

COPY . /src
WORKDIR /src/cmd/main

RUN CGO_ENABLED=0 GOOS=linux go build -o /app

FROM scratch

COPY --from=build_stage /app /

COPY --from=build_stage /src/key/*.pem /

EXPOSE 5555

CMD ["./app"]