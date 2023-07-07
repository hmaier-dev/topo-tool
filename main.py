import subprocess
from topology import Topology
from filter import Filtering

SWITCHES = [
    ("SW_A-Nord", "192.168.132.125"),
    ("SW_A-Sued", "192.168.132.126"),
    ("SW_C-Nord-Core", "192.168.132.120"),
    ("SW_C-Nord", "192.168.132.121"),
    ("SW_C-Sued", "192.168.132.122"),
    ("SW_B", "192.168.132.132"),
    ("SW_R", "192.168.132.141"),
    ("SW_R2", "192.168.132.142"),
    ("SW_D", "192.168.132.127"),
    ("SW_E1", "192.168.132.138"),
    ("SW_E2", "192.168.132.123"),
    ("SW_F", "192.168.132.202"),
    ("SW_G", "192.168.132.146"),
    ("SW_N1", "192.168.132.134"),
    # ("SW_N2", "192.168.132.144"), currently not reachable
    ("SW_WS", "192.168.132.130"),
    ("SW_KMS", "192.168.132.128"),
    ("SW_P", "192.168.132.136"),
    # ("SW_A121", "192.168.132.131"), 1/0/24 set as uplink, no BridgeAggregation filtering possible
]


def connection_test(switch_array):
    """
    Ping all Switches located in the array
    """
    for switch in switch_array:
        name = switch[0]
        ip = switch[1]
        response = subprocess.run("ping -n 1 " + ip, stdout=subprocess.DEVNULL)
        if response == 0:
            print(f"{name} is up!")
        else:
            print(f"{name} is down!")


def scan_for_hosts(ip_range):
    """
    Scan the given IP address range using Nmap
    """
    nmap_args = ['nmap', '-n', '-sP', ip_range]
    return subprocess.run(nmap_args)


if __name__ == "__main__":
    # username = ""
    # password = ""

    # topo = Topology(username, password)
    # for switch in SWITCHES:
    #     topo.pull_all_mac_tables(switch)  # Writes to ./mac-address-tables Directory

    for switch in SWITCHES:
        f = Filtering(switch)
        f.get_filtered_mac_table()  # Reads from ./mac-address-tables
