FROM golang:1.9
ARG BRANCH
ARG COMMIT
ARG PORT
WORKDIR /go/src/github.com/les/bluewhales/
COPY . .
RUN make all
ENV PORT="${PORT}"
HEALTHCHECK --interval=5s --retries=3 --timeout=3s \
            CMD curl -f "http://127.0.0.1:${PORT}/health/" || exit 1
CMD /go/bin/bluewhales
EXPOSE "${PORT}"
LABEL name="bluewhales" \
      description="Balaenoptera musculus" \
      maintainer="les" \
      license="GPL-3.0" \
      url="https://github.com/les/bluewhales" \
      version.branch="${BRANCH}" version.commit="${COMMIT}"
