import json
import os
import re

'''
This class should filter all unnecessary entrys as well, as sort the useful ones
'''


class Filtering:
    def __init__(self, file_name):
        """
        Logic is processed in the order of the following lists
        """
        self.noBridgeAggregation = []
        self.withStackAndInterface = []
        self.sortedByInterface = []

        self.table_name = os.path.join("./mac-address-tables", file_name)
        # Opening the desired file
        with open(self.table_name, "r") as file:
            self.table = json.loads(file.read())
            file.close()

    def filtering_mac_tables(self):
        for entry in self.table:
            mac = entry["mac"]
            interface = entry["interface"]
            vlan = entry["vlan"]

            if interface != "Bridge-Aggregation1":
                add = {"mac": mac,
                       "interface": interface,
                       "vlan": vlan}
                self.noBridgeAggregation.append(add)

    def get_stack_and_interface(self):
        for entry in self.noBridgeAggregation:
            mac = entry["mac"]
            interface_name = entry["interface"]
            vlan = entry["vlan"]
            matches = re.findall(r'\d+', interface_name)
            stack = matches[0]
            interface_num = matches[2]
            add = {"mac": mac,
                   "interface_name": interface_name,
                   "vlan": vlan,
                   "stack": stack,
                   "interface_num": matches[2]}

    def sort_by_interface(self):
        for entry in self.withStackAndInterface:
            mac = entry["mac"]
            interface = entry["interface"]
            vlan = entry["vlan"]


f = Filtering("SW_A-Sued.txt")
f.filtering_mac_tables()  # Reads from ./mac-address-tables
f.get_stack_and_interface()  #
