"""
Inspired from: https://github.com/CaptTofu/switchssh/blob/master/switchssh/switchssh.py
"""
import paramiko


# built for/around HP Comware 7
class SwitchSSH:
    cmd_disable_paging = "screen-length disable\n"

    def __init__(self, host, username, password, port=22, timeout=30):
        self.host = host
        self.username = username
        self.password = password
        self.port = port
        self.timeout = timeout

        ssh = paramiko.SSHClient()
        ssh.set_missing_host_key_policy(paramiko.AutoAddPolicy())

        # Connecting to the switch
        try:
            print(f"Connecting to {host}...")
            ssh.connect(self.host,
                        username=self.username,
                        password=self.password,
                        timeout=self.timeout,
                        look_for_keys=False,
                        )
        except Exception as e:
            print("*** Caught exception: " + str(e.__class__) + ": " + str(e))

        # To Handle Connection After Establishing
        self.ssh = ssh

        # Opening a channel
        try:
            self.channel = ssh.invoke_shell()
        except Exception as e:
            print("*** Caught exception: " + str(e.__class__) + ": " + str(e))

        try:
            (stdin, stdout, stderr) = self.ssh.exec_command("dis log")
            print(stdout.read())
        except Exception as e:
            print("*** Caught exception: " + str(e.__class__) + ": " + str(e))
            self.ssh.close()
    # END __init__

    def exec_command(self, command, msg=""):
        stderr = ""
        stdout = ""
        stderr = ""
        try:
            # self.channel.send(command)
            stin, stdout, stderr = self.ssh.exec_command(command)
            while not stdout.channel.exit_status_ready():
                if stdout.channel.recv_ready():
                    rl, wl, xl = select.select([stdout.channel], [], [], 0.0)
                    if len(rl) > 0:
                        tmp = stdout.channel.recv(1024)
                        output = tmp.decode()
                        print(output)
        except Exception as e:
            print("*** Caught exception: " + str(e.__class__) + ": " + str(e))
