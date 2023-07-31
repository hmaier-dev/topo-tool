
-- Exportiere Datenbank-Struktur für topology-tool
CREATE DATABASE IF NOT EXISTS `topology-tool` /*!40100 DEFAULT CHARACTER SET utf8 */ /*!80016 DEFAULT ENCRYPTION='N' */;
USE `topology-tool`;

-- Exportiere Struktur von Tabelle topology-tool.clearpass
CREATE TABLE IF NOT EXISTS `clearpass` (
                                           `id` int NOT NULL AUTO_INCREMENT,
                                           `hostname` char(50) DEFAULT NULL,
                                           `mac` char(50) DEFAULT NULL,
                                           `ip` char(16) DEFAULT NULL,
                                           PRIMARY KEY (`id`)
) ENGINE=InnoDB  DEFAULT CHARSET=utf8 COMMENT='Organized API Data ';

-- Exportiere Struktur von Tabelle topology-tool.sw_ws
CREATE TABLE IF NOT EXISTS `sw_ws` (
                                       `id` int NOT NULL AUTO_INCREMENT,
                                       `interface_name` char(50) DEFAULT NULL,
                                       `mac` char(50) DEFAULT NULL,
                                       `hostname` char(50) DEFAULT NULL,
                                       `ip` char(16) DEFAULT NULL,
                                       `vlan` int DEFAULT NULL,
                                       `stack` int DEFAULT NULL,
                                       `interface_num` int DEFAULT NULL,
                                       PRIMARY KEY (`id`)
) ENGINE=InnoDB  DEFAULT CHARSET=utf8 COMMENT='192.168.132.130';

-- Exportiere Struktur von Tabelle topology-tool.sw_r2
CREATE TABLE IF NOT EXISTS `sw_r2` (
                                       `id` int NOT NULL AUTO_INCREMENT,
                                       `interface_name` char(50) DEFAULT NULL,
                                       `mac` char(50) DEFAULT NULL,
                                       `hostname` char(50) DEFAULT NULL,
                                       `ip` char(16) DEFAULT NULL,
                                       `vlan` int DEFAULT NULL,
                                       `stack` int DEFAULT NULL,
                                       `interface_num` int DEFAULT NULL,
                                       PRIMARY KEY (`id`)
) ENGINE=InnoDB  DEFAULT CHARSET=utf8 COMMENT='192.168.132.142';

-- Exportiere Struktur von Tabelle topology-tool.sw_r
CREATE TABLE IF NOT EXISTS `sw_r` (
                                      `id` int NOT NULL AUTO_INCREMENT,
                                      `interface_name` char(50) DEFAULT NULL,
                                      `mac` char(50) DEFAULT NULL,
                                      `hostname` char(50) DEFAULT NULL,
                                      `ip` char(16) DEFAULT NULL,
                                      `vlan` int DEFAULT NULL,
                                      `stack` int DEFAULT NULL,
                                      `interface_num` int DEFAULT NULL,
                                      PRIMARY KEY (`id`)
) ENGINE=InnoDB  DEFAULT CHARSET=utf8 COMMENT='192.168.132.141';

-- Exportiere Struktur von Tabelle topology-tool.sw_n1
CREATE TABLE IF NOT EXISTS `sw_n1` (
                                       `id` int NOT NULL AUTO_INCREMENT,
                                       `interface_name` char(50) DEFAULT NULL,
                                       `mac` char(50) DEFAULT NULL,
                                       `hostname` char(50) DEFAULT NULL,
                                       `ip` char(16) DEFAULT NULL,
                                       `vlan` int DEFAULT NULL,
                                       `stack` int DEFAULT NULL,
                                       `interface_num` int DEFAULT NULL,
                                       PRIMARY KEY (`id`)
) ENGINE=InnoDB  DEFAULT CHARSET=utf8 COMMENT='192.168.132.134';

-- Exportiere Struktur von Tabelle topology-tool.sw_kms
CREATE TABLE IF NOT EXISTS `sw_kms` (
                                        `id` int NOT NULL AUTO_INCREMENT,
                                        `interface_name` char(50) DEFAULT NULL,
                                        `mac` char(50) DEFAULT NULL,
                                        `hostname` char(50) DEFAULT NULL,
                                        `ip` char(16) DEFAULT NULL,
                                        `vlan` int DEFAULT NULL,
                                        `stack` int DEFAULT NULL,
                                        `interface_num` int DEFAULT NULL,
                                        PRIMARY KEY (`id`)
) ENGINE=InnoDB  DEFAULT CHARSET=utf8 COMMENT='192.168.132.128';

