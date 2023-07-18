import os.path
import requests
from requests import packages
from requests.auth import HTTPBasicAuth
import xml.etree.ElementTree as ET
import cred


class Clearpass:
    def __init__(self, call_api=True):
        self.USER = cred.api_user
        self.PASSWORD = cred.api_password
        self.table_dir = "clearpass-tables"
        self.data = []

    def call_api(self):
        # Supress warning for unverified connection to clearpass
        requests.packages.urllib3.disable_warnings()
        basic = HTTPBasicAuth(self.USER, self.PASSWORD)
        print("Requesting data from clearpass...")
        resp = requests.get(
            "http://clearpass-cl.ipb-halle.de/tipsapi/config/read/Endpoint", verify=False, auth=basic)
        # Parsing the retrieved API XML Response
        root = ET.fromstring(resp.content)
        tree = ET.ElementTree(root)
        return tree

    def convert_to_json(self, tree):
        print("Parsing xml data...")
        # tree = ET.parse(self.path_xml).getroot()
        tmp_hostname = None
        tmp_mac = None
        print("Filtering xml data into json...")
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
                # array.append(add)
                self.data.append(add)
        return self.data

