#!/usr/bin/env bash

# exit on error
set -e

echo ""
echo "Install isc-dhcp-server hostapd ..."
apt install isc-dhcp-server hostapd -y
systemctl enable isc-dhcp-server
systemctl unmask hostapd
systemctl enable hostapd

echo ""
echo "Enable dhcp server for usb, wlan interfaces..."
cp /etc/default/isc-dhcp-server /etc/default/isc-dhcp-server.bak
sed -i -E 's/^#?(INTERFACESv4)=.*/\1="usb0 wlan0"/' "/etc/default/isc-dhcp-server"
diff --color "/etc/default/isc-dhcp-server.bak" "/etc/default/isc-dhcp-server" || true

echo ""
echo "Configure dhcp daemon..."
cp /etc/dhcp/dhcpd.conf /etc/dhcp/dhcpd.conf.bak
cat > /etc/dhcp/dhcpd.conf <<- EOM
# option definitions common to all supported networks...
option domain-name "local";
option domain-name-servers 8.8.8.8, 8.8.4.4;

default-lease-time 600;
max-lease-time 7200;

# The ddns-updates-style parameter controls whether or not the server will
# attempt to do a DNS update when a lease is confirmed. We default to the
# behavior of the version 2 packages ('none', since DHCP v2 didn't
# have support for DDNS.)
ddns-update-style none;

# If this DHCP server is the official DHCP server for the local
# network, the authoritative directive should be uncommented.
authoritative;

subnet 192.168.12.0 netmask 255.255.255.0 {
  range 192.168.12.100 192.168.12.200;
}

subnet 192.168.14.0 netmask 255.255.255.0 {
  range 192.168.14.100 192.168.14.200;
}
EOM

diff --color "/etc/dhcp/dhcpd.conf.bak" "/etc/dhcp/dhcpd.conf" || true

echo ""
echo "Configure usb0 interface..."
cat > /etc/network/interfaces.d/usb0 <<- EOM
auto usb0
allow-hotplug usb0
iface usb0 inet static
  address 192.168.12.1
  netmask 255.255.255.0
  broadcast 192.168.12.255
EOM
cat /etc/network/interfaces.d/usb0

echo ""
echo "Configure wlan0 interface..."
cat > /etc/network/interfaces.d/wlan0 <<- EOM
auto wlan0
allow-hotplug wlan0
iface wlan0 inet static
  address 192.168.14.1
  netmask 255.255.255.0
  broadcast 192.168.14.255
EOM
cat /etc/network/interfaces.d/wlan0

echo ""
echo "Remove initial wlan0 interface definition..."
cp /etc/network/interfaces /etc/network/interfaces.bak
sed -i -n '/# WiFi/q;p' /etc/network/interfaces
diff --color "/etc/network/interfaces.bak" "/etc/network/interfaces" || true

echo ""
echo "Configure hostap daemon (WLAN Hotspot)..."
cat > /etc/hostapd/hostapd.conf <<- EOM
interface=wlan0
driver=nl80211
ssid=CheckInHotspot
hw_mode=g
channel=11
macaddr_acl=0
ignore_broadcast_ssid=0
auth_algs=1
wpa=2
wpa_passphrase=check1n!
wpa_key_mgmt=WPA-PSK
wpa_pairwise=CCMP
wpa_group_rekey=86400
ieee80211n=1
wme_enabled=1
EOM
cat /etc/hostapd/hostapd.conf

echo ""
echo "downloading docker-compose file..."
cd /root
mkdir -p checkin-system
cd checkin-system

curl -O -L https://raw.githubusercontent.com/d-rk/checkin-system/main/docker-compose.yml

cat > .env <<- EOM
API_SECRET=08afe71014fa32e22fa115dc
API_ADMIN_PASSWORD=secret
EOM

echo "pulling docker images..."
docker compose pull

echo "starting docker containers..."
docker compose up -d

echo ""
echo "reboot to finalize the network changes!"