-- Exportiere Struktur von Tabelle topology-tool.sw_g
CREATE TABLE IF NOT EXISTS `sw_g` (
                                      `id` int NOT NULL AUTO_INCREMENT,
                                      `interface_name` char(50) DEFAULT NULL,
                                      `mac` char(50) DEFAULT NULL,
                                      `hostname` char(50) DEFAULT NULL,
                                      `ip` char(16) DEFAULT NULL,
                                      `vlan` int DEFAULT NULL,
                                      `stack` int DEFAULT NULL,
                                      `interface_num` int DEFAULT NULL,
                                      PRIMARY KEY (`id`)
) ENGINE=InnoDB  DEFAULT CHARSET=utf8 COMMENT='192.168.132.146';

-- Exportiere Struktur von Tabelle topology-tool.sw_f
CREATE TABLE IF NOT EXISTS `sw_f` (
                                      `id` int NOT NULL AUTO_INCREMENT,
                                      `interface_name` char(50) DEFAULT NULL,
                                      `mac` char(50) DEFAULT NULL,
                                      `hostname` char(50) DEFAULT NULL,
                                      `ip` char(16) DEFAULT NULL,
                                      `vlan` int DEFAULT NULL,
                                      `stack` int DEFAULT NULL,
                                      `interface_num` int DEFAULT NULL,
                                      PRIMARY KEY (`id`)
) ENGINE=InnoDB  DEFAULT CHARSET=utf8 COMMENT='192.168.132.202';


-- Exportiere Struktur von Tabelle topology-tool.sw_e2
CREATE TABLE IF NOT EXISTS `sw_e2` (
                                       `id` int NOT NULL AUTO_INCREMENT,
                                       `interface_name` char(50) DEFAULT NULL,
                                       `mac` char(50) DEFAULT NULL,
                                       `hostname` char(50) DEFAULT NULL,
                                       `ip` char(16) DEFAULT NULL,
                                       `vlan` int DEFAULT NULL,
                                       `stack` int DEFAULT NULL,
                                       `interface_num` int DEFAULT NULL,
                                       PRIMARY KEY (`id`)
) ENGINE=InnoDB  DEFAULT CHARSET=utf8 COMMENT='192.168.132.123';


-- Exportiere Struktur von Tabelle topology-tool.sw_e1
CREATE TABLE IF NOT EXISTS `sw_e1` (
                                       `id` int NOT NULL AUTO_INCREMENT,
                                       `interface_name` char(50) DEFAULT NULL,
                                       `mac` char(50) DEFAULT NULL,
                                       `hostname` char(50) DEFAULT NULL,
                                       `ip` char(16) DEFAULT NULL,
                                       `vlan` int DEFAULT NULL,
                                       `stack` int DEFAULT NULL,
                                       `interface_num` int DEFAULT NULL,
                                       PRIMARY KEY (`id`)
) ENGINE=InnoDB  DEFAULT CHARSET=utf8 COMMENT='192.168.132.138';

-- Exportiere Struktur von Tabelle topology-tool.sw_d
CREATE TABLE IF NOT EXISTS `sw_d` (
                                      `id` int NOT NULL AUTO_INCREMENT,
                                      `interface_name` char(50) DEFAULT NULL,
                                      `mac` char(50) DEFAULT NULL,
                                      `hostname` char(50) DEFAULT NULL,
                                      `ip` char(16) DEFAULT NULL,
                                      `vlan` int DEFAULT NULL,
                                      `stack` int DEFAULT NULL,
                                      `interface_num` int DEFAULT NULL,
                                      PRIMARY KEY (`id`)
) ENGINE=InnoDB  DEFAULT CHARSET=utf8 COMMENT='192.168.132.127';

-- Exportiere Struktur von Tabelle topology-tool.sw_c-sued
CREATE TABLE IF NOT EXISTS `sw_c-sued` (
                                           `id` int NOT NULL AUTO_INCREMENT,
                                           `interface_name` char(50) DEFAULT NULL,
                                           `mac` char(50) DEFAULT NULL,
                                           `hostname` char(50) DEFAULT NULL,
                                           `ip` char(16) DEFAULT NULL,
                                           `vlan` int DEFAULT NULL,
                                           `stack` int DEFAULT NULL,
                                           `interface_num` int DEFAULT NULL,
                                           PRIMARY KEY (`id`)
) ENGINE=InnoDB  DEFAULT CHARSET=utf8 COMMENT='192.168.132.122';

