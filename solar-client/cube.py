
import serial
import requests
import json
import time

def main():

  ser = serial.Serial('/dev/cu.usbserial-A60205DM', 9600)

  time.sleep(2)

  while True:
    r = requests.get('http://do-jo.co/solar/analysis')

    j = json.loads(r.text)

    count = 0
    r = 0
    g = 0
    b = 0

    for e in j:
      er = e['AverageColor']['R']
      eg = e['AverageColor']['G']
      eb = e['AverageColor']['B']

      esum = er + eg + eb

      if esum >= 50:
        count += 1
        r += er
        g += eg
        b += eb

    r = r / count
    g = g / count
    b = b / count

    s = "%d,%d,%d," % (r, g, b)

    print s, count
    #ser.write(s)

    time.sleep(20)

if __name__ == '__main__':
  main()
