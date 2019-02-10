# Source: https://medium.com/@chemidy/create-the-smallest-and-secured-golang-docker-image-based-on-scratch-4752223b7324

############################
# STEP 1 build executable binary
############################
FROM golang@sha256:8dea7186cf96e6072c23bcbac842d140fe0186758bcc215acb1745f584984857 AS builder
# Install git + SSL ca certificates.
# Git is required for fetching the dependencies.
# Ca-certificates is required to call HTTPS endpoints.

RUN    apk update \
    && apk add --no-cache git ca-certificates \
    && update-ca-certificates

# Create appuser
RUN adduser -D -g '' appuser

WORKDIR /usr/local/src/simple-go-echo
COPY . .

# Fetch dependencies.
# Using go mod with go 1.11
RUN go mod download
# Build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go install

############################
# STEP 2 build a small image
############################
FROM scratch
ENV PORT=1919
# Import from builder.
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd
# Copy our static executable
COPY --from=builder /go/bin/simple-go-echo /go/bin/simple-go-echo
# Use an unprivileged user.
USER appuser
EXPOSE $PORT
# Run the hello binary.
ENTRYPOINT ["/go/bin/simple-go-echo"]
