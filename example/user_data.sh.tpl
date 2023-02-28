#!/bin/bash

mkdir -p /opt/openvpn
cat > /opt/openvpn/profile.ovpn <<EOF
${profile}
EOF

set -ex
exec > >(tee /var/log/user-data.log|logger -t user-data -s 2>/dev/console) 2>&1

apt install -y apt-transport-https
wget https://swupdate.openvpn.net/repos/openvpn-repo-pkg-key.pub
apt-key add openvpn-repo-pkg-key.pub
wget -O /etc/apt/sources.list.d/openvpn3.list https://swupdate.openvpn.net/community/openvpn3/repos/openvpn3-focal.list
apt update
apt install -y openvpn3

# does not work on reboot
openvpn3 session-start --config /opt/openvpn/profile.ovpn
# add auto-reload