# Simple cloud control utility
*Yet he commanded the skies above...*


## Installation

Download the repo, compile the code for your platform.

```bash
$ git clone https://github.com/bitcoin-software/nubectl
$ cd nubectl
$ ./build.sh
$ mv nubectl /usr/local/bin/nubectl

```

## Usage

**nubectl** can be used with any CBSD-based platform, self-hosted, cloud-hosted or both

### with on-premise [CBSD](https://github.com/cbsd/cbsd) cluster

via env(1):

```bash
$ export CLOUDURL="https://your-cbsd-api.endpoint.com"
$ export CLOUDKEY="/path/to/your/ssh/key.pub"
$ nubectl help
```

via args:

```bash
$ nubectl --cloudurl https://your-cbsd-api.endpoint.com --cloudkey /path/to/your/ssh/key.pub
$ nubectl -cloudurl=https://your-cbsd-api.endpoint.com -cloudkey=/path/to/your/ssh/key.pub
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
