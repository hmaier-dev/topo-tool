# topo-tool
This tool pulls mac-addresses from a list of **HPE7-Comware-Switches** 
and compares them with those saved on a **Aruba-Clearpass** (**N**etwork**A**ccess**C**ontrol).

To get data from the switches, the `napalm-h3c-comware` python-library is used. 
For scraping the clearpass, the standard clearpass-api is used.

From this comparison, we can allocate which **switch-port** belongs to which **hostname**.

## Running without docker
Before running this tool you need several setup/in-mind:

- network-connection to comware-switches and to a clearpass

- username + password for the switches and the clearpass (create an api-reader)

- `data/cred.py` containing the credentials

```` python
db_user = ""
db_password = ""
api_user = ""
api_password = ""
sw_user = ""
sw_password = ""
````

- a mysql/mariadb-Database called `topology-tool` (schema will get generated from python)
- latest htmx-library in `web/js/htmx.min.js`

  - Linux: `wget https://unpkg.com/htmx.org@1.9.4/dist/htmx.min.js > web/js/htmx.min.js  `
  - Windows: `(wget https://unpkg.com/htmx.org@1.9.4/dist/htmx.min.js).Content | Out-File -Path web/js/htmx.min.js  `



## Usage
By now there a two use-case implemented.
### CLI
You can run the `data/main.py` with the following options.
- `--scanner`: connect to the Clearpass-API, to pull information about all discovered endpoints, and to the network-switches, to get all mac-addresses in the network
- `--search <hostname>`: search through the database to get all collected information about a device

### Web-Interface
You can run the `web/server.go` to get a web-interface for searching through the database. Remember that you first have to run `data/main.py --scanner` to populate the database
