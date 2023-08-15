import socket
import time
import datetime
import sys

from server import server
from db import Database
from switches import SWITCHES_LIST
from scanner import Scanner

SWITCHES = SWITCHES_LIST  # importing a list containing switch_name + ip
# list = [
#        (<name>,<ip>),
#        (<name>,<ip>),
#        ]


def searcher(hostname):
    db_host = "localhost"
    db_port = 3306
    try:
        sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)  # presumably
        sock.settimeout(2)
        sock.connect((db_host, db_port))
        sock.close()
    except Exception as e:
        print(f"Problem with db: {e}")
        return

    db = Database(db_host, db_port)
    SWITCHES = db.show_switch_tables()
    count = 0
    for sw in SWITCHES:
        resp = db.select_entry_by_hostname(
            sw, hostname)  # BEWARE FOR SQL INJECTION
        if resp:
            for output in resp:
                # Building a nice dict from this
                count += 1
                cols = db.column_name_of(sw)
                out = prettify_response(cols, output)
                print("----------------------------------------------")
                print(f"Found on: {sw}\n", out)
                print("----------------------------------------------")
    if count <= 0:
        print("Not found... :(")


# Just for the eyes!
def prettify_response(col, data):
    out = {}
    for index, entry in enumerate(col):
        c = col[index][0]
        v = data[index]
        out[c] = v
    return out


def mock_scanner():
    time.sleep(1)
    yield "Testing connection to the database..."
    time.sleep(1)
    yield "Connection to database successful!"
    time.sleep(1)
    yield "Setup tables..."
    time.sleep(1)
    yield "Connecting to the clearpass api..."
    time.sleep(1)
    yield "Scanning switches!"
    time.sleep(1)
    yield "Scanning finished!"
    time.sleep(1)
    now = datetime.datetime.now().strftime("%Y-%m-%d %H:%M:%S")
    yield f"Current date and time : {now}"


if __name__ == "__main__":
    y = len(sys.argv)
    for x in range(1, y):
        if sys.argv[x] == "--scanner":
            scan = Scanner(SWITCHES_LIST).start()
            while True:
                for out in scan:  # generator object
                    print(out)
                print("Next scan in an hour...")
                try:
                    time.sleep(60*60)  # re-generate every hour
                except KeyboardInterrupt:
                    break
            break
        elif sys.argv[x] == "--query-switches":
            Scanner(SWITCHES_LIST).query_switches()
            break
        elif sys.argv[x] == "--query-clearpass":
            Scanner(SWITCHES_LIST).query_clearpass()
            break
        elif sys.argv[x] == "--search":
            hostname = sys.argv[x + 1]
            searcher(hostname)
            break
        elif sys.argv[x] == "--server":
            # takes generator object and sends it to client
            server(mock_scanner())
            break
