import subprocess
from topology import Topology

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
    topo = Topology()
    topo.pull_mac_tables()
