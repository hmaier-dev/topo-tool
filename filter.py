import json
import os

'''
This class should filter all unnecessary entrys as well as sort the useful ones
'''

class Filtering:
    def __init__(self, file_name):
        self.table_name = os.path.join("./mac-address-tables", file_name)

        # Opening the desired file
        with open(self.table_name, "r") as file:
            self.table = json.loads(file.read())
            file.close()




