import os.path
import subprocess
import requests
import re
from requests import packages
from requests.auth import HTTPBasicAuth
import xml.etree.ElementTree as ET


import cred  # create a cred.py with username = <user> and password = <password>
from discover import Discovery  # Connection to HP Switches
from filter import Filtering  # Filtering/Formatting Output from MAC-Tables

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


def connection_test(switch_array):
    """
    Ping all Switches located in the array
    """
    for switch in switch_array:
        name = switch[0]
        ip = switch[1]
        response = subprocess.run("ping -n 1 " + ip, stdout=subprocess.DEVNULL)
        if response == 0:
            print(f"{name} is up!")
        else:
            print(f"{name} is down!")


def scan_for_hosts(ip_range):
    """
    Scan the given IP address range using Nmap
    """
    nmap_args = ['nmap', '-n', '-sP', ip_range]
    return subprocess.run(nmap_args)


"""
Run functions for the classes
"""


def run_mac_discovery():
    username = input("Enter username:")
    password = input("Enter password:")

    netconf = Discovery(username, password)
    for switch in SWITCHES:
        netconf.get_mac_table(switch)  # Writes to ./switch-tables Directory


def run_mac_table_filter():
    for switch in SWITCHES:
        f = Filtering(switch)
        f.get_filtered_mac_table()  # Reads from ./switch-tables

def call_clearpass_api():
    requests.packages.urllib3.disable_warnings()  # Supress warning for unverified connection to clearpass
    basic = HTTPBasicAuth(cred.username, cred.password)
    # mac = "AC-91-A1-83-B5-2B"
    mac = "0030180891C7"
    resp = requests.get("http://clearpass-cl.ipb-halle.de/tipsapi/config/read/Endpoint", verify=False, auth=basic)
    root = ET.fromstring(resp.content)  # Parsing the retrieved API XML Response
    tree = ET.ElementTree(root)
    tree.write("clearpass-tables/clear_out.xml")

    # for elem in root.iter():
    #     # print(elem.tag, elem.attrib)
    #     m = elem.findall("macAddress")
    #     h = elem.findall("hostname")
    #     print(m, h

def handle_clearpass_table():
    tree = ET.parse("clearpass-tables\clear_out.xml").getroot()
    array = []
    tmp_hostname = None
    tmp_mac = None
    for elem in tree.iter():
        if "hostname" in elem.attrib:
            tmp_hostname = elem.attrib["hostname"]
        if "macAddress" in elem.attrib:
            tmp_mac = elem.attrib["macAddress"]

        if tmp_hostname is not None and tmp_mac is not None:
            add = {
                "hostname": tmp_hostname,
                "mac": tmp_mac
            }
            tmp_hostname = None
            tmp_mac = None
            array.append(add)

    print(array)

if __name__ == "__main__":
    # run_mac_discovery()
    # run_mac_table_filter()
    # call_clearpass_api()
    handle_clearpass_table()
