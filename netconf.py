from napalm import get_network_driver


class Comware7:
    def __init__(self, **kvargs):
        self.hostname = kvargs.get("hostname")
        self.username = kvargs.get("username")
        self.password = kvargs.get("password")

        driver = get_network_driver("h3c_comware")
        driver = driver(self.hostname, self.username, self.password)
        driver.open()
        self.driver = driver
        # ret = driver.get_arp_table()
        # ret = driver.get_mac_address_table()


    def get_mac_table(self):
        return self.driver.get_mac_address_table()

