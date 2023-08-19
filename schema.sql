CREATE TABLE house(
    id int NOT NULL AUTO_INCREMENT,
    name varchar(255) NOT NULL,
    address varchar(255),
    country varchar(255),
    description varchar(500),
    price float(10,7),
    photo varchar(255),
    PRIMARY KEY(id) 
);
