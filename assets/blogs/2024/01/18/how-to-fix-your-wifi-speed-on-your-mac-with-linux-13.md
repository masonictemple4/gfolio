---
title: How to fix your wifi speed on your mac with Linux
subtitle: improve linux wifi speed on your mac
description: "If you have a certain wifi card, and you've noticed your wifi speeds are laughably slow compared to that of when you used macos. This solution might help you!"
authors:
    - username: "masonictemple4"
      profilepicture: ""
tags:
    - "linux"
    - "wifi"
    - "os"
    - "bcm43602"
    - "broadcom"
    - "firmware"
---

# How to fix slow wifi speeds for linux on mac 
Hey all, I recently dove pretty deep into the wifi issues on Pop_OS using my
old macbook pro 15" with touchbar. 

For me this was a multi-part solution. First I had to implement a fix just to 
connect to the wifi without the endless authentication loop.

By no means am I some sort of expert here, most of my solutions came from all-night
deep dives through the depths of the internet to find anyone who had similar issues
that came up with a solution.

The fix for this was suprisingly silly and I'm still not sure why it fixes the problem, however `$ sudo iwconfig` will list your interfaces. For me I have `wl3ps0`, next run `$ sudo iwconfig wl3ps0 txpower 10dBm`, and try reconnecting to your network.

Others have said turning power management off all together has worked as well, for me however this was what did the trick. To try turning power management off you can run `$ sudo iwconfig wl3ps0 power off`. 

***Note:*** You will have to run this command every time you reboot.


Okay, so now we at least have internet connection. After a while you might notice your
internet speed is laggy, run a couple speed tests to see how it's performing compared
to what you expect.

For me I noticed I was getting 1/10th of the speed I would normally get when booting
into macOS. I pieced together bits of debugging steps I found from various forms you will
see those below.

##### Debugging tools
1. Find the exact version of wifi card you have by running `$ lspci -vvnn | grep Broadcom`, because I already know that
my wifi card is from Broadcom I can grep to filter the output. This command returns:

    `03:00.0 Network controller [0280]: Broadcom Inc. and subsidiaries BCM43602 802.11ac Wireless LAN SoC [14e4:43ba] (rev 02)`

    You will soon find, this model relies on a non-free propiatary version of their firmware, that Apple will not let them release.

2. Verify that you do not have: `bcmwl-kernel-source` installed, and if you do `purge` it and reboot.

3. Determine what firmware your card is actually using: `$ sudo readlink /sys/class/net/wlp3s0/device/driver`. You'll see
something similar to this `../../../../bus/pci/drivers/brcmfmac` in the output. **brcmfmac** Is what we're looking for here.

4. To see what uses this firmware here: `$ lsmod | grep brcmfmac` 

5. After a reboot, check your **dmesg** logs. This in part with the information above is what helped me fix the wifi speed.
There are various ways to check this log. If you're actively iterating on solutions, it might not be a bad idea to
save each log so you can visit them later. By doing something like `$ sudo dmesg > postfirmware-upgrade-log.txt`.

    You can also see the file located at `/var/log/dmesg`

    **NOTE:** Now is a good time to start grepping or searching through this file with the info we found previously.

Desparate, I sent the following snippet from my **dmesg** output to ChatGPT to see what any of it meant:
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

Sure enough, the errors are from missing files. Countless hours later, I found this hidden gem `brcmfmac43602-pcie.txt` file using the gist below. 

### The solution
The solution that worked for me ended up slightly different than what the gist suggests. The first couple of times I attempted this, 
I attempted the firmware instructions as well. However, I wasn't able to install the suggested firmware with `apt` and had to dig
to find the files separately. Digging through the linux-firmware repo and downloading `.deb` files from debian to try and --force-overwrite.
Which in the end I think only complicated things, and caused other issues behind the scenes.

The working solution was:
1. Download the [brcmfmac43602-pcie.txt](https://gist.github.com/MikeRatcliffe/9614c16a8ea09731a9d5e91685bd8c80#file-brcmfmac43602-pcie-txt) template.
2. Replace the macaddr=XX:XX:XX:XX:XX:XX on line 8, with your own.
    To find your mac address any of the following can work:
    ```
    # Remember the interface id from before (mine was wlp3s0)
    $ sudo ifconfig
    # mine was located immediately to the right of `ether` under the ipv6 addresses

    or 

    $ sudo ip link show
    # mine was the second device in the list, and on the second line following `link/ether` was the mac address.

    ```
    If you're not sure, my Mac address also shows in my wifi-settings when clicking on the
    cog of the network i'm connected to labeled **Hardware Address**
3. Save the file. Then copy to firmware location `$ cp brcmfmac43602-pcie.txt /lib/firmware/brcm/brcmfmac43602-pcie.txt`.
4. `$ sudo rmmod brcmfmac` (if this throws an error because it's in use by brcmf_wcc then run `$sudo rmmod brcmfmac_wcc` and then try again.
5. `$ sudo modprobe`

My network services restarted on their own, however to be safe you can always reboot. This solution enabled the card to see our 5GHz network. I also noticed I no longer had to
run our previous `$ sudo iwconfig ... txpower 10dBm` command, it connects on it's own and finally runs at acceptable speeds compared to what it is supposed to. 

### In conclusion
This is not a perfect solution, in fact I didn't even realize it had improved until making a forum post and I saw the frequency and bit rate from the `iwconfig` output that I had made an improvement. Which, led to frequent speed tests. They're not all great but I can confidently say it has made a 10x difference. Every now and again, it will drop off.

Make sure to keep an eye on it and check it after system updates to see if there were any breaking changes. 


Honestly, I probably should have just bought a dongle and been done with it. But, hopefully someone out there finds this useful!! 


Feel free to reach out if you have had a similar or completely different experience, or if you'd like to add anything. 

Happy hunting :)



## Credits:
Huge thank you to Mike for the [brcmfmac43602-pcie.txt gist](https://gist.github.com/MikeRatcliffe/9614c16a8ea09731a9d5e91685bd8c80#file-brcmfmac43602-pcie-txt) that ended up fixing both speed and connection issues.

A big thanks to Miguel for the [tx power fix](https://gist.github.com/torresashjian/e97d954c7f1554b6a017f07d69a66374) that helps resolve infinite auth loops when connecting to wifi network.



