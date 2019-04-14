BEGIN TRANSACTION;
CREATE TABLE IF NOT EXISTS `scandata` (
	`uuid`	TEXT NOT NULL UNIQUE,
	`scanargs`	TEXT,
	`scanstart`	TEXT,
	`scantype`	TEXT,
	`scanprotocol`	TEXT,
	`scanservices`	TEXT,
	`scanend`	TEXT,
	`summary`	TEXT
);
CREATE TABLE IF NOT EXISTS `portdata` (
	`uuid`	TEXT,
	`ip`	TEXT,
	`port`	TEXT,
	`state`	TEXT,
	`reason`	TEXT,
	`service`	TEXT,
	`method`	TEXT
);
CREATE TABLE IF NOT EXISTS `hostdata` (
	`uuid`	TEXT NOT NULL,
	`ip`	TEXT NOT NULL,
	`host_state`	TEXT,
	`reason`	TEXT,
	`hostname`	TEXT
);
COMMIT;
