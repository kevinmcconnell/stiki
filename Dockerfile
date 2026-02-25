FROM --platform=$BUILDPLATFORM golang:1.25 AS build
ARG TARGETOS
ARG TARGETARCH
WORKDIR /src
COPY . .
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -o /stiki .

FROM scratch
COPY --from=build /stiki /stiki
COPY templates/ /templates/
COPY public/ /public/
EXPOSE 80
ENTRYPOINT ["/stiki"]
