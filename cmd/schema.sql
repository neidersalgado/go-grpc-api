CREATE TABLE user (
	Id VARCHAR(20),
    name VARCHAR(25) not null,
	pwdhash VARCHAR(100),
	age INT,
	aditional_information VARCHAR(255),
	PRIMARY KEY (Id)
);

CREATE TABLE parent (
	Id INT NOT NULL AUTO_INCREMENT,
	parent VARCHAR(20),
	son VARCHAR(20),
	PRIMARY KEY (Id),
    FOREIGN KEY (parent) REFERENCES user(Id),
    FOREIGN KEY (son) REFERENCES user(Id)
);