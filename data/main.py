from db import Database
from discover import Discovery  # Connection to HP Switches
from filter import Filtering  # Filtering/Formatting Output from MAC-Tables
from api import Clearpass  # Pulling hostname + mac from NAC (Network Access Control)
import time

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
    ("SW_P", "192.168.132.136"),
    # ("SW_A121", "192.168.132.131"), 1/0/24 set as uplink, no BridgeAggregation filtering possible
]

if __name__ == "__main__":
    # Database-fun
    print("Initializing the database...")
    db = Database()
    db.setup(SWITCHES)
    time.sleep(60)
    while True:
        # Pulling the MAC-Tables
        netconf = Discovery()
        for sw in SWITCHES:
            netconf.get_mac_table(sw)  # Writes to ./switch-tables Directory

        # Cleaning the MAC-Tables
        for switch in SWITCHES:
            f = Filtering(switch)
            f.get_filtered_mac_table()  # Reads from ./switch-tables

        # Pulling hostname + mac Combo
        cp = Clearpass()



