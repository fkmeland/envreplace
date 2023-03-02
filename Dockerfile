FROM scratch
WORKDIR /app
COPY envreplace /app/
ENTRYPOINT ["/app/envreplace"]