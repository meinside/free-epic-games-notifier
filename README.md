# Free Epic Games Notifier

Fetch new free games from Epic Games' store and notify them through Pushbullet and etc.

## Install

```bash
$ git clone https://github.com/meinside/free-epic-games-notifier.git
$ cd free-epic-games-notifier/
$ go build
```

## Configure

```bash
$ cp epic_notifier.json.sample epic_notifier.json
$ vi epic_notifier.json
```

and set your tokens and urls there.

## Run

### Directly

```bash
$ ./free-epic-games-notifier
```

### With Docker

Build image,

```bash
$ docker build -t MY_IMAGE_TAG .
```

then run it with:

```bash
$ docker run -t MY_IMAGE_TAG
```

## License

MIT

