# ASD
ASD is A ~~SH-itty~~ SH-script Downloader. When installing things 90% of the time we Google around for the exact list of instructions and copy paste them into the Terminal. This process gets redundant and cumbersome over time when installing the same things for multiple devices.

ASD is a very simple tool intended to retrieve a curated collection of `.sh` scripts from [asd-modules](https://github.com/zweicoder/asd-modules) and install them after some minor dependency resolution. It falls back to `apt-get install` if the script is not found.

The scripts in [asd-modules](https://github.com/zweicoder/asd-modules) currently mainly support Linux, and MacOS to a lesser extent, but will get better as I add more scripts for personal use / get PRs from Internet strangers :)


## Usage
Currently the only thing this does is simple dependency resolution and generating the script. It can't run the installation commands because there's no good way to pass in a large string of `bash` commands into `os/exec`

Generate the script

```
# Generates an install script with commands to install docker, then drone
asd gen drone
# Run it locally
bash -s < install.sh

# Or run it remotely over ssh
ssh user@remote.com 'bash -s' < install.sh
```
