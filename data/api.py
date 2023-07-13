import requests
from requests import packages
from requests.auth import HTTPBasicAuth
import xml.etree.ElementTree as ET
import cred  # create a cred.py with username = <user> and password = <password>


class Clearpass:
    def __init__(self):
        self.USER = cred.api_user
        self.PASSWORD = cred.api_password

        self.data = []
        self.call_api()
        self.filter_xml()

    def call_api(self):
        requests.packages.urllib3.disable_warnings()  # Supress warning for unverified connection to clearpass
        basic = HTTPBasicAuth(self.USER, self.PASSWORD)
        resp = requests.get("http://clearpass-cl.ipb-halle.de/tipsapi/config/read/Endpoint", verify=False, auth=basic)
        root = ET.fromstring(resp.content)  # Parsing the retrieved API XML Response
        tree = ET.ElementTree(root)
        tree.write(".\clearpass-tables\clearpass_table.xml")

    def filter_xml(self):
        tree = ET.parse(".\clearpass-tables\clearpass_table.xml").getroot()
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
                # array.append(add)
                self.data.append(add)

