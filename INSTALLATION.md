# Installation

* Download and unpack the latest release of Goradius
```
curl -OL "https://github.com/funzoneq/goradius/releases/download/v0.0.x/goradius_Linux_x86_64.tar.gz"
tar -zxf goradius_Linux_x86_64.tar.gz
```

* Move the goradius file and create a goradius user
```
sudo mv goradius /usr/local/bin/goradius
sudo useradd -s /sbin/nologin -M goradius
```

* Make a systemd file and start the service
```
vi /usr/lib/systemd/system/goradius.service
systemctl enable goradius
systemctl start goradius
journalctl -fu goradius
```


