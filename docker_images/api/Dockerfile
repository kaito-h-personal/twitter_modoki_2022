## 使用するGoのバージョンを指定
ARG GO_DOCKER_IMAGE_VERSION=${GO_DOCKER_IMAGE_VERSION}
## コンテナ内で使用するディレクトリを指定
ARG API_CODE_WORKDIR=${API_CODE_WORKDIR}

FROM golang:${GO_DOCKER_IMAGE_VERSION}

RUN apk update && apk add git alpine-sdk

# ワーキングディレクトリの設定(TODO: WORKDIRはENVを使わないと変数展開できない)
WORKDIR /src
