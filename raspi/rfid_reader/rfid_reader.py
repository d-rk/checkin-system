#!/usr/bin/env python
#
import RPi.GPIO as GPIO
from mfrc522 import SimpleMFRC522
import json
import os
import traceback
from time import sleep

import requests

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
API_USER = os.getenv("API_USER")
API_PASSWORD = os.getenv("API_PASSWORD")

session = requests.session()
session.headers = {"Content-type": "application/json"}


def requests_call(method, url, **kwargs):
    try:
        response = session.request(method, url, **kwargs)
    except BaseException as exception:
        # anticipate giant data string: curtail for logging purposes
        if 'data' in kwargs and len(kwargs['data']) > 500:
            kwargs['data'] = f'{kwargs["data"][:500]}...'
        print(f'request: {method.upper()} {url} {kwargs}')
        print(f'exception: {exception}')
        raw_tb = traceback.extract_stack()
        msg = 'Stack trace:\n' + ''.join(traceback.format_list(raw_tb[:-1]))
        print(msg)
        return False, exception
    return True, response


def parse_json(response):
    try:
        return response.json()
    except BaseException as ex:
        print(f"no valid json: {response.text}")
        raise ex


def get_token():
    success, response = requests_call('post', f"{API_BASE_URL}/api/login",
                                      data=json.dumps({"username": f"{API_USER}", "password": f"{API_PASSWORD}"})
                                      )

    if success and response.status_code == 200:
        return parse_json(response)["token"]
    elif success:
        print("login failed", "status", response.status_code, "response", response.json())
        return None
    else:
        print("login failed")
        return None


def refresh_token(response, *args, **kwargs):
    if response.status_code == 401:
        print("Fetching new token as the previous token expired")
        token = get_token()
        session.headers.update({"Authorization": f"Bearer {token}"})
        response.request.headers["Authorization"] = session.headers["Authorization"]
        return session.send(response.request, verify=False)


session.hooks['response'].append(refresh_token)


def wait_for_backend():
    healthy = False

    while not healthy:
        success, response = requests_call('get', f"{API_BASE_URL}/api/v1/users/me")

        healthy = success and response.status_code == 200

        GPIO.output(LED_GREEN_GPIO, healthy)
        GPIO.output(LED_RED_GPIO, not healthy)

        if not healthy:
            print("backend not reachable...")
            sleep(5)


def post_rfid_id(id):
    print("rfid_uid:", id)
    headers = {"Content-type": "application/json"}
    success, response = requests_call('post', f"{API_BASE_URL}/api/v1/checkins?rfid={id}", headers=headers)

    if success:
        print(response.status_code)
        print(parse_json(response))

    if success and response.status_code == 200:
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
