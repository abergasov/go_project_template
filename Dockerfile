FROM golang:1.19 AS build
RUN echo "Based on commit: $GIT_HASH"
COPY . /app
WORKDIR /app
RUN go mod tidy && go mod vendor
RUN make build

# step 2 - create container to run
FROM debian:buster-slim
WORKDIR /app
COPY --from=build /app /app
COPY --from=build /app/bin/binary /app/
COPY --from=build /app/configs/* /app/configs/
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
RUN chmod +x /app/binary
EXPOSE 8000/tcp
CMD /app/binary