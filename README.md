# topo-tool
This tool pulls mac-addresses from a list of **HPE7-Comware-Switches** 
and compares them with those saved on a **Aruba-Clearpass** (**N**etwork**A**ccess**C**ontrol).

From this comparison, the program can allocate which **switch-port** belongs to which **hostname**.

I wrote this tool, because I was instructed to inventory the location all computers at work. Even though I could not figure out the exact room, with this tool I could get the house in which the device is standing. 
This was a huge time save for me. 

## Running without docker
Before running this tool you need several setup/in-mind:

- You **need** a network-connection to **comware-switches** and to a **clearpass-api**. 

- username + password for the switches and the clearpass (create an api-reader)

- `data/cred.txt` containing the credentials, which looks like the following.

```` txt
## Clearpass API-User
api_user=<user>
api_password=<password>

## User for Comware-Switches
sw_user=<user>
sw_password=<password>
````

- A mysql/mariadb-database called `topology-tool` setup from `db/setup.sql`. 
- latest htmx-library in `web/js/htmx.min.js`

  - Linux: `wget https://unpkg.com/htmx.org@1.9.4/dist/htmx.min.js > web/js/htmx.min.js  `
  - Windows: `(wget https://unpkg.com/htmx.org@1.9.4/dist/htmx.min.js).Content | Out-File -Path web/js/htmx.min.js  `


### Web-Interface
You can run/build the `web/server.go` to get a web-interface for searching through the database. 
The address is `localhost:8181/topotool`.
