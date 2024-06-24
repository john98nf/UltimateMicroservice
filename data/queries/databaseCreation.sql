create table
  `Companies` (
    `ID` int unsigned not null comment 'Primary Key',
    `NAME` VARCHAR(255) not null comment 'Name of the company (Primary Key)',
    `DESCRIPTION` TEXT null comment 'Description of the company (optional)',
    `EMPLOYEES` INTEGER UNSIGNED not null comment 'Number of Employees',
    `REGISTRATION_STATUS` BOOLEAN not null comment 'A flag indigating if the company is registered or not.',
    `TYPE` ENUM('Corporations','NonProfit','Cooperative','Sole Proprietorship') not null comment 'Type of the company',
    primary key (`ID`, `NAME`)
  );

alter table
  `Companies`
modify column
  `ID` int unsigned not null auto_increment

create table
  `USERS` (
    `ID` int unsigned not null comment 'Primary Key',
    `USERNAME` VARCHAR(255) not null comment 'Username (Primary Key)',
    `DISPLAY_NAME` TINYTEXT null,
    primary key (`ID`, `USERNAME`)
  );

alter table
  `USERS`
modify column
  `ID` int unsigned not null auto_increment