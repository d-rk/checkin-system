# checkin-system

over-engineered checkin-system that can be deployed to a raspberry pi.

![Raspi Board](.github/images/raspi.jpg)

## Hardware

The system uses a raspberry pi with a custom shield to read RFID tokens
and provide visual and acoustic feedback.

### Building the raspi shield

For the shield the following parts are needed:

#### parts list

| part                       | description                             | link                                        |
| -------------------------- | --------------------------------------- | ------------------------------------------- |
| RC522 RFID Reader module   | module to read ids from cards or tokens | https://www.amazon.de/gp/product/B01M28JAAZ |
| PCF8523 RTC module         | module to record correct timestamps     | https://www.amazon.de/gp/product/B07LGTX8M9 |
| KY-012 active piezo buzzer | acoustic feedback after rfid read       | https://www.amazon.de/gp/product/B07ZYVH6XM |
| green/red LED              | visual feedback after rfid read         |                                             |
| 2 resistors                | needed for LED circuits                 |                                             |

#### Schematic

The shield can be build based on the following schematic:

![Fritzing](.github/images/fritzing.png)

Download the [Fritzing File](raspi/fritzing/rfid_reader.fzz).

### Install and configure the raspi

#### DietPi

The easiest way to install and configure a raspi is by using DietPi.

1. Download DietPi Image from https://dietpi.com/#download
2. Flash `DietPi_RPi-ARMv8-Bookworm.img.xz` to SD card
3. Go to `raspi/dietpi` and patch the image `./patch.sh <folder-where-sd-card-is-mounted-to>`
4. After first boot of the PI, run `/boot/checkin-system/post-install.sh` locally.

After the setup the raspi is available via:
1. USB Ethernet with ip 192.168.12.1
2. WLAN Hotspot with ip 192.168.14.1

#### USB Ethernet on Windows Subsystem for Linux (WSL)

##### Make USB Devices available in WSL

in an administrator command prompt run:

```shell
# list available devices
usbipd list

# share device (change 1-4 to your device) needs to be done once
usbipd bind --busid 1-4

# attach device (change 1-4 to your device) to wsl
usbipd attach --wsl --busid 1-4

# detach device (change 1-4 to your device) to wsl
usbipd detach --busid 1-4
```

##### Configure USB Ethernet in WSL

```shell
# list usb devices
lsusb

# additional interface should show up (normally usb0 something something cryptic)
ip a

# configure interface
sudo ifconfig usb0 192.168.12.2 up

# connect
ssh root@192.168.12.1
```

#### Mounting USB Flash Drives on Windows Subsystem for Linux (WSL)

To use a USB Flash Drive with WSL, you need to connect it to WSL first.
https://devblogs.microsoft.com/commandline/connecting-usb-devices-to-wsl/

Afterward, you can mount it with the following command:

```shell
sudo mkdir /media/usb1
sudo mkdir /media/usb2

lsblk # find the correct devices, e.g. /dev/sde1 /dev/sde2
sudo mount /dev/sde1 /media/usb1
sudo mount /dev/sde2 /media/usb2

# unmount
sudo umount /media/usb1
sudo umount /media/usb2
```

#### Manual

When installing a different distribution on the Raspi, the following steps
describe the steps needed to get everything running.

Install dependencies:

```
sudo apt-get update
sudo apt-get upgrade

sudo apt-get install build-essential git python3-dev python3-pip python3-smbus i2c-tools
sudo pip3 install spidev mfrc522
```

Configure interfaces:

```
sudo raspi-config
```

