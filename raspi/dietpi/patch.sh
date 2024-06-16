#!/usr/bin/env bash

# exit on error
set -e

usage() {
    cat <<EOM
Usage: $(basename $0) <target-folder>
where
  <target-folder> folder where diet-pi is mounted
EOM
    exit 1
}

if [[ $# -eq 0 ]]; then
  usage
fi

TARGET_DIR=$1

if [ ! -d "$TARGET_DIR" ]; then
  echo "$TARGET_DIR does not exist."
  exit 1
fi

# load env variables
if [ -f .env ]; then
   source .env
else
   echo ".env does not exist"
fi

start_patch() {
  file="$1"
  echo "patching $file..."

  if [ ! -f "$file" ]; then
    echo "$file does not exist."
    exit 1
  fi

  # backup
  if [ ! -f "$file.bak" ]; then
    cp "$file" "$file.bak"
  else
    cp "$file.bak" "$file"
  fi
}

finish_patch() {
  file="$1"
  diff --color "$file.bak" "$file" || true
  echo ""
}


# ---
CONFIG_FILE="$TARGET_DIR/config.txt"
start_patch "$CONFIG_FILE"

# enable spi interface
sed -i -E 's/^#?(dtparam=spi)=.*/\1=on/' "$CONFIG_FILE"

# enable i2c interface
sed -i -E 's/^#?(dtparam=i2c_arm)=.*/\1=on/' "$CONFIG_FILE"
sed -i -E 's/^#?(dtparam=i2c_arm_baudrate=100000).*/\1/' "$CONFIG_FILE"

# use dwc2 USB driver (usb network interface)
echo "dtoverlay=dwc2" >> "$CONFIG_FILE"

finish_patch "$CONFIG_FILE"


# ---
CMDLINE_FILE="$TARGET_DIR/cmdline.txt"
start_patch "$CMDLINE_FILE"

# load required modules for usb network interface
sed -i 's/rootwait/rootwait modules-load=dwc2,g_ether/g' "$CMDLINE_FILE"

finish_patch "$CMDLINE_FILE"


# ---
DIETPI_FILE="$TARGET_DIR/dietpi.txt"
start_patch "$DIETPI_FILE"

# keyboard layout
sed -i -E 's/^#?(AUTO_SETUP_KEYBOARD_LAYOUT)=.*/\1=de/' "$DIETPI_FILE"

# timezone
sed -i -E 's/^#?(AUTO_SETUP_TIMEZONE)=.*/\1=Europe\/Berlin/' "$DIETPI_FILE"

# disable ethernet
sed -i -E 's/^#?(AUTO_SETUP_NET_ETHERNET_ENABLED)=.*/\1=0/' "$DIETPI_FILE"

# enable wifi
sed -i -E 's/^#?(AUTO_SETUP_NET_WIFI_ENABLED)=.*/\1=1/' "$DIETPI_FILE"

# wifi country code
sed -i -E 's/^#?(AUTO_SETUP_NET_WIFI_COUNTRY_CODE)=.*/\1=DE/' "$DIETPI_FILE"

# headless
sed -i -E 's/^#?(AUTO_SETUP_HEADLESS)=.*/\1=1/' "$DIETPI_FILE"

# no browser
sed -i -E 's/^#?(AUTO_SETUP_BROWSER_INDEX)=.*/\1=0/' "$DIETPI_FILE"

# enable auto setup
sed -i -E 's/^#?(AUTO_SETUP_AUTOMATED)=.*/\1=1/' "$DIETPI_FILE"

# survey opt-out
sed -i -E 's/^#?(SURVEY_OPTED_IN)=.*/\1=0/' "$DIETPI_FILE"

# hostname
sed -i -E 's/^#?(AUTO_SETUP_NET_HOSTNAME)=.*/\1=CheckInPi/' "$DIETPI_FILE"

# global password
sed -i -E 's/^#?(AUTO_SETUP_GLOBAL_PASSWORD)=.*/\1=check1n!/' "$DIETPI_FILE"

# public key for ssh
DEFAULT_PUBKEY_LOCATION=~/.ssh/id_rsa.pub
read -r -e -p "Copy ssh public key from: " -i "$DEFAULT_PUBKEY_LOCATION" PUBKEY_LOCATION

if [ ! -f "$PUBKEY_LOCATION" ]; then
  echo "$PUBKEY_LOCATION does not exist."
  exit 1
fi

PUBKEY=$(cat "$PUBKEY_LOCATION")
# use ~ instead of / as delimiter for sed to avoid escaping issues
sed -i -E "s~^#?(AUTO_SETUP_SSH_PUBKEY)=.*~\1=$PUBKEY~" "$DIETPI_FILE"

# the SOFTWARE_TO_INSTALL text-block contains two markers to make the replacement repeatable
# delete everything between markers, including markers
sed -i '/### SOFTWARE_START/,/### SOFTWARE_END/d' "$DIETPI_FILE"

SOFTWARE_TO_INSTALL=$(cat <<EOF
### SOFTWARE_START

# install docker
AUTO_SETUP_INSTALL_SOFTWARE_ID=162

# install docker compose
AUTO_SETUP_INSTALL_SOFTWARE_ID=134

### SOFTWARE_END
EOF
)

SOFTWARE_TO_INSTALL_LINE_BREAKS=$(echo "$SOFTWARE_TO_INSTALL" | awk '{printf "%s\\n", $0}')

# install software
sed -i -E "s/^(#?(AUTO_SETUP_INSTALL_SOFTWARE_ID)=.*)/\1\n$SOFTWARE_TO_INSTALL_LINE_BREAKS/" "$DIETPI_FILE"

finish_patch "$DIETPI_FILE"

# ---
DIETPI_WIFI_FILE="$TARGET_DIR/dietpi-wifi.txt"
start_patch "$DIETPI_WIFI_FILE"

if [[ -z "${WIFI_SSID}" ]]; then
  read -r -e -p "Enter WIFI SSID (leave empty to skip): " WIFI_SSID
fi

if [[ -z "${WIFI_PASSWORD}" ]]; then
  read -r -e -p "Enter WIFI Password: " WIFI_PASSWORD
fi

# WIFI SSID
sed -i -E "s/^(aWIFI_SSID\[0\])='.*'/\1='$WIFI_SSID'/" "$DIETPI_WIFI_FILE"

# WIFI password
sed -i -E "s/^(aWIFI_KEY\[0\])='.*'/\1='$WIFI_PASSWORD'/" "$DIETPI_WIFI_FILE"

finish_patch "$DIETPI_WIFI_FILE"


# --- copy additional files
SCRIPT_DIR="$TARGET_DIR/checkin-system"

mkdir -p "$SCRIPT_DIR"
cp post-install.sh "$SCRIPT_DIR/post-install.sh"
echo "copied $SCRIPT_DIR/post-install.sh"

echo "done"