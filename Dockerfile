FROM alpine
COPY envreplace /

USER 1001
ENTRYPOINT ["envreplace"]
CMD ["-h"]