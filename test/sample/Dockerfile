FROM ngrok/ngrok as ngrok_image

FROM alpine:3.18
ENV GO_VERSION=1.22.0

WORKDIR /workspace
RUN apk --no-cache add bash curl docker-cli-buildx git jq make \
    && curl -Lo go.tar.gz https://dl.google.com/go/go${GO_VERSION}.linux-amd64.tar.gz \
    && tar -C /usr/local -xzf go.tar.gz \
    && rm go.tar.gz \
    && ln -s /usr/local/go/bin/go /usr/bin/go

COPY --from=ngrok_image /bin/ngrok /bin/ngrok
