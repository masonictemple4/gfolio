---
title: "Installing Popos on macbook pro"
description: "NOT A TUTORIAL. This is a running log I use more as a cheat sheet for when I run into issues or document processes to make my life easier. It contains various fixes, configuration steps, and debugging notes."
subtitle: "cheat sheet for installing pop_os on macbook pro"
authors:
    - username: "masonictemple4"
tags:
    - "general"
    - "golang"
    - "docker"
    - "linux"
    - "os"
    - "postgres"
    - "neovim"
---

# Installing Pop_OS on Macbook pro
NOT A Tutorial. This is a running log I use more as a cheat sheet for when I run into issues or document processes to make my life easier. It contains various fixes, configuration steps, and debugging notes.


Hope you find something here helpful :) 


### Linux Wifi Fix
[gist here](https://gist.github.com/torresashjian/e97d954c7f1554b6a017f07d69a66374)
```
1. $ sudo apt-get purge bcmwl-kernel-source
2. $ sudo apt update
3. $ sudo update-pciids
4. $ sudo apt install firmware-b43-installer
5. $ sudo reboot # reboots your machine
6. $ sudo iwconfig wlp3s0 txpower 10dBm
```

A couple of things to note here: replace `wlp3s0` with the correct wireless device,
mine is defaulting to this for now so that's what i put. Second, the last command
`sudo iwconfig wlp3s0 txpower 10dBm` has to be run every time you power off and on.

**Upload/download speed debugging**:
Thanks so much for sharing! The above solution fixed the auth loop problem,  however after running some speedtests I found that I'm only pulling about 1/10th of what I do when I boot into macOS. 

Some quick system information for context.


System Info
```
OS: Pop!_OS 22.04 LTS x86_64 
Host: MacBookPro14,3 1.0 
Kernel: 6.5.4-76060504-generic 
Uptime: 15 hours, 52 mins 
Packages: 1884 (dpkg), 7 (flatpak) 
Shell: bash 5.1.16 
Resolution: 2880x1800, 3840x2160 
DE: GNOME 42.5 
WM: Mutter 
WM Theme: Pop 
Theme: Pop-dark [GTK2/3] 
Icons: Pop [GTK2/3] 
Terminal: tmux 
CPU: Intel i7-7820HQ (8) @ 3.900GHz 
GPU: AMD ATI Radeon RX 460/560D / Pro 450/455/460/555/555X/560/560X 
Memory: 10450MiB / 15881MiB
```

Wireless Card: `03:00.0 Network controller [0280]: Broadcom Inc. and subsidiaries BCM43602 802.11ac Wireless LAN SoC [14e4:43ba] (rev 02)`
OS: Pop!_OS 22.04 LTS x86_64 
Host: MacBookPro14,3 1.0 
Kernel: 6.5.4-76060504-generic 
Uptime: 15 hours, 52 mins 
Packages: 1884 (dpkg), 7 (flatpak) 
Shell: bash 5.1.16 
Resolution: 2880x1800, 3840x2160 
DE: GNOME 42.5 
WM: Mutter 
WM Theme: Pop 
Theme: Pop-dark [GTK2/3] 
Icons: Pop [GTK2/3] 
Terminal: tmux 
CPU: Intel i7-7820HQ (8) @ 3.900GHz 
GPU: AMD ATI Radeon RX 460/560D / Pro 450/455/460/555/555X/560/560X 
Memory: 10450MiB / 15881MiB 
Confirmed the wlp3s0 devices is using brcmfmac `$ sudo readlink /sys/class/net/wlp3s0/device/driver` I see that it's using `brcmfmac`. 

Then ran `$ lsmod | grep brcmfmac` which outputs:
```
brcmfmac_wcc           12288  0
brcmfmac              507904  1 brcmfmac_wcc
brcmutil               16384  1 brcmfmac
cfg80211             1257472  1 brcmfmac
```

I've also verified that `bcmwl-kernel-source` is not installed.

My `nmcli` output:
```
wlp3s0: connected to ....
	"Broadcom and subsidiaries BCM43602"
	wifi (brcmfmac), XX:XX:XX:XX:XX:XX, hw, mtu 1500
	ip4 default, ip6 default
	inet4 192.168.1.195/24
	route4 192.168.1.0/24 metric 600
	route4 169.254.0.0/16 metric 1000
	route4 default via 192.168.1.1 metric 600
	inet6 .....
	inet6 .....
	inet6 .....
	inet6  .....
	route6 .....
	route6 .....
	route6 .....
	route6 default via .... metric 600

docker0: connected (externally) to docker0
	"docker0"
	bridge, 02:42:68:7D:99:42, sw, mtu 1500
	inet4 172.17.0.1/16
	route4 172.17.0.0/16 metric 0

p2p-dev-wlp3s0: disconnected
	"p2p-dev-wlp3s0"
	wifi-p2p, hw

lo: unmanaged
	"lo"
	loopback (unknown), 00:00:00:00:00:00, sw, mtu 65536

DNS configuration:
	servers: 192.168.1.1
	domains: lan
	interface: wlp3s0

	servers: .....
	domains: lan
	interface: wlp3s0

Use "nmcli device show" to get complete information about known devices and
"nmcli connection show" to get an overview on active connection profiles.

Consult nmcli(1) and nmcli-examples(7) manual pages for complete usage details.
```

Iwconfig output:
```
lo        no wireless extensions.

wlp3s0    IEEE 802.11  ESSID:"....."  
          Mode:Managed  Frequency:5.785 GHz  Access Point: XX:XX:XX:XX:XX:XX   
          Bit Rate=650 Mb/s   Tx-Power=31 dBm   
          Retry short limit:7   RTS thr:off   Fragment thr:off
          Encryption key:off
          Power Management:off
          Link Quality=61/70  Signal level=-49 dBm  
          Rx invalid nwid:0  Rx invalid crypt:0  Rx invalid frag:0
          Tx excessive retries:86  Invalid misc:0   Missed beacon:0

docker0   no wireless extensions.
```

I've done some digging and solved one issue along the way improving my signal by adding the `brcmfmac43602-pcie.txt` file manually. However, which helped a quite a bit especially once I switched to the 5GH network bumping me to about 1/3 of what I get in macOS. 

I'm guessing it has something to do with the missing files/errors i found in the dmesg logs. 
```
[    7.494020] brcmfmac: brcmf_fw_alloc_request: using brcm/brcmfmac43602-pcie for chip BCM43602/2
[    7.494775] brcmfmac 0000:03:00.0: Direct firmware load for brcm/brcmfmac43602-pcie.Apple Inc.-MacBookPro14,3.bin failed with error -2
[    7.501904] brcmfmac 0000:03:00.0: Direct firmware load for brcm/brcmfmac43602-pcie.clm_blob failed with error -2
[    7.502152] brcmfmac 0000:03:00.0: Direct firmware load for brcm/brcmfmac43602-pcie.txcap_blob failed with error -2
[    7.623188] applesmc: key=911 fan=2 temp=46 index=45 acc=0 lux=0 kbd=0
[    7.623288] applesmc applesmc.768: hwmon_device_register() is deprecated. Please convert the driver to use hwmon_device_register_with_info().
[    7.699958] mc: Linux media interface: v0.10
[    7.987325] Bluetooth: hci0: BCM: failed to write update baudrate (-16)
[    7.987335] Bluetooth: hci0: Failed to set baudrate
[    7.988472] brcmfmac_wcc: brcmf_wcc_attach: executing
[    7.991241] brcmfmac: brcmf_c_process_clm_blob: no clm_blob available (err=-2), device may have limited channels available
[    7.991247] brcmfmac: brcmf_c_process_txcap_blob: no txcap_blob available (err=-2)
[    7.992042] brcmfmac: brcmf_c_preinit_dcmds: Firmware: BCM43602/2 wl0: Nov 10 2015 06:38:10 version 7.35.177.61 (r598657) FWID 01-ea662a8c
```

Has anyone experienced this before? I'm not sure what to do next, I've looked all over for `brcm/brcmfmac43602-pcie.clm_blob` and ` brcm/brcmfmac43602-pcie.txcap_blob` I haven't been able to find them, even in the linux-firmware repos.

As suggested [here](https://gist.github.com/MikeRatcliffe/9614c16a8ea09731a9d5e91685bd8c80#file-brcmfmac43602-pcie-txt) I tried install `firmware-brcm80211`  with apt but wasn't able to find a source for apt so I ended up just installing it from the .deb file. I had to --force-overwrite any files that already existed from linux-firmware. Unfortunately that didn't fix the missing pcie.clm_blob and txcap_blob files it did make a significant difference. And I began seeing 5GH networks in my wifi list, after connecting to the it the speeds have been about 10x what they were, and I haven't had to redo the power config after rebooting.

** Note after powering off and sitting for a day or two, and then running sytemupdates this has broken again and i'm not seeing the log messages to fix it. Now it's stuck on 2.5 GH. I'm assuming some files got overwritten. At this point it may make more sense to just buy a dongle.

### Inverted scrolling like mac
To change the settings so that scrolling is the same as on the mac,
you're looking in settings `Mouse & Touchpad` and want to enable
`Natural scrolling`

### Install Snapd
Because most of the libraries are behind that are available via the default 
sources in apt, install `$ sudo apt install snapd` this will allow us
to install more up to date tools.

Keep in mind, these may or may not be complete with security updates to install
at your own risk.

Now we can use `$ sudo snap install ...` when we need to.

### Installing Neovim
There are a few different ways to go about installing Neovim, most of which seem reasonable and
correct. However, they all have issues. Installing via the pop os store is just some limited
environment terminal version of it so that won't do. Installing via `apt` is 3+ versions behind.
The next logical thought is to use snap.. **DO NOT install via `snap`** It mostly works,
however none of your more complicated LSP servers will attach to the buffers if installed via snap.

The best option to install neovim is to **BUILD FROM SOURCE**

1. Setup prerequisites [here](https://github.com/neovim/neovim/wiki/Building-Neovim#build-prerequisites) `$ sudo apt-get install ninja-build gettext cmake unzip curl`
2. Clone the repo `$ git clone https://github.com/neovim/neovim`
3. Navigate to the repo `$ cd neovim`
4. (building stable release) `$ git checkout stable`
5. Next run make `$ make CMAKE_BUILD_TYPE=Release`
6. Finally, install `$ sudo make install`

See also [documentation](https://github.com/neovim/neovim/wiki/Installing-Neovim#install-from-source)

Note you can configure different install locations to make uninstalling easier, however i chose to go with the default basic installation to keep things kosher.

### Installing Plug
Follow [these instructions](https://github.com/junegunn/vim-plug#unix-linux) to install
vim plugged before setting up your nvim config.

Neovim install:
```
$ curl -fLo ~/.vim/autoload/plug.vim --create-dirs \
    https://raw.githubusercontent.com/junegunn/vim-plug/master/plug.vim
```


### Tmux Fixes
In order for theme settings to transfer when using tmux make sure you add
the following to your `~/.tmux.conf` file:

```
set-option -g default-terminal "screen-256color"
set -ga terminal-overrides ",*256col*:Tc"
```

It's important to note that these settings are for when your environment variable
`TERM=screen-256color` is set.

Also quick note, installing with snap had all kinds of bugs, had to **install with
apt instead**.

### Installing go
Follow the instructions [here](https://go.dev/doc/install). Important to note
I just set the `PATH` in `/etc/profile` so it was global.

Verify with `$ go version` and you're all set.


### Gopls installation
After battling it out with Gopls and the snap installation, for some reason 
it would not attach to the buffer in neovim. So instead of using that version
make sure it's uninstalled. Then instead make sure your `$ go env GOBIN` is in 
your `$PATH` and install it like so `$ go install golang.org/x/tools/gopls@latest`.

Check to make sure it worked: `$ gopls verison`. If not double check that
your `go env GOBIN` is set to `$GOPATH/bin` and that it is in your `$PATH`.
You can add it somewhere like `/etc/profile`


### [Lua language server](https://luals.github.io/wiki/build/)
Installing the lua language server has to be done manually. After you build the
repo from scratch, you can create a wrapper to use it correctly.

This requires `ninja` so make sure you have that installed before diving through
this.

See below for instructions:
1. Make an lsp config directory inside of `~/.config/nvim/` such as `lsps`
2. Cd into the directory from above.
3. Clone in the [lua-language-server repo](https://github.com/LuaLS/lua-language-server): `$ git clone https://github.com/LuaLS/lua-language-server`
4. Cd into lua-language-server
5. run `$ ./make.sh` and wait for it to finish.
6. Create a new file `lua-language-server` just like that, no extensions.
```
#!/bin/bash
# Note this is the path to the lua-language-server directory you
# previous ran ./make.sh in. You can find the bin directory inside there
# which contains the generated executable.
exec "$HOME/.config/nvim/lsps/bin/lua-language-server" "$@"
```
7. Make the new file executable `$ chmod +x lua-language-server`
8. Move it to your `/usr/local/bin` directory and you're all set.

### Installing fonts
Installing fonts is pretty easy, just download the files you want `.ttf`s and then
double click it and at the top of the window click install.

Once inside of the terminal you can select that font inside the profile of the terminal
preferences you're using.

### Setting up pbcopy
This will require the `xsel` utility which can be installed via `$ sudo apt install xsel`

Update your `~/.bashrc` file with the following:
```
alias pbcopy='xsel --clipboard --input'
alias pbpaste='xsel --clipboard --output'
```

save the file and restart your terminal or just `$ source ~/.bashrc` to activate it.

### Docker Engine && Docker Desktop setup
- First step is to install the docker engine on your system following this [documentation](https://docs.docker.com/engine/install/ubuntu/).
- After setting up the docker engine you can follow the [documentation](https://docs.docker.com/desktop/install/ubuntu/) to setup docker desktop.

##### Postgresql Docker setup
```
$ docker run --name postgres -e POSTGRES_PASSWORD=mysecretpassword -d -p 5432:5432 postgres

# To attach yourself to the docker container
$ docker exec it postgres bash

$ psql -U postgres -W postgres
```

### Installing PGADMIN
Installing PgAdmin can be done following this [documentation](https://www.pgadmin.org/download/pgadmin-4-apt/). 

The following is the instruction exerpt from the documentation linked above:
```
# Install the public key for the repository (if not done previously):
$ curl -fsS https://www.pgadmin.org/static/packages_pgadmin_org.pub | sudo gpg --dearmor -o /usr/share/keyrings/packages-pgadmin-org.gpg

# Create the repository configuration file:
$ sudo sh -c 'echo "deb [signed-by=/usr/share/keyrings/packages-pgadmin-org.gpg] https://ftp.postgresql.org/pub/pgadmin/pgadmin4/apt/$(lsb_release -cs) pgadmin4 main" > /etc/apt/sources.list.d/pgadmin4.list && apt update'

# Install pgAdmin

# Install for both desktop and web modes:
$ sudo apt install pgadmin4
```
