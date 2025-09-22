#!/usr/bin/env bash

# WiFi Network Manager for wpa_supplicant.conf
# Provides functions to add, update, remove, and manage WiFi networks

set -e

WPA_CONFIG="${WPA_CONFIG:-/etc/wpa_supplicant/wpa_supplicant.conf}"
INTERFACE_CONFIG=/etc/network/interfaces.d/wlan0

function usage() {
    echo "Usage: $0 [COMMAND] [OPTIONS]"
    echo ""
    echo "Commands:"
    echo "  add SSID PASSWORD    Add network (moves to top priority)"
    echo "  remove SSID          Remove network by SSID"
    echo "  list                 List all configured networks"
    echo "  mode                 Show current WiFi mode (client or hotspot)"
    echo "  toggle-mode          Toggle between client and hotspot mode"
    echo "  help                 Show this help message"
    echo ""
    echo "Examples:"
    echo "  $0 add \"MyWiFi\" \"mypassword\""
    echo "  $0 remove \"OldNetwork\""
    echo "  $0 list"
}

function ensure_config_exists() {
    if [[ ! -f "$WPA_CONFIG" ]]; then
        echo "Creating basic wpa_supplicant.conf..."
        cat > "$WPA_CONFIG" <<EOF
country=DE
ctrl_interface=DIR=/var/run/wpa_supplicant GROUP=netdev

EOF
    fi
}

function generate_network_block() {
    local ssid="$1"
    local password="$2"

    if [[ -z "$ssid" ]]; then
        echo "Error: SSID cannot be empty" >&2
        return 1
    fi

    if [[ -z "$password" ]]; then
        echo "Error: Password cannot be empty" >&2
        return 1
    fi

    # Use wpa_passphrase to generate the PSK
    wpa_passphrase "$ssid" "$password" | grep -v "^[[:space:]]*#psk="
}

