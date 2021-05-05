# Simple cloud control utility
*Yet he commanded the skies above...*


## Installation

Download the repo, compile the code for your platform.

```bash
$ git clone https://github.com/bitcoin-software/nubectl
$ cd nubectl
$ go build main.go
$ mv main /usr/local/bin/nubectl

```

## Usage

**nubectl** can be used with any CBSD-based platform, self-hosted, cloud-hosted or both

### with on-premise [CBSD](https://github.com/cbsd/cbsd) cluster

```bash
$ export CLOUDURL="https://your-cbsd-api.endpoint.com"
$ export CLOUDKEY="/path/to/your/ssh/key.pub"
$ nubectl help
```

### with [bitclouds](https://bitclouds.sh)
```bash
$ export CLOUDKEY="/path/to/your/ssh/key.pub"
$ nubectl help
```

## WIP
This project is under heavy development. Anything can be changed rapidly for no reason.

## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.

## License
[WTFPL](http://www.wtfpl.net/)