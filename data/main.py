import socket
import time
import re
import datetime
import sys

from server import server
from db import Database
from discover import Discovery  # Connection to HP Switches
# Pulling hostname + mac from NAC (Network Access Control)
from api import Clearpass
from switches import SWITCHES_LIST

SWITCHES = SWITCHES_LIST  # importing a list containing switch_name + ip
# list = [
#        (<name>,<ip>),
#        (<name>,<ip>),
#        ]

# db_host = "localhost"
db_host = "db"
db_port = 3306


def check(host, port, timeout=2):
    sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)  # presumably
    sock.settimeout(timeout)
    sock.connect((host, port))
    sock.close()


def cleaning_mac_table(mac_table, regex_filter="Bridge-Aggregation"):
    tmp_list = []
    for entry in mac_table:
        mac = entry["mac"]
        interface_name = entry["interface"]
        vlan = entry["vlan"]
        # Step 1
        # Drop all non-local mac-addresses
        if not re.match(regex_filter, interface_name):
            # Step 2
            # Resolving num of stack and interface
            matches = re.findall(r'\d+', interface_name)
            # Little Hack to enable sorting
            if int(matches[2]) < 10:
                sort_helper = matches[0] + "0" + matches[2]
            else:
                sort_helper = matches[0] + matches[2]
            add = {"mac": mac,  # string
                   "interface_name": interface_name,  # string
                   "vlan": vlan,  # int
                   "stack": matches[0],  # int
                   "interface_num": matches[2],  # int
                   "sort_helper": int(sort_helper)}  # int
            # Step 3
            # Adding the new entry to a temporary list
            tmp_list.append(add)
    # Step 4
    # returning the list sorted by the sort_helper key
    return sorted(tmp_list, key=lambda x: (x["sort_helper"]))


def scanner():
    try:
        try:
            print("Waiting for database for 60 secs...")
            time.sleep(60)  # wait 15 sec for the db to start
        except KeyboardInterrupt:
            print("Skipping wait time!")
        yield "Testing connection to the database..."
        check(db_host, db_port)
    except Exception as e:
        yield (f"Problem with db: {e}")
        return

    yield "Connection to database successful!"
    db = Database(db_host, db_port)
    # Commented out for container-usage
    # db.drop(SWITCHES)
    # db.setup_clearpass_table()
    # db.setup_switch_tables(SWITCHES)

    yield "Connecting to the clearpass api..."
    cp = Clearpass()
    xml = cp.call_api()
    json = cp.convert_to_json(xml)  # constructs json with hostname + mac + ip
    db.truncate("clearpass")
    db.insert_api_data(json)

    max = len(SWITCHES)
    c = 1
    for sw in SWITCHES:
        name = sw[0]
        ip = sw[1]
        yield f"[{c}/{max}] Connecting to {ip} a.k.a {name}..."
        access = Discovery(sw)
        dirty = access.get_mac_table()  # Connecting to a single access switch
        clean = cleaning_mac_table(dirty)
        db.truncate(name)
        db.insert_switch_data(sw, clean)  # Saving clean mac-table
        sw_data = db.select_switch_data(name)
        for entry in sw_data:  # Iterating through every switch
            id = entry[0]
            mac = entry[2]
            mac = mac.replace(":", "")
            hostname = db.select_hostname_by_mac(
                mac)  # Calls the Clearpass table
            ip = db.select_ip_by_mac(mac)
            if hostname is not None:
                db.update_hostname_by_id(hostname[0], id, sw[0])
            if ip is not None:
                db.update_ip_by_id(ip[0], id, sw[0])
            # yield f"{hostname[0] + ip[0]} -> {mac}"
        c += 1

    yield "Scanning finished!"
    now = datetime.datetime.now().strftime("%Y-%m-%d %H:%M:%S")
    yield f"Current date and time : {now}"


def searcher(hostname):
    try:
        check(db_host, db_port)
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
            while True:
                for out in scanner():  # generator object
                    print(out)
                now = datetime.datetime.now().strftime("%Y-%m-%d %H:%M:%S")
                print("Next scan in an hour...")
                time.sleep(60*60)  # re-generate every hour
            break
        elif sys.argv[x] == "--search":
            hostname = sys.argv[x + 1]
            searcher(hostname)
            break
        elif sys.argv[x] == "--server":
            # takes generator object and sends it to client
            server(mock_scanner())
            break
