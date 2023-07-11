import os.path
# Package: napalm-h3c-comware
from napalm import get_network_driver

'''
Main Class where All the good Stuff Happens
'''


class Discovery:
    def __init__(self,
                 username,
                 password):
        self.username = username
        self.password = password


    def get_mac_table(self, switch):
        """
        Takes tuple with string: name and string: ip
        """
        name = switch[0]
        ip = switch[1]
        print(f"Connecting to {name} with {ip}...")
        driver = get_network_driver("h3c_comware")
        driver = driver(ip, self.username, self.password)
        driver.open()
        print(f"Pulling MAC-Table from {name}...")
        netconf_output = driver.get_mac_address_table()
        # Write mac-address-table to a file
        print(f"Formatting and writing to file...")
        file = "RAW_" + name + ".json"
        path = os.path.join("switch-tables", file)
        with open(path, "w") as out:
            mac = str(netconf_output)
            # Formatting the output-string to be useful for json
            mac = mac.replace("\'", "\"")
            mac = mac.replace(" True", "true")
            mac = mac.replace(" False", "false")
            out.write(mac)
            out.close()
