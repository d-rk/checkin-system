#!/usr/bin/env python
#
import json
import sys
import os
import traceback
from time import sleep

import requests

from raspi import RaspiAccess
#from raspi_dummy import RaspiAccess

raspi = RaspiAccess()

API_BASE_URL = os.getenv("API_BASEURL", "http://localhost:8080")
API_USER = os.getenv("API_USER")
API_PASSWORD = os.getenv("API_PASSWORD")

if API_USER is None:
    print("env variable missing: API_USER")
    sys.exit(-2)

if API_PASSWORD is None:
    print("env variable missing: API_PASSWORD")
    sys.exit(-2)

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
    if not response.request.url.endswith("/api/login") and response.status_code == 401:
        print("Fetching new token as the previous token expired")
        token = get_token()

        if token is None:
            raise Exception("unable to refresh token")
        else:
            session.headers.update({"Authorization": f"Bearer {token}"})
            response.request.headers["Authorization"] = session.headers["Authorization"]
            return session.send(response.request, verify=False)


session.hooks['response'].append(refresh_token)


def wait_for_backend():
    healthy = False

    while not healthy:
        success, response = requests_call('get', f"{API_BASE_URL}/api/v1/users/me")

        healthy = success and response.status_code == 200

        raspi.set_lights(healthy, not healthy)

        if not healthy:
            print("backend not reachable...")
            sleep(5)


def show_result(success):

    raspi.set_buzzer(True)
    time = 0
    increment = 50

    while time < 600:
        toggle = (time % 100) == 0
        raspi.set_lights(success and toggle, (not success) and toggle)
        time += increment
        sleep(increment / 1000.0)

        if success and time > 150:
            raspi.set_buzzer(False)


def post_rfid_id(id):
    print("rfid_uid:", id)
    headers = {"Content-type": "application/json"}
    success, response = requests_call('post', f"{API_BASE_URL}/api/v1/checkins?rfid={id}", headers=headers)

    if success:
        print(response.status_code)
        print(parse_json(response))

    if success and (response.status_code == 200 or response.status_code == 201):
        show_result(True)
    else:
        show_result(False)

    raspi.set_buzzer(False)
    raspi.set_lights(True, False)
    sleep(2.0)


print("waiting for backend....")

wait_for_backend()

print("waiting for card....")

try:
    while True:
        rfid_id, text = raspi.read_rfid()
        post_rfid_id(rfid_id)
finally:
    raspi.set_buzzer(False)
    raspi.set_lights(False, False)
    raspi.cleanup()
