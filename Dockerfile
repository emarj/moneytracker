FROM scratch

COPY ./bin/moneytracker_linux ./moneytracker

EXPOSE 3245
ENTRYPOINT ["/moneytracker"]