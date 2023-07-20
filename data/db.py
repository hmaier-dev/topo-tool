import cred
import mariadb


class Database:
    def __init__(self, host, port):
        self.USERNAME = cred.db_user
        self.PASSWORD = cred.db_password
        self.DB_NAME = "topology-tool"

        self.host = host
        self.port = port
        self.conn = self.connect()
        self.cursor = self.conn.cursor()
        self.conn.autocommit = True

    def connect(self):
        try:
            conn = mariadb.connect(
                user=self.USERNAME,
                password=self.PASSWORD,
                host=self.host,
                port=self.port,
                database=self.DB_NAME
            )
        except mariadb.Error as e:
            # if authentication-plugin-caching-sha2-password-cannot-be-loaded
            # ALTER USER 'root'@'localhost' IDENTIFIED WITH mysql_native_password BY 'root';
            print(f"Error connecting to MariaDB Platform: {e}")
            return

        return conn

    # Setting up the database
    def setup_switch_tables(self, switch_array):
        def cmd(name, ip): return (
            f"CREATE TABLE IF NOT EXISTS `{name}` ("
            f"`id` int NOT NULL AUTO_INCREMENT,"
            f"`interface_name` char(50) DEFAULT NULL,"
            f"`mac` char(50) DEFAULT NULL, "
            f"`hostname` char(50) DEFAULT NULL,"
            f"`vlan` int DEFAULT NULL,"
            f"`stack` int DEFAULT NULL,"
            f"`interface_num` int DEFAULT NULL,"
            f" PRIMARY KEY (`id`)"
            f") ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='{ip}';"
        )

        for entry in switch_array:
            sw_name = entry[0]
            sw_ip = entry[1]
            sql = cmd(sw_name, sw_ip)
            self.cursor.execute(sql)

    def setup_clearpass_table(self):
        sql = (
            f"CREATE TABLE IF NOT EXISTS `clearpass` ("
            f"`id` int NOT NULL AUTO_INCREMENT,"
            f"`hostname` char(50) DEFAULT NULL,"
            f"`mac` char(50) DEFAULT NULL, "
            f" PRIMARY KEY (`id`)"
            f") ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='Organized API Data ';"
        )
        self.cursor.execute(sql)

    # Remove table
    def drop(self, switch_array):
        def cmd(name): return f" DROP TABLE `{name}`;"

        for entry in switch_array:
            sw_name = entry[0]
            sql = cmd(sw_name)
            self.cursor.execute(sql)

    # Remove data from table
    def truncate(self, name):
        def cmd(name): return f" TRUNCATE TABLE `{name}`;"
        cmd2 = lambda name: f"ALTER TABLE `{name}` AUTO_INCREMENT = 1 ;"
        sql = cmd(name)
        sql2 = cmd2(name)
        self.cursor.execute(sql)
        self.cursor.execute(sql2)

    # Data containing Mac + Switch-Port combo
    def insert_switch_data(self, switch, mac_list):
        name = switch[0]
        cmd = lambda name, interface_name, mac, vlan, stack, interface_num: (
            f"INSERT INTO `{name}` "
            f" (interface_name, mac, vlan, stack, interface_num)"
            f" VALUES "
            f" ('{interface_name}','{mac}', {vlan}, {stack}, {interface_num});"

        )
        for entry in mac_list:
            int_name = entry["interface_name"]
            m = entry["mac"]
            v = entry["vlan"]
            s = entry["stack"]
            int_num = entry["interface_num"]
            sql = cmd(name, int_name, m, v, s, int_num)
            print(sql)
            try:
                resp = self.cursor.execute(sql)
            except mariadb.Error as e:
                print(f"Error: {e}")

    # Data containing Hostname + Mac combo
    def insert_api_data(self, json):
        cmd = lambda hostname, mac: (
            f"INSERT INTO `clearpass` "
            f" (hostname, mac)"
            f" VALUES "
            f"('{hostname}', '{mac}') ;"

        )
        for entry in json:
            h = entry["hostname"]
            m = entry["mac"]
            sql = cmd(h, m)
            print(sql)
            try:
                resp = self.cursor.execute(sql)
            except mariadb.Error as e:
                print(f"Error: {e}")

    def select_clearpass_data(self):
        sql = f"SELECT * FROM clearpass ;"
        self.cursor.execute(sql)
        rows = self.cursor.fetchall()
        return rows

    def select_switch_data(self, switch):
        name = switch[0]
        sql = f"SELECT * FROM `{name}` ;"
        self.cursor.execute(sql)
        rows = self.cursor.fetchall()
        return rows

    def select_hostname_by_mac(self, mac):
        sql = f"SELECT hostname FROM `clearpass` WHERE mac = ? "
        self.cursor.execute(sql, (mac,))
        row = self.cursor.fetchone()
        return row

    def update_hostname_by_id(self, hostname, id, switch):
        name = switch[0]
        sql = f"UPDATE `{name}` SET hostname = '{hostname}' WHERE id = {id} ;"
        self.cursor.execute(sql)

    def show_switch_tables(self):
        add = []
        sql = f"SHOW TABLES;"
        self.cursor.execute(sql)
        row = self.cursor.fetchone()
        for entry in row:
            if entry != "clearpass":
                add.append(entry)
        return add