- [enable spi interface](https://www.raspberrypi-spy.co.uk/2014/08/enabling-the-spi-interface-on-the-raspberry-pi/)
- [enable i2c interface](https://www.raspberrypi-spy.co.uk/2014/11/enabling-the-i2c-interface-on-the-raspberry-pi/)

Install docker:

```
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh

sudo usermod -aG docker $USER
newgrp docker

sudo pip3 install docker-compose
```

## Software

The software consists of the following parts:

- postgres database for persistence
- backend that is attached to the database and provides rest and websocket api
- frontend for the application, which talks via rest,websocket with the backend
- python script to interact with the raspi-shield and send rfid reading to the backend

The software can be installed via [Docker Compose](docker-compose.yml).

### installation

```shell
# clone the repository
git clone https://github.com/d-rk/checkin-system.git

# change to the repo dir
cd checkin-system

# build images yourself
docker-compose build --pull

# or pull them
docker-compose pull

# bring up the containers
docker-compose up -d
```

The frontend will then be accessible under http://localhost:3000/

## Connecting to the RASPI

The raspi can be accessed via ssh or http.
The WI-FI module of the raspi can be in one of two modes: Hotspot or Wlan Client.

After connecting via ssh you can switch to between Hotspot/Client mode with the following command:
```shell
wlan-switch
```

### Via WI-FI (Hotspot)

When the device is in Hotspot mode, connect to the following network:
 - SSID: `CheckInHotspot`
 - Password: `check1n!`

Afterward, connect to the raspi:

```shell
# http
http://192.168.14.1

# ssh (password: see above)
ssh root@192.168.14.1
```

### Via WI-FI (Client)

When the device is in Client mode, you can look up the ip of the raspi on your router.
With this ip you can connect like in Hotspot mode.

### Via USB

You can connect the raspi to a PC using a Micro-USB to USB-A Cable. Make sure to use the
inner one of the two Micro-USB Ports.

For Windows: You have to install a special driver for this to work. see: https://github.com/d-rk/windows_10_raspi_usb_otg_fix 

Afterward, connect to the raspi:

```shell
# http
http://192.168.12.1

# ssh (password: see above)
ssh root@192.168.14.1
```

## Troubleshooting

### Calibrate Hardware Clock

When the hardware clock is out of sync, connect the pi to the internet and update the time as follows:

```shell
# check current time
date
sudo hwclock -r

# update time via internet
sudo ntpdate ntp.ubuntu.com

# set time to hardware clock
sudo hwclock -w
```

## development

this chapter describes steps needed when developing the software.

### local dev environment setup

1. Create a postgres in docker:

```
docker run --name postgres \
    -e POSTGRES_USER=postgres \
    -e POSTGRES_PASSWORD=postgres \
    -e POSTGRES_DB=checkin-system \
    -p 5432:5432 -d postgres
```

2. Create an `.env` file for the backend

```
cat > backend/.env <<- EOM

# sqlite db
#DB_DRIVER=sqlite3
#DB_NAME=checkin.db

# postgres db
DB_HOST=127.0.0.1
DB_DRIVER=postgres
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=checkin-system
DB_PORT=5432
DB_SSL_MODE=disable

CORS_ALLOWED_ORIGINS=*

# days after which checkIn will be deleted
CHECKIN_RETENTION_DAYS=100

# secret to sign bearer tokens with
API_SECRET=yoursecretstring

# token expiry duration
TOKEN_EXPIRY_MINUTES=60

# password for initial admin account
ADMIN_PASSWORD=secret
EOM
```

3. Run backend

```
cd backend
go run ./...
```

4. Create an `.env` file for the frontend

```
cat > frontend/.env <<- EOM
VITE_API_BASE_URL=http://localhost:8080
# optional admin credentials for auto login 
VITE_API_USER=admin
VITE_API_PASSWORD=secret
EOM
```

5. Run frontend

```
cd frontend
npm start
```

### simplify working with different raspis

To simplify connecting to different raspis via USB and avoiding host key check
errors add the following to your `~/.ssh/config`:

```shell
Host checkin
    HostName 192.168.12.1
    User root
    StrictHostKeyChecking no
    UserKnownHostsFile /dev/null
    Port 22
```

And then just use: `ssh checkin`

### build/publish docker images

#### publish latest image for development

```shell
# login to quay.io / user-settings / Generate Encrypted Password
export QUAY_IO_PASSWORD=xxx
echo $QUAY_IO_PASSWORD | docker login -u d_rk --password-stdin quay.io

# build and push latest
GIT_COMMIT=$(git rev-parse --short HEAD) docker buildx bake --push

# build and push latest and backend only
GIT_COMMIT=$(git rev-parse --short HEAD) docker buildx bake backend --push
```

#### publish release

```shell
# login to quay.io / user-settings / Generate Encrypted Password
export QUAY_IO_PASSWORD=xxx
echo $QUAY_IO_PASSWORD | docker login -u d_rk --password-stdin quay.io

# build and push with tag
git tag v3.0.0 # create tag

# build and push versioned docker images
COMMIT=$(git rev-parse --short HEAD)
TAG=$(git describe --tags --abbrev=0)
MAJOR="${TAG%%.*}"           # v3
MINOR="${TAG%.*}"            # v3.3
PATCH="$TAG"                 # v3.3.3

VERSION=$PATCH GIT_COMMIT=$COMMIT docker buildx bake --push

docker tag quay.io/d_rk/checkin-system-backend:$PATCH quay.io/d_rk/checkin-system-backend:$MINOR
docker tag quay.io/d_rk/checkin-system-backend:$PATCH quay.io/d_rk/checkin-system-backend:$MAJOR
docker tag quay.io/d_rk/checkin-system-frontend:$PATCH quay.io/d_rk/checkin-system-frontend:$MINOR
docker tag quay.io/d_rk/checkin-system-frontend:$PATCH quay.io/d_rk/checkin-system-frontend:$MAJOR
docker tag quay.io/d_rk/checkin-system-reader:$PATCH quay.io/d_rk/checkin-system-reader:$MINOR
docker tag quay.io/d_rk/checkin-system-reader:$PATCH quay.io/d_rk/checkin-system-reader:$MAJOR

docker push quay.io/d_rk/checkin-system-backend:$MINOR
docker push quay.io/d_rk/checkin-system-backend:$MAJOR
docker push quay.io/d_rk/checkin-system-frontend:$MINOR
docker push quay.io/d_rk/checkin-system-frontend:$MAJOR
docker push quay.io/d_rk/checkin-system-reader:$MINOR
docker push quay.io/d_rk/checkin-system-reader:$MAJOR

# push git tags
git push --tags
```
