ARG USER_NAME='app'
ARG USER_UID='1000'
ARG USER_GID='1000'
ARG USER_GROUP='app_users'
ARG REPO='github.com/kohirens/git-tool-belt'
ARG GLIBC_VER='2.34-r0'
ARG GH_VER='2.4.0'

FROM golang:1.20-alpine3.18 AS dev

ARG USER_NAME
ARG USER_UID
ARG USER_GID
ARG USER_GROUP
ARG REPO
ARG GLIBC_VER

ENV GOPATH /home/${USER_NAME}
ENV CGO_ENABLED=0
ENV WORK_DIR=/home/${USER_NAME}/src/${REPO}
ENV SHELL=/bin/sh
ENV LD_LIBRARY_PATH="/usr/lib:/lib"

# VS Code Requirements: openssh, musl, libgcc, libstdc++
RUN apk --no-progress --purge --no-cache upgrade \
 && apk --no-progress --purge --no-cache add --upgrade \
    curl \
    git \
    gnupg \
    libgcc \
    libstdc++ \
    openssh \
 && apk --no-progress --purge --no-cache upgrade \
 && rm -vrf /var/cache/apk/* \
 && rm -rf /tmp/* \
 && curl --version \
 && git --version

# Install vanilla GLibC: https://github.com/sgerrand/alpine-pkg-glibc
RUN curl -o /etc/apk/keys/sgerrand.rsa.pub https://alpine-pkgs.sgerrand.com/sgerrand.rsa.pub \
 && curl -LO https://github.com/sgerrand/alpine-pkg-glibc/releases/download/${GLIBC_VER}/glibc-${GLIBC_VER}.apk \
 && apk add glibc-${GLIBC_VER}.apk \
 && rm glibc-${GLIBC_VER}.apk \
       /etc/apk/keys/sgerrand.rsa.pub

# Add a non-root group and user, helpful if you dev on Linux.
RUN addgroup --system --gid ${USER_GID} ${USER_GROUP} \
 && adduser --system \
    --disabled-password \
    --ingroup ${USER_GROUP} \
    --uid ${USER_UID} \
    ${USER_NAME}

USER ${USER_NAME}

# Make directories with correct permissions
RUN mkdir -p ~/bin ~/src

WORKDIR "/home/${USER_NAME}/src/${REPO}"

RUN go install github.com/git-chglog/git-chglog/cmd/git-chglog@v0.15.4 \
 && go clean -cache -modcache -testcache

COPY .docker/start.sh /usr/local/bin/start.sh

ENV PATH="/home/${USER_NAME}/bin:${PATH}"

ENTRYPOINT [ "sh" ]

CMD [ "/usr/local/bin/start.sh" ]

FROM dev AS build

ARG USER_NAME
ARG USER_UID
ARG USER_GID
ARG USER_GROUP
ARG REPO

COPY . .

USER root

RUN chown -R "${USER_NAME}:${USER_GROUP}" .

USER ${USER_NAME}

RUN go install \
 && go generate \
 && go build \
 && go install

FROM alpine:3.15 AS release

ARG USER_NAME
ARG USER_UID
ARG USER_GID
ARG USER_GROUP
ARG GH_VER
ARG GH_FILE="gh_${GH_VER}_linux_amd64"

ENV GH_VER=${GH_VER}
ENV GH_FILE=${GH_FILE}

RUN apk --no-progress --purge --no-cache upgrade \
 && apk --no-progress --purge --no-cache add --upgrade \
    bash \
    curl \
    git \
    jq \
    openssh \
 && apk --no-progress --purge --no-cache upgrade \
 && rm -vrf /var/cache/apk/* \
 && rm -rf /tmp/*

# Install Github CLI tool.
RUN cd tmp \
 && curl -LO https://github.com/cli/cli/releases/download/v${GH_VER}/${GH_FILE}.tar.gz \
 && tar -xzf ${GH_FILE}.tar.gz \
 && chmod +x ${GH_FILE}/bin/gh \
 && mv ${GH_FILE}/bin/gh /usr/local/bin \
 && rm -rf ${GH_FILE}.tar.gz ${GH_FILE} \
 && gh --version

# Add a non-root group and user, helpful if you dev on Linux.
RUN addgroup --system --gid ${USER_GID} ${USER_GROUP} \
 && adduser --system \
    --disabled-password \
    --ingroup ${USER_GROUP} \
    --uid ${USER_UID} \
    ${USER_NAME}

COPY --from=build /home/${USER_NAME}/bin/git-tool-belt /usr/local/bin
COPY --from=build /home/${USER_NAME}/bin/git-chglog /usr/local/bin

USER ${USER_NAME}

ENTRYPOINT [ "git-tool-belt" ]

CMD [ "" ]