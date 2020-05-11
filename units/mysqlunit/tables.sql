SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
SET time_zone = "+00:00";

DROP DATABASE IF EXISTS test_database;
CREATE DATABASE test_database;

use test_database;

CREATE TABLE `test_table`
(
    `id` int(11),
    `name` varchar(500)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8;
