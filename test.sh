#!/bin/bash

echo "User-Name=cr02-laca01.ps100:2048-4094,User-Password=secret" | radclient -x localhost:1812 auth secret
