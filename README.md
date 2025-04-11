# goradius
A radius server for Juniper BNG written in Go

We authenticate customers based on the customers.json file. The customers.json file contains multiple customers:

```
{
    "identifier": "olt01:10-1:6b86b273ff",
    "status": "suspended",
    "speed_up": "250m",
    "speed_down": "250m",
    "vrf": "BNG-Users",
    "static_routes": ["10.0.0.0/29"]
},
{
    "identifier": "olt01:10-8:2c624232cd",
    "status": "active",
    "speed_up": "1g",
    "speed_down": "1g",
    "vrf": "BNG-NAT",
    "static_routes": []
}
```

The identifier field is whatever the BNG router has configured to send as the username. Usually this is some option82 information that the OLT adds identifying the circuit-id. Status can be active or suspended. The speed up and down will get added as a Input and Output interface filter. The VRF will decide what routing instance the customer ends up in, think: NAT instance, public IP instance or Portal. Optionally the customer can receive a list of 0 or more static routes that will be routed to the the /32 IP of the customer.

# Configuration modes
You can run goradius in two modes: Authenticated or Unauthenticated.

## Unauthenticated mode
In unauthenticated mode, we will allow everybody in and give them a default filter of 1G and place them in the BNG-NAT VRF.

### goradius.conf
```
{
    "AuthEnabled": false,
    "DefaultVRF": "BNG-NAT",
    "DefaultUploadSpeed": "1G",
    "DefaultDownloadSpeed": "1G"
}
```

## Authenticated mode
In authenticated mode goradius will compare the customer to the customers.json file. If a customer is not found or not active, they will be denied access.

### goradius.conf
```
{
    "CustomerFile": "/etc/goradius/customers.json",
    "CaptivePortalEnabled": false,
    "AuthEnabled": true
}
```

## Captive portal mode
In authenticated mode with captive portal enabled, customers that are not found or not active will be allowed access to the Portal VRF. In this VRF, any request they make via a webbrowser will be transported to a captive portal where customers can either sign up or update their billing information.

### goradius.conf
```
{
    "CustomerFile": "/etc/goradius/customers.json",
    "CaptivePortalEnabled": true,
    "AuthEnabled": true
}
```

# Installation on Debian/Ubuntu

```
curl -OL https://github.com/funzoneq/goradius/releases/download/v0.0.<version>/goradius_Linux_x86_64.tar.gz
tar -zxf goradius_Linux_x86_64.tar.gz
sudo mv goradius /usr/local/bin/goradius
sudo chmod +x /usr/local/bin/goradius
sudo vi /etc/goradius/goradius.conf
sudo curl -o /usr/lib/systemd/system/goradius.service -L https://raw.githubusercontent.com/funzoneq/goradius/refs/heads/main/contrib/goradius.service
sudo systemctl daemon-reload
sudo systemctl enable goradius
sudo systemctl start goradius
```