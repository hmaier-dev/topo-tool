import json
import os
import re

'''
This class should filter all unnecessary entrys as well, as sort the useful ones
'''


class Filtering:
    def __init__(self):
        """
        Logic is processed in the order of the following lists
        """
        self.noBridgeAggregation = []
        self.withStackAndInterface = []
        self.sortedByInterface = []

    def cleaning_mac_table(self, switch, table, regex_filter="Bridge-Aggregation"):
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
