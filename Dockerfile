FROM golang:1.16-bullseye

RUN apt-get update -y && apt-get upgrade -y && \
apt-get install bash git  -y

WORKDIR /app
RUN git config --global url."https://core-deploy:a4NaMfxHGhtfEtuGSuKX@ssi-gitlab.teda.th".insteadOf "https://ssi-gitlab.teda.th"
ADD . /app/
ADD libcs_pkcs11_R2.so /usr/local/lib/libcs_pkcs11_R2.so
ADD cs_pkcs11_R2.cfg /etc/utimaco/cs_pkcs11_R2.cfg
ADD go.mod go.sum /app/
RUN go mod download
ADD . /app/
RUN go build -o main

FROM ubuntu:20.04
COPY --from=0 /app/main /main
ADD libcs_pkcs11_R2.so /usr/local/lib/libcs_pkcs11_R2.so
ADD cs_pkcs11_R2.cfg /etc/utimaco/cs_pkcs11_R2.cfg
CMD ./main

