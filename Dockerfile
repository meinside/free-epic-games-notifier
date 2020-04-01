# https://github.com/chromedp/docker-headless-shell
FROM chromedp/headless-shell:latest

# install dumb-init and ca-certificates
RUN apt-get update && \
    apt-get install -y dumb-init ca-certificates

# copy pre-built binary and config file
COPY /free-epic-games-notifier /
COPY /epic_notifier.json /

ENTRYPOINT ["dumb-init", "--"]
CMD ["/free-epic-games-notifier"]

