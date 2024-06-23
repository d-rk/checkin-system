import sys


class RaspiAccess:

    def set_lights(self, green_on, red_on):
        green = "🟩" if green_on else "⬜"
        red = "🟥" if red_on else "⬜"
        print(green + red)

    def set_buzzer(self, on):
        if on:
            print("beep")

    def read_rfid(self):
        for line in sys.stdin:
            if 'Exit' == line.rstrip():
                sys.exit(0)

            return line.rstrip(), None

    def cleanup(self):
        pass
