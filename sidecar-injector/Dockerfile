# Use distroless as minimal base image to package the manager binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
FROM harbor.hzxingzai.cn/tools/busybox:latest
WORKDIR /
COPY injector .
COPY certs/ ./certs
USER root:root

ENTRYPOINT ["/injector"]
