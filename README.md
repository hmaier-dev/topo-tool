# topo-tool
This tool pulls mac-addresses from a list of **HPE7-Comware-Switches** 
and compares them with those saved on a **Aruba-Clearpass** (**N**etwork**A**ccess**C**ontrol).

From this comparison, we can allocate which **switch-port** belongs to which **hostname**.

## Running without docker
Before running this tool you need several setup/in-mind:

- network-connection to comware-switches and to clearpass

- username + password for the switches and the clearpass (create an api-reader)

- `data/cred.txt` containing the credentials

```` txt
api_user=<user>
api_password=<password>

sw_user=<user>
sw_password=<password>
````

- a mysql/mariadb-Database called `topology-tool` 
- latest htmx-library in `web/js/htmx.min.js`

  - Linux: `wget https://unpkg.com/htmx.org@1.9.4/dist/htmx.min.js > web/js/htmx.min.js  `
  - Windows: `(wget https://unpkg.com/htmx.org@1.9.4/dist/htmx.min.js).Content | Out-File -Path web/js/htmx.min.js  `

### Data-Scrapper
Make sure the 

### Web-Interface
You can run the `web/server.go` to get a web-interface for searching through the database. 
The address is `localhost:8181/topotool`.
