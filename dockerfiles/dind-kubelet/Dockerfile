FROM debian:9.8-slim as downloader

RUN apt-get update && \
    apt-get install -y wget tar && \
    rm -rf /var/lib/apt/lists*

RUN mkdir /glibc && \
    wget -q -O /glibc/sgerrand.rsa.pub https://alpine-pkgs.sgerrand.com/sgerrand.rsa.pub && \
    wget -q -O /glibc/glibc.apk https://github.com/sgerrand/alpine-pkg-glibc/releases/download/2.29-r0/glibc-2.29-r0.apk

ARG KUBERNETES_VERSION

RUN wget -q -O - https://dl.k8s.io/$KUBERNETES_VERSION/kubernetes-node-linux-amd64.tar.gz | \
    tar xz kubernetes/node/bin/kubelet && \
    mv kubernetes/node/bin/kubelet /kubelet && \
    chmod +x /kubelet && \
    rm -rf kubernetes

FROM docker:18.09-dind

COPY --from=downloader /glibc/sgerrand.rsa.pub /etc/apk/keys/sgerrand.rsa.pub
COPY --from=downloader /glibc/glibc.apk glibc.apk
COPY --from=downloader /kubelet /kubelet

RUN apk add glibc.apk

COPY entrypoint.sh /

ENTRYPOINT /entrypoint.sh
