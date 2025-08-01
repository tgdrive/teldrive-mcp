FROM alpine:latest as certs
RUN apk --no-cache add ca-certificates

FROM scratch
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY teldrive-mcp /teldrive-mcp
EXPOSE 8080
ENTRYPOINT ["/teldrive-mcp"]
