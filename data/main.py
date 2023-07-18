import socket
import time
import re
import datetime

from db import Database
from discover import Discovery  # Connection to HP Switches
# Pulling hostname + mac from NAC (Network Access Control)
from api import Clearpass

SWITCHES = [
    ("SW_A-Nord", "192.168.132.125"),
    ("SW_A-Sued", "192.168.132.126"),
    ("SW_C-Nord-Core", "192.168.132.120"),
    ("SW_C-Nord", "192.168.132.121"),
    ("SW_C-Sued", "192.168.132.122"),
    ("SW_B", "192.168.132.132"),
    ("SW_R", "192.168.132.141"),
    ("SW_R2", "192.168.132.142"),
    ("SW_D", "192.168.132.127"),
    ("SW_E1", "192.168.132.138"),
    ("SW_E2", "192.168.132.123"),
    ("SW_F", "192.168.132.202"),
    ("SW_G", "192.168.132.146"),
    ("SW_N1", "192.168.132.134"),
    # ("SW_N2", "192.168.132.144"), currently not reachable
    ("SW_WS", "192.168.132.130"),
    ("SW_KMS", "192.168.132.128"),
    # ("SW_P", "192.168.132.136"), 1/0/28 set as uplink
    # ("SW_A121", "192.168.132.131"), 1/0/24 set as uplink
]

db_host = "localhost"
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
            sort_helper = matches[0] + matches[2]  # combine two strings
            add = {"mac": mac,                          # string
                   "interface_name": interface_name,    # string
                   "vlan": vlan,                        # int
                   "stack": matches[0],                 # int
                   "interface_num": matches[2],         # int
                   "sort_helper": int(sort_helper)}     # int
            # Step 3
            # Adding the new entry to a temporary list
            tmp_list.append(add)
    # Step 4
    # returning the list sorted by the sort_helper key
    return sorted(tmp_list, key=lambda x: (x["sort_helper"]))


def main():
    try:
        print("Testing connection to the database...")
        check(db_host, db_port)
    except Exception as e:
        print(f"Problem with db: {e}")
        return

    print("Connection to database successful!")
    db = Database(db_host, db_port)
    print("Cleaning the database tables...")
    # db.truncate(SWITCHES)

    # Pulling the MAC-Tables
    # access = Discovery()
    # max = len(SWITCHES)
    # c = 1
    # for sw in SWITCHES:
    #     print(f"[{c}/{max}] Connecting to {sw[1]} with {sw[0]}...")
    #     dirty = access.get_mac_table(sw)
    #     clean = cleaning_mac_table(dirty)
    #     db.insert_switch_data(sw, clean)
    #     c += 1

# for sw in SWITCHES:
    # Pulling hostname + mac Combo
    db.setup_clearpass_table()
    cp = Clearpass()
    xml = cp.call_api()
    json = cp.convert_to_json(xml)
    db.insert_api_data(json)

    print("Scanning finished!")
    now = datetime.datetime.now().strftime("%Y-%m-%d %H:%M:%S")
    print(f"Current date and time : {now}")


if __name__ == "__main__":
    main()
