USE users;

CREATE TABLE IF NOT EXISTS user (
	Id int not null AUTO_INCREMENT,
	email varchar(60) not null,
    name VARCHAR(25) not null,
	pwdhash VARCHAR(100),
	age INT,
	aditional_information VARCHAR(255),
	PRIMARY KEY (Id)
);

