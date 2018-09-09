#!/bin/bash

set -x

sysctl net.ipv4.ip_forward=1
sysctl -p

export DEBIAN_FRONTEND=noninteractive

apt update && apt-get -y \
    -o Dpkg::Options::="--force-confdef" \
    -o Dpkg::Options::="--force-confold" \
    upgrade

apt install -y --no-install-recommends \
    mosh \
    software-properties-common

add-apt-repository -y ppa:wireguard/wireguard

apt update && apt install -y --no-install-recommends \
    wireguard

apt autoremove -y

ufw allow 8270/udp
ufw allow 51820/udp
ufw allow 22/tcp
ufw enable

umask 077 && cat > /etc/wireguard/wg0.conf << EOF
[Interface]
Address = 192.168.242.1/28
Address = fd86:ea04:1115::1/64
SaveConfig = true
PostUp = iptables -A FORWARD -i wg0 -j ACCEPT; iptables -t nat -A POSTROUTING -o eth0 -j MASQUERADE; ip6tables -A FORWARD -i wg0 -j ACCEPT; ip6tables -t nat -A POSTROUTING -o eth0 -j MASQUERADE
PostDown = iptables -D FORWARD -i wg0 -j ACCEPT; iptables -t nat -D POSTROUTING -o eth0 -j MASQUERADE; ip6tables -D FORWARD -i wg0 -j ACCEPT; ip6tables -t nat -D POSTROUTING -o eth0 -j MASQUERADE
ListenPort = 51820
PrivateKey = __WG_PRIVATE_KEY__
EOF

sed -ir "s/__WG_PRIVATE_KEY__/$( wg genkey )/" /etc/wireguard/wg0.conf

systemctl enable wg-quick@wg0
wg-quick up wg0

update-locale LANG=en_US.utf-8 LC_MESSAGE=POSIX
locale-gen
LC_ALL=en_US.utf-8 mosh-server

clear
wg show
