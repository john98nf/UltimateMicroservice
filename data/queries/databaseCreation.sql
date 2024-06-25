create table
  `COMPANIES` (
    `ID` binary(16) not null comment 'Primary Key',
    `NAME` VARCHAR(15) not null unique comment 'Name of the company (unique)',
    `DESCRIPTION` TEXT null comment 'Description of the company (optional)',
    `EMPLOYEES` INTEGER UNSIGNED not null comment 'Number of Employees',
    `REGISTRATION_STATUS` BOOLEAN not null comment 'A flag indigating if the company is registered or not.',
    `LEGAL_TYPE` ENUM('Corporations','NonProfit','Cooperative','Sole Proprietorship') not null comment 'Type of the company',
    primary key (`ID`)
  );

create table
  `USERS` (
    `ID` binary(16) not null comment 'Primary Key',
    `USERNAME` VARCHAR(255) unique not null comment 'Username (Unique)',
    `DISPLAY_NAME` TINYTEXT null,
    primary key (`ID`)
  );