CREATE TABLE user (
	Id int not null AUTO_INCREMENT,
	email varchar(60) not null,
    name VARCHAR(25) not null,
	pwdhash VARCHAR(100),
	age INT,
	aditional_information VARCHAR(255),
	PRIMARY KEY (Id)
);

CREATE TABLE parent (
	Id INT NOT NULL AUTO_INCREMENT,
	parent int,
	son int,
	PRIMARY KEY (Id),
    FOREIGN KEY (parent) REFERENCES user(Id),
    FOREIGN KEY (son) REFERENCES user(Id)
);