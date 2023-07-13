import os.path
import requests
from requests import packages
from requests.auth import HTTPBasicAuth
import xml.etree.ElementTree as ET
import cred  # create a cred.py with username = <user> and password = <password>


class Clearpass:
    def __init__(self, call_api=True):
        self.USER = cred.api_user
        self.PASSWORD = cred.api_password
        self.table_dir = "clearpass-tables"

        if not os.path.exists(self.table_dir):
            os.makedirs(self.table_dir)

        self.data = []
        self.path_xml = os.path.join(self.table_dir, "clearpass_table.xml")
        self.path_json = os.path.join(self.table_dir, "clearpass_table.json")
        if call_api:
            self.call_api()
        self.filter_xml()


    def call_api(self):
        requests.packages.urllib3.disable_warnings()  # Supress warning for unverified connection to clearpass
        basic = HTTPBasicAuth(self.USER, self.PASSWORD)
        print("Requesting data from clearpass...")
        resp = requests.get("http://clearpass-cl.ipb-halle.de/tipsapi/config/read/Endpoint", verify=False, auth=basic)
        root = ET.fromstring(resp.content)  # Parsing the retrieved API XML Response
        tree = ET.ElementTree(root)
        print("Write data to xml...")
        tree.write(self.path_xml)

    def filter_xml(self):
        print("Parsing xml data...")
        tree = ET.parse(self.path_xml).getroot()
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

        with open(self.path_json, "w") as file:
            file.write(str(self.data))
            file.close()
