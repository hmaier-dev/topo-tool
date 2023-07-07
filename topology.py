from netconf import Comware7

class Topology:
    def __init__(self):
        self.username = "exmaple user"
        self.password = "example password"
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
        # ("SW_N2", "192.168.132.144"),
        ("SW_WS", "192.168.132.130"),
        ("SW_KMS", "192.168.132.128"),
        ("SW_P", "192.168.132.136"),
        ("SW_A121", "192.168.132.131"),
    ]

    def pull_mac_tables(self):
        dict = self.SWITCHES
        for switch in dict:
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
                out.write(str(mac))
                out.close()
