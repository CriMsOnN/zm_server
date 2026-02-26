#!/bin/bash

if ! command -v node &> /dev/null; then
    echo "nodejs could not be found, generate secrets manually"
    exit 1
fi

secret=$(node -e "console.log(require('crypto').randomBytes(32).toString('hex'))")
secret=$(echo $secret | tr -d '\n')
sed -i "s/FIVEM_BACKEND_SECRET=.*/FIVEM_BACKEND_SECRET='$secret'/" backend/.env