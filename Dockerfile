FROM alpine
WORKDIR /app
COPY envreplace /app/

USER 1001
ENTRYPOINT ["/app/envreplace"]
CMD ["-h"]