-- Exportiere Struktur von Tabelle topology-tool.sw_c-nord-core
CREATE TABLE IF NOT EXISTS `sw_c-nord-core` (
                                                `id` int NOT NULL AUTO_INCREMENT,
                                                `interface_name` char(50) DEFAULT NULL,
                                                `mac` char(50) DEFAULT NULL,
                                                `hostname` char(50) DEFAULT NULL,
                                                `ip` char(16) DEFAULT NULL,
                                                `vlan` int DEFAULT NULL,
                                                `stack` int DEFAULT NULL,
                                                `interface_num` int DEFAULT NULL,
                                                PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='192.168.132.120';

-- Exportiere Struktur von Tabelle topology-tool.sw_c-nord
CREATE TABLE IF NOT EXISTS `sw_c-nord` (
                                           `id` int NOT NULL AUTO_INCREMENT,
                                           `interface_name` char(50) DEFAULT NULL,
                                           `mac` char(50) DEFAULT NULL,
                                           `hostname` char(50) DEFAULT NULL,
                                           `ip` char(16) DEFAULT NULL,
                                           `vlan` int DEFAULT NULL,
                                           `stack` int DEFAULT NULL,
                                           `interface_num` int DEFAULT NULL,
                                           PRIMARY KEY (`id`)
) ENGINE=InnoDB  DEFAULT CHARSET=utf8 COMMENT='192.168.132.121';

-- Exportiere Struktur von Tabelle topology-tool.sw_b
CREATE TABLE IF NOT EXISTS `sw_b` (
                                      `id` int NOT NULL AUTO_INCREMENT,
                                      `interface_name` char(50) DEFAULT NULL,
                                      `mac` char(50) DEFAULT NULL,
                                      `hostname` char(50) DEFAULT NULL,
                                      `ip` char(16) DEFAULT NULL,
                                      `vlan` int DEFAULT NULL,
                                      `stack` int DEFAULT NULL,
                                      `interface_num` int DEFAULT NULL,
                                      PRIMARY KEY (`id`)
) ENGINE=InnoDB  DEFAULT CHARSET=utf8 COMMENT='192.168.132.132';

-- Exportiere Struktur von Tabelle topology-tool.sw_a-sued
CREATE TABLE IF NOT EXISTS `sw_a-sued` (
                                           `id` int NOT NULL AUTO_INCREMENT,
                                           `interface_name` char(50) DEFAULT NULL,
                                           `mac` char(50) DEFAULT NULL,
                                           `hostname` char(50) DEFAULT NULL,
                                           `ip` char(16) DEFAULT NULL,
                                           `vlan` int DEFAULT NULL,
                                           `stack` int DEFAULT NULL,
                                           `interface_num` int DEFAULT NULL,
                                           PRIMARY KEY (`id`)
) ENGINE=InnoDB  DEFAULT CHARSET=utf8 COMMENT='192.168.132.126';

-- Exportiere Struktur von Tabelle topology-tool.sw_c-nord-core
CREATE TABLE IF NOT EXISTS `sw_c-nord-core` (
                                                `id` int NOT NULL AUTO_INCREMENT,
                                                `interface_name` char(50) DEFAULT NULL,
                                                `mac` char(50) DEFAULT NULL,
                                                `hostname` char(50) DEFAULT NULL,
                                                `ip` char(16) DEFAULT NULL,
                                                `vlan` int DEFAULT NULL,
                                                `stack` int DEFAULT NULL,
                                                `interface_num` int DEFAULT NULL,
                                                PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='192.168.132.120';

-- Exportiere Struktur von Tabelle topology-tool.sw_a-nord
CREATE TABLE IF NOT EXISTS `sw_a-nord` (
                                           `id` int NOT NULL AUTO_INCREMENT,
                                           `interface_name` char(50) DEFAULT NULL,
                                           `mac` char(50) DEFAULT NULL,
                                           `hostname` char(50) DEFAULT NULL,
                                           `ip` char(16) DEFAULT NULL,
                                           `vlan` int DEFAULT NULL,
                                           `stack` int DEFAULT NULL,
                                           `interface_num` int DEFAULT NULL,
                                           PRIMARY KEY (`id`)
) ENGINE=InnoDB  DEFAULT CHARSET=utf8 COMMENT='192.168.132.125';

-- giving rights to www-data
# GRANT delete, insert, select, update, drop ON `topology-tool`.* TO 'www-data'@'localhost';