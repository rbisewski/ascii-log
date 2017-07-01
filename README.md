# ascii-log- harvests data from server logs and displays it via ASCII

A cut and dry golang application to generate daily data from log file
entries of common servers, such as nginx or apache and output it to a
text file so it can be accessed via w3m or lynx or wget.

Specifically it takes IP address data and executes a simple hostname
lookup for the purposes of examining what name is currently assigned
at the moment of the server request.

Feel free to fork it and use it for other projects if you find it
useful.


# Requirements

The following is needed in order for this to function as intended:

* Linux kernel 4.0+
* cron
* golang 1.6+
* apache / nginx

Older kernels could still give some kind of result, but I *think* most of
the newer versions of golang require newer kernels. Feel free to email me if
this is incorrect.


# Installation

0) Build this program as you would a simple golang module.

    make

1) Install this program on your server.

    make install

2) Adjust the cron job to set the choice of server (default is nginx).

    vim /etc/cron.d/ascii-log

Alternatively, if you are running Arch Linux w/ systemd, you can use the
included ascii-log.service instead. However, the cron job is recommended
since it has greater compatibility with more distros.

# Uninstallation

1) To remove this program from your system.

    make uninstall

2) Consider cleaning up any remaining logs, if they are no longer needed.

    rm /var/www/html/data/ip.log


# TODOs

* Consider adding the ability to display the info as a HTML5/CSS graph

* Hostname check on the IP addresses as soon as they access the server

* Create a systemd service that works better with alternative distros

# Author

Written by Robert Bisewski at Ibis Cybernetics. For more information, contact:

* Website -> www.ibiscybernetics.com

* Email -> contact@ibiscybernetics.com
