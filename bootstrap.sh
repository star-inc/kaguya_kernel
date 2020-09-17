#!/bin/sh

cat <<EOF
Kaguya
 ---
The opensource instant messaging framework.

(c)2020 Star Inc.
EOF

echo ""
echo "Installing Go Packages"
go get -u github.com/google/uuid
go get -u github.com/gorilla/websocket
go get -u go.mongodb.org/mongo-driver/mongo

echo ""
echo "Let\`s Gopher (>w<)"
echo ""