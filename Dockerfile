FROM go AS build

ARG VERSION

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY cmd /build/cmd
COPY internal /build/internal
COPY pkg /build/pkg

RUN CGO_ENABLED=0 go build \
  -ldflags "-s -w -X main.Version=${VERSION}" \
  -o pincer \
  ./cmd/pincer/main.go

FROM rocky AS run

COPY --from=build /build/pincer /usr/local/bin/pincer

RUN echo "pincer:x:1001:1001:pincer:/:/sbin/nologin" >> /etc/passwd && \
  echo "pincer:x:1001:" >> /etc/group && \
  echo "pincer:*:19820:0:99999:7:::" >> /etc/shadow

USER pincer

EXPOSE 8080
EXPOSE 9000
ENTRYPOINT [ "/usr/local/bin/pincer" ]
CMD [ "run" ]
