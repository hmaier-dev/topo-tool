import cred
import mariadb


class Database:
    def __init__(self, host, port):
        self.USERNAME = cred.db_user
        self.PASSWORD = cred.db_password
        self.DB_NAME = "topology-tool"
        self.host = host
        self.port = port
        self.cursor = self.connect()

    def connect(self):
        try:
            conn = mariadb.connect(
                user=self.USERNAME,
                password=self.PASSWORD,
                host=self.host,
                port=self.port,
                database=self.DB_NAME
            )
            return conn.cursor()
        except mariadb.Error as e:
            # if authentication-plugin-caching-sha2-password-cannot-be-loaded
            # ALTER USER 'root'@'localhost' IDENTIFIED WITH mysql_native_password BY 'root';
            print(f"Error connecting to MariaDB Platform: {e}")

    # Setting up the database
    def setup(self, switch_array):
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

    # Remove table
    def drop(self, switch_array):
        def cmd(name): return f" DROP TABLE `{name}`;"
        for entry in switch_array:
            sw_name = entry[0]
            sql = cmd(sw_name)
            self.cursor.execute(sql)

    # Remove data from table
    def truncate(self, switch_array):
        def cmd(name): return f" TRUNCATE TABLE `{name}`;"
        for entry in switch_array:
            sw_name = entry[0]
            sql = cmd(sw_name)
            self.cursor.execute(sql)

    # Data containing Mac + Switch-Port combo
    def insert_switch_data(self, json):

        pass

    # Data containing Hostname + Mac combo
    def insert_api_data(self):
        pass
