#!/usr/bin/env bash

# exit on error
set -e

wlan_config=/etc/network/interfaces.d/wlan0

function switch_to_client_mode() {
  echo ""
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

  if ! ip addr show wlan0 | grep -q 'inet '; then
    echo ""
    echo "FAILED: ifup wlan0 failed, fallback to hotspot..."
    switch_to_hotspot_mode
  fi
}

function switch_to_hotspot_mode() {
  echo ""
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
}

function get_current_mode() {
  if grep -q 192.168.14.1 "$wlan_config"; then
    echo hotspot
  else
    echo client
  fi
}

MODE=$(get_current_mode)

echo "current mode: $MODE"
if [[ "${MODE}" == "client" ]]; then
  switch_to_hotspot_mode
else
  switch_to_client_mode
fi

MODE=$(get_current_mode)
echo ""
echo "wlan0 now configured as: ${MODE}"
