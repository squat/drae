FROM scratch
MAINTAINER Lucas Servén <lserven@gmail.com>
COPY bin/drae /
ENTRYPOINT ["/drae"]
CMD ["api"]
