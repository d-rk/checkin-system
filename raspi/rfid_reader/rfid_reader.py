#!/usr/bin/env python

import RPi.GPIO as GPIO
from mfrc522 import SimpleMFRC522
from time import sleep
import requests
import os
import json

GPIO.setmode(GPIO.BCM)

BUZZER_GPIO = 4
LED_GREEN_GPIO = 17
LED_RED_GPIO = 27

GPIO.setup(BUZZER_GPIO, GPIO.OUT)
GPIO.setup(LED_GREEN_GPIO, GPIO.OUT)
GPIO.setup(LED_RED_GPIO, GPIO.OUT)

GPIO.output(LED_GREEN_GPIO, GPIO.LOW)
GPIO.output(LED_RED_GPIO, GPIO.LOW)


API_BASE_URL = os.getenv("API_BASEURL", "http://localhost:8080")


def wait_for_backend():

    healthy = False

    while not healthy:
        r = requests.get(f"{API_BASE_URL}/api/v1/users")

        healthy = r.status_code == 200

        GPIO.output(LED_GREEN_GPIO, healthy)
        GPIO.output(LED_RED_GPIO, not healthy)

        if not healthy:
            print("backend not reachable...")
            sleep(5)


def post_rfid_id(id):
    print("rfid_uid:", id)
    headers = {"Content-type": "application/json"}
    r = requests.post(
        f"{API_BASE_URL}/api/v1/checkins",
        data=json.dumps({"rfid_uid": f"{id}"}),
        headers=headers,
    )
    print(r.status_code)
    print(r.json())

    if r.status_code == 200:
        GPIO.output(BUZZER_GPIO, GPIO.HIGH)
        GPIO.output(LED_GREEN_GPIO, GPIO.HIGH)
        GPIO.output(LED_RED_GPIO, GPIO.LOW)
        sleep(0.3)
    else:
        GPIO.output(BUZZER_GPIO, GPIO.HIGH)
        GPIO.output(LED_GREEN_GPIO, GPIO.LOW)
        GPIO.output(LED_RED_GPIO, GPIO.HIGH)
        sleep(0.3)

    GPIO.output(BUZZER_GPIO, GPIO.LOW)
    GPIO.output(LED_GREEN_GPIO, GPIO.LOW)
    GPIO.output(LED_RED_GPIO, GPIO.LOW)
    sleep(2.0)


wait_for_backend()

print("waiting for card....")

reader = SimpleMFRC522()

try:
    while True:
        id, text = reader.read()
        post_rfid_id(id)
finally:
    GPIO.output(BUZZER_GPIO, GPIO.LOW)
    GPIO.output(LED_GREEN_GPIO, GPIO.LOW)
    GPIO.output(LED_RED_GPIO, GPIO.LOW)
    GPIO.cleanup()
