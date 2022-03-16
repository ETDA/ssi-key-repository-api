FROM golang:1.16-bullseye

RUN apt-get update -y && apt-get upgrade -y && \
apt-get install bash git openssh-server -y
WORKDIR /app
RUN git config --global url."https://core-deploy:a4NaMfxHGhtfEtuGSuKX@ssi-gitlab.teda.th".insteadOf "https://ssi-gitlab.teda.th"
ADD go.mod go.sum /app/
RUN go mod download
RUN go get -u github.com/pilu/fresh
ADD . /app/
ADD libcs_pkcs11_R2.so /usr/local/lib/libcs_pkcs11_R2.so
ADD cs_pkcs11_R2.cfg /etc/utimaco/cs_pkcs11_R2.cfg
CMD fresh -c other_runner.conf
