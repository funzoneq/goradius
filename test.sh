#!/bin/bash

echo "User-Name=cr02-laca01.ps100:2048-4094,User-Password=secret" | radclient -x localhost:1812 auth secret
echo "User-Name=olt01:10-1:6b86b273ff,User-Password=secret" | radclient -x localhost:1812 auth secret
echo "User-Name=olt01:10-2:d4735e3a26,User-Password=secret" | radclient -x localhost:1812 auth secret
echo "User-Name=olt01:10-2:d4735e3a26,User-Password=secret" | radclient -x localhost:1812 auth secret
echo "User-Name=olt01:10-0:5feceb66ff,User-Password=secret" | radclient -x localhost:1812 auth secret
echo "User-Name=olt01:10-3:4e07408562,User-Password=secret" | radclient -x localhost:1812 auth secret
