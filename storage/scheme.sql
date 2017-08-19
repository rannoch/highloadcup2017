drop database if exists hlcup2017;

create database hlcup2017 CHARACTER SET utf8 COLLATE utf8_general_ci;

use hlcup2017;

create table user (
	id int primary key,
  email varchar(100),
  first_name varchar(50),
  last_name varchar(50),
  gender varchar(1),
  birth_date int
) ENGINE=InnoDB DEFAULT CHARSET=utf8;


create table location (
	id int primary key,
  place text,
  country varchar(50),
  city varchar(50),
  distance int
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

create table visit (
	id int primary key,
  location int,
  user int,
  visited_at int,
  mark tinyint(1)
)
ENGINE=InnoDB DEFAULT CHARSET=utf8;