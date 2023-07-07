from netconf import Comware7


class Topology:
    def __init__(self, username, password):
        self.username = username
        self.password = password
        self.SWITCHES = [
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

    def pull_mac_tables(self):
        sw_dict = self.SWITCHES
        for switch in sw_dict:
            name = switch[0]
            ip = switch[1]
            # Connecting to the Comware via netconf
            print(f"Connecting to {name} with {ip}...")
            netconf = Comware7(
                hostname=ip,
                username=self.username,
                password=self.password)

            mac = netconf.get_mac_table()
            # Write mac-address-table to a file
            with open(f"./mac-address-tables/{name}.txt", "w") as out:
                mac = str(mac)
                # Formatting the output-string to be useful for json
                mac = mac.replace("\'", "\"")
                mac = mac.replace(" True", "true")
                mac = mac.replace(" False", "false")
                out.write(mac)
                out.close()
