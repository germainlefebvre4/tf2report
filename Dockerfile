FROM alpine:3 AS build

ARG DISTRIBUTION=linux
ARG CPU_ARCH=amd64
ARG TF2REPORT_VERSION=0.1.0

WORKDIR /app

RUN apk update && \
    apk add --no-cache curl ca-certificates

RUN curl -L --retry 5 --retry-delay 2 --retry-all-errors \
    --output tf2report.tar.gz \
    "https://github.com/germainlefebvre4/tf2report/releases/download/v${TF2REPORT_VERSION}/tf2report_${DISTRIBUTION}_${CPU_ARCH}.tar.gz" && \
    tar -zxvf tf2report.tar.gz && \
    chmod +x tf2report


FROM alpine:3

COPY --from=build /app/tf2report /usr/local/bin/tf2report

WORKDIR /app

ENTRYPOINT ["tf2report"]
CMD ["--help"]
