# GRVM [![Build Status](https://travis-ci.org/Bugagazavr/grvm.svg?branch=master)](https://travis-ci.org/Bugagazavr/grvm)

RVM replacement

## Features

* Less shell scripts
* Builded on Go with BoltDB
* Ruby-build
* Hooks

## TODO

* Gemsets maybe

## Installation

```sh
curl -s https://raw.githubusercontent.com/Bugagazavr/grvm/master/install.sh | bash --
```

Do not forget add to your `zshrc` or `bashrc`:

```sh
# GRVM
[[ -s "$HOME/.grvm/scripts/grvm" ]] && source $HOME/.grvm/scripts/grvm
```

## Development

You need Go 1.5 + for development

```sh
mkdir -p $GOPATH/src/github.com/Bugagazavr/grvm
git clone https://github.com/Bugagazavr/grvm.git $GOPATH/src/github.com/Bugagazavr/grvm
cd $GOPATH/src/github.com/Bugagazavr/grvm
make devinstall
source ~/.grvm/scripts/grvm
```

## Contribution

Your welcome
