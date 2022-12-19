FROM debian:latest
RUN mkdir -p /service
WORKDIR /service
ADD ./EWallet.tar.gz /service/
EXPOSE 8080
CMD ./EWallet