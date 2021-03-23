FROM golang:1.16.0-buster as builder

# install git
RUN apt-get install -y git

RUN mkdir /build

COPY . /build/

WORKDIR /build

RUN CGO_ENABLED=0 GOOS=linux go build -o bin .

FROM golang:1.16.0-buster

RUN mkdir /api
RUN addgroup --system dinamicka
RUN adduser --system --disabled-password --no-create-home --home /api --ingroup dinamicka dinamicka
RUN chown dinamicka:dinamicka /api

# Switch current root user to transcoder, prevent running service from root
USER dinamicka

COPY --from=builder /build/bin /api/

WORKDIR /api

LABEL   Name="Dinamicka Api"

#Run service
ENTRYPOINT ["./bin"]