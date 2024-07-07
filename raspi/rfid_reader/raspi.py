import RPi.GPIO as GPIO
from mfrc522 import SimpleMFRC522

MODE = GPIO.BCM

BCM_BUZZER_GPIO = 4
BCM_LED_GREEN_GPIO = 17
BCM_LED_RED_GPIO = 27

BOARD_BUZZER_PIN = 7
BOARD_LED_GREEN_PIN = 11
BOARD_LED_RED_PIN = 13


class RaspiAccess:

    reader = None
    buzzer_pin = 0
    led_green_pin = 0
    led_red_pin = 0

    def __init__(self):
        GPIO.setmode(MODE)
        GPIO.setwarnings(False)

        if MODE == GPIO.BCM:
            self.buzzer_pin = BCM_BUZZER_GPIO
            self.led_green_pin = BCM_LED_GREEN_GPIO
            self.led_red_pin = BCM_LED_RED_GPIO
        else:
            self.buzzer_pin = BOARD_BUZZER_PIN
            self.led_green_pin = BOARD_LED_GREEN_PIN
            self.led_red_pin = BOARD_LED_RED_PIN

        self.reader = SimpleMFRC522()

        GPIO.setup(self.buzzer_pin, GPIO.OUT)
        GPIO.setup(self.led_green_pin, GPIO.OUT)
        GPIO.setup(self.led_red_pin, GPIO.OUT)

        GPIO.output(self.led_green_pin, GPIO.LOW)
        GPIO.output(self.led_red_pin, GPIO.LOW)

    def set_lights(self, green_on, red_on):
        GPIO.output(self.led_green_pin, green_on)
        GPIO.output(self.led_red_pin, red_on)

    def set_buzzer(self, on):
        GPIO.output(self.buzzer_pin, on)

    def read_rfid(self):
        return self.reader.read()

    def cleanup(self):
        GPIO.cleanup()
