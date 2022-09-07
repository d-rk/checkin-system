# checkin-system

over-engineered checkin-system that can be deployed to a raspberry pi.

![Raspi Board](.github/images/raspi.jpg)

## installation

### run in docker

```shell
docker-compose build
docker-compose up -d
```

The frontend will then be accessible under http://localhost:3000/

### run manually

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
DB_HOST=127.0.0.1
DB_DRIVER=postgres
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=checkin-system
DB_PORT=5432
DB_SSL_MODE=disable

CORS_ALLOWED_ORIGINS=*
EOM
```

3. Run backend

```
cd backend
go run .
```

4. Run frontend

```
cd frontend
npm start
```

## raspi setup

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

## raspi shield

The shield can be build based on the following schematic:

![Fritzing](.github/images/fritzing.png)

Download the [Fritzing File](raspi/fritzing/rfid_reader.fzz).
