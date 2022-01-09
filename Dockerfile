FROM golang as builder

RUN mkdir /build

COPY . /build/

WORKDIR /build

RUN CGO_ENABLED=0 GOOS=linux go build -o bin .

FROM golang

ENV PORT=$PORT
ENV SERVER_ADDRESS=:5555
ENV CLIENT_LOCATION=/api/public
ENV REDIS_ADDRESS=:6379
ENV REDIS_PASSWORD=""

RUN mkdir /api

WORKDIR /build

COPY --from=builder /build/bin /api/
COPY client/build /api/public

WORKDIR /api

LABEL   Name="Chat Api"

#Run service
ENTRYPOINT ["./bin"]