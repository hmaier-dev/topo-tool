import os.path

from netconf import Comware7

'''
Main Class where All the good Stuff Happens
'''


class Discovery:
    def __init__(self, username, password):
        self.username = username
        self.password = password

    def get_mac_table(self, switch):
        """
        Takes tuple with string: name and string: ip
        """
        name = switch[0]
        ip = switch[1]
        # Connecting to the Comware via netconf
        print(f"Connecting to {name} with {ip}...")
        netconf = Comware7(
            hostname=ip,
            username=self.username,
            password=self.password)

        print(f"Pulling MAC-Table from {name}...")
        netconf_output = netconf.get_mac_table()
        # Write mac-address-table to a file
        print(f"Formatting and writing to file...")
        file = "RAW_" + name + ".json"
        path = os.path.join("mac-address-tables", file)
        with open(path, "w") as out:
            mac = str(netconf_output)
            # Formatting the output-string to be useful for json
            mac = mac.replace("\'", "\"")
            mac = mac.replace(" True", "true")
            mac = mac.replace(" False", "false")
            out.write(mac)
            out.close()