function remove_network() {
    local target_ssid="$1"

    if [[ -z "$target_ssid" ]]; then
        echo "Error: SSID cannot be empty" >&2
        return 1
    fi

    if [[ ! -f "$WPA_CONFIG" ]]; then
        echo "Error: $WPA_CONFIG does not exist" >&2
        return 1
    fi

    # Create a temporary file to build the new config
    local temp_file=$(mktemp)
    local in_target_network=false
    local network_found=false

    while IFS= read -r line; do
        if [[ "$line" =~ ^network=\{ ]]; then
            in_target_network=true
            current_network_block="$line"$'\n'
            continue
        fi

        if [[ "$in_target_network" == true ]]; then
            current_network_block+="$line"$'\n'

            if [[ "$line" =~ ^[[:space:]]*ssid=\"(.*)\" ]]; then
                current_ssid="${BASH_REMATCH[1]}"
            fi

            if [[ "$line" =~ ^\} ]]; then
                in_target_network=false

                if [[ "$current_ssid" != "$target_ssid" ]]; then
                    # Keep this network
                    echo -n "$current_network_block" >> "$temp_file"
                else
                    # Skip this network (remove it)
                    network_found=true
                    echo "Removed network: $target_ssid"
                fi

                current_network_block=""
                current_ssid=""
                continue
            fi
        else
            # Not in a network block, copy line as-is
            echo "$line" >> "$temp_file"
        fi
    done < "$WPA_CONFIG"

    if [[ "$network_found" == false ]]; then
        echo "Warning: Network '$target_ssid' not found"
        rm "$temp_file"
        return 1
    fi

    # Replace the original config with the new one
    mv "$temp_file" "$WPA_CONFIG"
    echo "Network '$target_ssid' removed successfully"
    return 0
}

function add_network() {
    local target_ssid="$1"
    local password="$2"

    if [[ -z "$target_ssid" ]]; then
        echo "Error: SSID cannot be empty" >&2
        return 1
    fi

    if [[ -z "$password" ]]; then
        echo "Error: Password cannot be empty" >&2
        return 1
    fi

    ensure_config_exists

    # Generate the new network block
    local new_network_block
    new_network_block=$(generate_network_block "$target_ssid" "$password")

    if [[ $? -ne 0 ]]; then
        echo "Error: Failed to generate network block" >&2
        return 1
    fi

    # Create a temporary file to build the new config
    local temp_file=$(mktemp)
    local network_inserted=false

    while IFS= read -r line; do
        # Check if this is the first network block
        if [[ "$line" =~ ^network=\{ ]] && [[ "$network_inserted" == false ]]; then
            # Insert new network before this existing network
            echo "$new_network_block" >> "$temp_file"
            echo "" >> "$temp_file"
            network_inserted=true
        fi

        # Copy the current line
        echo "$line" >> "$temp_file"
    done < "$WPA_CONFIG"

    # If we never found a network block, add the new network at the end
    if [[ "$network_inserted" == false ]]; then
        echo "" >> "$temp_file"
        echo "$new_network_block" >> "$temp_file"
    fi

    # Replace the original config with the new one
    mv "$temp_file" "$WPA_CONFIG"
    echo "Added new network '$target_ssid' with top priority"
}

function list_networks() {
    if [[ ! -f "$WPA_CONFIG" ]]; then
        echo "No wpa_supplicant.conf file found at ($WPA_CONFIG)"
        return 1
    fi

    local network_count=0
    local in_network=false
    local current_ssid=""

    while IFS= read -r line; do
        if [[ "$line" =~ ^network=\{ ]]; then
            in_network=true
            current_ssid=""
            continue
        fi

        if [[ "$in_network" == true ]]; then
            if [[ "$line" =~ ^[[:space:]]*ssid=\"(.*)\" ]]; then
                current_ssid="${BASH_REMATCH[1]}"
            fi

            if [[ "$line" =~ ^\} ]]; then
                in_network=false
                if [[ -n "$current_ssid" ]]; then
                    echo "$current_ssid"
                    network_count=$((network_count + 1))
                fi
            fi
        fi
    done < "$WPA_CONFIG"

    if [[ $network_count -eq 0 ]]; then
        echo "No networks configured"
    fi
}

function switch_to_client_mode() {
  echo ""
  echo ""
  echo "Configure wlan0 interface as client..."
  cat > "$INTERFACE_CONFIG" <<- EOM
auto wlan0
allow-hotplug wlan0
iface wlan0 inet dhcp
wpa-conf /etc/wpa_supplicant/wpa_supplicant.conf
EOM
  cat "$INTERFACE_CONFIG"

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
  cat > "$INTERFACE_CONFIG" <<- EOM
auto wlan0
allow-hotplug wlan0
iface wlan0 inet static
  address 192.168.14.1
  netmask 255.255.255.0
  broadcast 192.168.14.255
EOM
  cat "$INTERFACE_CONFIG"

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

function get_wifi_mode() {
  if grep -q 192.168.14.1 "$INTERFACE_CONFIG"; then
    echo hotspot
  else
    echo client
  fi
}

function toggle_wifi_mode() {
  local MODE=$(get_wifi_mode)

  echo "current mode: $MODE"
  if [[ "${MODE}" == "client" ]]; then
    switch_to_hotspot_mode
  else
    switch_to_client_mode
  fi

  MODE=$(get_wifi_mode)
  echo ""
  echo "wlan0 now configured as: ${MODE}"
}

# Main script logic
case "${1:-}" in
    "add")
        if [[ $# -ne 3 ]]; then
            echo "Error: add command requires SSID and PASSWORD"
            echo "Usage: $0 add SSID PASSWORD"
            exit 1
        fi
        add_network "$2" "$3"
        ;;
    "remove")
        if [[ $# -ne 2 ]]; then
            echo "Error: remove command requires SSID"
            echo "Usage: $0 remove SSID"
            exit 1
        fi
        remove_network "$2"
        ;;
    "list")
        list_networks
        ;;
    "mode")
        get_wifi_mode
        ;;
    "toggle-mode")
        toggle_wifi_mode
        ;;
    "help"|"--help"|"-h")
        usage
        ;;
    "")
        usage
        exit 1
        ;;
    *)
        echo "Error: Unknown command '$1'"
        echo ""
        usage
        exit 1
        ;;
esac
