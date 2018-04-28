
# Open relay

TCP server written for [Warrior project](https://cheikhseck.itch.io/warrior-project). Download the windows build and open it to run a WP server.

## Builds
- Windows [Download](https://github.com/cheikhshift/open-relay/raw/master/open-relay.exe). Learn how to find your IP address [here](https://support.microsoft.com/en-us/help/15291/windows-find-pc-ip-address)

## Building from source

#### Requirements 
- Any version of GO.

#### Install package
Download and install package with command : `go get github.com/cheikhshift/open-relay`

## About

Open relay using delimiters to parse and execute requests. NO JSON.

## CMD usage

	  -hostname string
    	Hostname the server should listen on. (default "localhost")
 	  -port int
    	Port number to listen on. (default 3333)
	
## TODO
Document open-relay dialog format.
