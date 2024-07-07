#!/usr/bin/env bash

# exit on error
set -e

wlan_config=/etc/network/interfaces.d/wlan0

MODE=client
NEW_MODE=hotspot

if grep -q 192.168.14.1 "$wlan_config"; then
  MODE=hotspot
  NEW_MODE=client
fi

echo "current mode: $MODE"

if [[ "client" == "${MODE}" ]]; then

  echo ""
  echo "Configure wlan0 interface as hotspot..."
  cat > "$wlan_config" <<- EOM
auto wlan0
allow-hotplug wlan0
iface wlan0 inet static
  address 192.168.14.1
  netmask 255.255.255.0
  broadcast 192.168.14.255
EOM
  cat "$wlan_config"

  echo ""
  echo "Enable dhcp server for usb, wlan interfaces..."
  cp /etc/default/isc-dhcp-server /etc/default/isc-dhcp-server.bak
  sed -i -E 's/^#?(INTERFACESv4)=.*/\1="usb0 wlan0"/' "/etc/default/isc-dhcp-server"
  diff --color "/etc/default/isc-dhcp-server.bak" "/etc/default/isc-dhcp-server" || true

  systemctl restart isc-dhcp-server

  echo ""
  echo "Try to bring up interface:"
  ifdown wlan0 || true
  ifup wlan0 || true

  echo ""
  echo "Restart hostapd:"
  systemctl restart hostapd

  exit 0
fi

echo ""
echo "Configure wlan0 interface as client..."
cat > "$wlan_config" <<- EOM
auto wlan0
allow-hotplug wlan0
iface wlan0 inet dhcp
wpa-conf /etc/wpa_supplicant/wpa_supplicant.conf
EOM
cat "$wlan_config"

echo ""
echo "Disable dhcp server for wlan interfaces..."
cp /etc/default/isc-dhcp-server /etc/default/isc-dhcp-server.bak
sed -i -E 's/^#?(INTERFACESv4)=.*/\1="usb0"/' "/etc/default/isc-dhcp-server"
diff --color "/etc/default/isc-dhcp-server.bak" "/etc/default/isc-dhcp-server" || true

systemctl restart isc-dhcp-server

echo ""
echo "Try to bring up interface:"
ifdown wlan0 || true
ifup wlan0 || true

echo ""
echo "Available networks:"
wpa_cli list_networks

read -p "Add an additional network? [y/n] " -n 1 -r
if [[ $REPLY =~ ^[Yy]$ ]]; then
    if [[ -z "${WIFI_SSID}" ]]; then
      read -r -e -p "Enter WIFI SSID: " WIFI_SSID
    fi

    if [[ -z "${WIFI_PASSWORD}" ]]; then
      read -r -e -p "Enter WIFI Password: " WIFI_PASSWORD
    fi

    wpa_passphrase "${WIFI_SSID}" "${WIFI_PASSWORD}" >> /etc/wpa_supplicant/wpa_supplicant.conf
    cat /etc/wpa_supplicant/wpa_supplicant.conf
fi

echo ""
echo "wlan0 now configured as: ${NEW_MODE}"
