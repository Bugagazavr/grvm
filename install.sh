#!/bin/sh

function getArch {
	arch=`uname -m`
	case $arch in
		"386"|"i386")
			echo "386"
			;;
		"amd64"|"x86_64")
			echo "amd64"
			;;
		"arm"|"arm7")
			echo "arm"
			;;
		*)
			echo "You have unsupported arch: $arch"
			exit 1
			;;
	esac
}

function getPlatform {
	platform=`uname -s`
	case $platform in
		"Darwin"|"darwin")
			echo "darwin"
			;;
		"Linux"|"linux")
			echo "linux"
			;;
		*)
			echo "You have unsupported platform: $platform"
			exit 1
			;;
	esac
}

function getLatesDownloadUrl {
	echo $(curl -L -s https://api.github.com/repos/Bugagazavr/grvm/releases/latest | grep 'browser_' | grep $(getPlatform) | grep $(getArch) | cut -d\" -f4)
}

function tmpCleanup {
	[ -d "/tmp/grvm" ] && rm -rf /tmp/grvm
	[ -f "/tmp/grvm.tar.gz" ] && rm -rf /tmp/grvm.tar.gz
}

case $1 in
	"devinstall")
		go build grvm.go
		rm -rf $HOME/.grvm/bin
		rm -rf $HOME/.grvm/scripts
		mkdir -p $HOME/.grvm/bin
		mkdir -p $HOME/.grvm/scripts
		cp scripts/grvm $HOME/.grvm/scripts/grvm
		cp grvm $HOME/.grvm/bin/grvm
		$HOME/.grvm/bin/grvm doctor
		;;
	*)
		echo "Prepare to install GRVM"
		tmpCleanup

		downloadUrl=$(getLatesDownloadUrl)
		echo "Downloading $downloadUrl"
		curl -L -s $downloadUrl > /tmp/grvm.tar.gz
		echo "Downloading finished"

		echo "Extracting"
		mkdir -p /tmp/grvm
		tar -xf /tmp/grvm.tar.gz --directory=/tmp/grvm

		if [ ! -d "$HOME/.grvm" ]; then
			echo "Creating $HOME/.grvm directory"
			mkdir -p $HOME/.grvm/{bin,scripts}
		else
			echo "Cleanup existing $HOME/.grvm directory"
			[ -d "$HOME/.grvm/bin" ] && rm -f $HOME/.grvm/bin/*
			[ -d "$HOMR/.grvm/scripts" ] && rm -f $HOME/.grvm/scripts/*
			mkdir -p $HOME/.grvm/{bin,scripts}
		fi

		echo "Install GRVM"
		cp /tmp/grvm/bin/grvm $HOME/.grvm/bin/grvm
		cp /tmp/grvm/scripts/grvm $HOME/.grvm/scripts/grvm

		echo "Delete installtion files"
		tmpCleanup

		source $HOME/.grvm/scripts/grvm
		echo "GRVM has been installed"
		;;
esac
