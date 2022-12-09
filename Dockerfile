FROM scratch

COPY ./bin/moneytracker_linux_amd64 ./moneytracker

EXPOSE 3245
ENTRYPOINT ["/moneytracker"]