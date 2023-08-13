FROM golang:1.21 AS build
RUN echo "Based on commit: $GIT_HASH"

# All these steps will be cached
RUN mkdir /app
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
# COPY the source code as the last step
COPY . .

RUN make build_in_docker

# step 2 - create container to run
FROM alpine:latest
# It's essential to regularly update the packages within the image to include security patches
RUN apk update && apk upgrade

# Reduce image size
RUN rm -rf /var/cache/apk/* && rm -rf /tmp/*

# Avoid running code as a root user
RUN adduser -D appuser
USER appuser
WORKDIR /app
COPY --from=build /app /app
COPY --from=build /app/bin/binary /app/
COPY --from=build /app/configs/* /app/configs/
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
# RUN chmod +x /app/binary
EXPOSE 8000/tcp
CMD /app/binary