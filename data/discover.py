# Package: napalm-h3c-comware
from napalm import get_network_driver

import os.path
import cred

'''
Main Class where All the good Stuff Happens
'''


class Discovery:
    def __init__(self, switch):
        self.USERNAME = cred.sw_user
        self.PASSWORD = cred.sw_password
        self.table_dir = "switch-tables"
        self.mac_table = []
        # Name and IP of Switch
        try:
            self.name = switch[0]
            self.ip = switch[1]
        except Exception as e:
            print(f"No name or ip given: {e}")
        #
        # if not os.path.exists(self.table_dir):
        #     os.makedirs(self.table_dir)
        self.driver = get_network_driver("h3c_comware")
        self.driver = self.driver(self.ip, self.USERNAME, self.PASSWORD, timeout=5)
        self.driver.open()

    def get_mac_table(self):
        """
        Takes tuple with string: name and string: ip
        """
        mac = []
        try:
            mac = self.driver.get_mac_address_table()
        except Exception as e:
            print(f"Problems access switch: {e}")
        return mac

    def get_arp_table(self):
        """
        Takes tuple with string: name and string: ip
        """
        arp = []
        try:
            arp = self.driver.get_arp_table()
        except Exception as e:
            print(f"Problems pulling the arp table: {e}")
        return arp


        # mac = str(mac)
        # # Formatting the output-string to be useful for json
        # mac = mac.replace("\'", "\"")
        # mac = mac.replace(" True", "true")
        # mac = mac.replace(" False", "false")
        # # self.write_to_file()
        # print("-------------------------------------------------------")
        # return mac

    def write_to_file(self):
        print("Formatting and writing to file...")
        file = "RAW_" + self.name + ".json"
        path = os.path.join(self.table_dir, file)
        with open(path, "w") as out:
            mac = str(self.mac_table)
            # Formatting the output-string to be useful for json
            mac = mac.replace("\'", "\"")
            mac = mac.replace(" True", "true")
            mac = mac.replace(" False", "false")
            out.write(mac)
            out.close()
