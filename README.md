# topo-tool
This tool pulls all MAC-Tables from **HPE7-Comware-Switches** and compares 
them with those on the **Clearpass (Network Access Control)**. 
From this comparision, we can allocate which **switch-port** belongs to which **hostname**.
## Pre-Requisites
Before running this tool you need several setup/in-mind:

- (obviously) network-connection to comware-switches aswell to a clearpass

- a mariadb

- `cred.py`-file to import username + credentials for 
    - api-user of the clearpass
    - user for ssh-access to the comware-switches
    - database-user
