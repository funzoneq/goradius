# Installation

* Download and unpack the latest release of Goradius
```
curl -OL "https://github.com/funzoneq/goradius/releases/download/v0.0.x/goradius_Linux_x86_64.tar.gz"
tar -zxf goradius_Linux_x86_64.tar.gz
```

* Move the goradius file and create a goradius user
```
sudo mv goradius /usr/local/bin/goradius
sudo chmod +x /usr/local/bin/goradius
sudo useradd -s /sbin/nologin -M goradius
```

* Configure goradius
```
sudo mkdir /etc/goradius
sudo vi /etc/goradius/goradius.conf # See goradius.example.conf
```

* Make a systemd file and start the service
```
curl -o /usr/lib/systemd/system/goradius.service -L https://raw.githubusercontent.com/funzoneq/goradius/refs/heads/main/contrib/goradius.service
sudo systemctl daemon-reload
systemctl enable goradius
systemctl start goradius
journalctl -fu goradius
```


