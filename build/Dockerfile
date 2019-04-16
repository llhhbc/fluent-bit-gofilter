FROM fluent/fluent-bit:1.0.4

# CMD ["/fluent-bit/bin/fluent-bit", "-c", "/fluent-bit/etc/fluent-bit.conf"]

COPY ./fluent-bit /fluent-bit/bin/

#RUN mkdir /fluent-bit/lib/

COPY ./cgolib.so /fluent-bit/lib/



