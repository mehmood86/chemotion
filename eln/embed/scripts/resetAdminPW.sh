#!/bin/bash

cd /chemotion/app && echo "u = User.find_by(type: 'Admin'); u.password='${1:-chemotion}'; u.account_active=true; u.save" | rails c
