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

        self.raw_file_name = "RAW_" + switch[0] + ".json"
        self.clean_file_name = "CLEAN_" + switch[0] + ".json"
        self.raw_file_ip = switch[1]

        table_name = os.path.join("switch-tables", self.raw_file_name)
        # Opening the desired file
        with open(table_name, "r") as file:
            self.raw_mac_table = json.loads(file.read())
            file.close()

    def cleaning_mac_table(self, regex_filter="Bridge-Aggregation"):
        table = self.raw_mac_table
        tmp_list = []
        for entry in table:
            mac = entry["mac"]
            interface_name = entry["interface"]
            vlan = entry["vlan"]
            # Step 1
            # Drop all non-local mac-addresses
            if not re.match(regex_filter, interface_name):
                # Step 2
                # Resolving num of stack and interface
                matches = re.findall(r'\d+', interface_name)
                sort_helper = matches[0] + matches[2]  # combine two strings
                add = {"mac": mac,
                       "interface": interface_name,
                       "vlan": vlan,
                       "stack_num": matches[0],
                       "interface_num": matches[2],
                       "sort_helper": int(sort_helper)}
                # Step 3
                # Adding the new entry to a temporary list
                tmp_list.append(add)
            # Step 4
            # returning the list sorted by the sort_helper key
            return sorted(tmp_list, key=lambda x: (x["sort_helper"]))

    def write_to_mac_directory(self, mac_table):
        path = os.path.join("switch-tables", self.clean_file_name)
        with open(path, "w") as out:
            out.write(str(mac_table))
            out.close()

    def get_filtered_mac_table(self):
        # Calling all filtering and sorting functions
        clean = self.cleaning_mac_table()
        self.write_to_mac_directory(clean)  # Writes clean table
        return clean
