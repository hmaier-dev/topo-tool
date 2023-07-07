import json
import os
import re

'''
This class should filter all unnecessary entrys as well, as sort the useful ones
'''


class Filtering:
    def __init__(self, switch):
        """
        Logic is processed in the order of the following lists
        """
        self.noBridgeAggregation = []
        self.withStackAndInterface = []
        self.sortedByInterface = []

        self.switch_name = switch[0] + "-clean.json"
        self.switch_ip = switch[1]

        self.table_name = os.path.join("./mac-address-tables", self.switch_name)
        # Opening the desired file
        with open(self.table_name, "r") as file:
            self.raw_mac_table = json.loads(file.read())
            file.close()

    # Step 1
    def filtering_mac_tables(self):
        for entry in self.raw_mac_table:
            mac = entry["mac"]
            interface = entry["interface"]
            vlan = entry["vlan"]
            if not re.match("Bridge-Aggregation",interface):
                add = {"mac": mac,
                       "interface": interface,
                       "vlan": vlan}
                self.noBridgeAggregation.append(add)

    # Step 2
    def resolve_stack_and_interface(self):
        for entry in self.noBridgeAggregation:
            mac = entry["mac"]
            interface_name = entry["interface"]
            vlan = entry["vlan"]
            matches = re.findall(r'\d+', interface_name)
            stack = matches[0]
            interface_num = matches[2]

            helper = stack + interface_num
            add = {"mac": mac,
                   "interface_name": interface_name,
                   "vlan": vlan,
                   "stack": stack,
                   "interface_num": interface_num,
                   "sort_helper": int(helper)}
            self.withStackAndInterface.append(add)

    # Step 3
    def sort_by_interface(self):
        # solution found on GitHub... do not really know how sorted + lambda works (?!)
        self.sortedByInterface = sorted(self.withStackAndInterface, key=lambda x: (x["sort_helper"]))

    def write_to_mac_directory(self):
        name = self.switch_name
        with open(f"./mac-address-tables/{name}", "w") as out:
            out.write(str(self.sortedByInterface))
            out.close()

    def get_filtered_mac_table(self):
        # Calling all filtering and sorting functions
        self.filtering_mac_tables()  # Reads from ./mac-address-tables
        self.resolve_stack_and_interface()  #
        self.sort_by_interface()
        self.write_to_mac_directory()  # writes sortedByInterface Array
        return self.sortedByInterface
