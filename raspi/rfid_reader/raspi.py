import RPi.GPIO as GPIO
from mfrc522 import SimpleMFRC522


class RaspiAccess:

    BUZZER_GPIO = 4
    LED_GREEN_GPIO = 17
    LED_RED_GPIO = 27

    reader = SimpleMFRC522()

    def __init__(self):
        GPIO.setmode(GPIO.BCM)

        GPIO.setup(self.BUZZER_GPIO, GPIO.OUT)
        GPIO.setup(self.LED_GREEN_GPIO, GPIO.OUT)
        GPIO.setup(self.LED_RED_GPIO, GPIO.OUT)

        GPIO.output(self.LED_GREEN_GPIO, GPIO.LOW)
        GPIO.output(self.LED_RED_GPIO, GPIO.LOW)

    def set_lights(self, green_on, red_on):
        GPIO.output(self.LED_GREEN_GPIO, green_on)
        GPIO.output(self.LED_RED_GPIO, red_on)

    def set_buzzer(self, on):
        GPIO.output(self.BUZZER_GPIO, on)

    def read_rfid(self):
        return self.reader.read()

    def cleanup(self):
        GPIO.cleanup()
