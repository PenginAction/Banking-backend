# Banking


# mysql

 CREATE TABLE User (
    UserID INT AUTO_INCREMENT PRIMARY KEY,
    Name VARCHAR(255) NOT NULL,
    Email VARCHAR(255) UNIQUE NOT NULL,
    Password VARCHAR(255) NOT NULL,
    Pin VARCHAR(255) NOT NULL
);

CREATE TABLE Account (
    AccountID VARCHAR(10) PRIMARY KEY,
    Balance DECIMAL(15, 2) NOT NULL,
    UserID INT,
    FOREIGN KEY (UserID) REFERENCES User(UserID)
);

CREATE TABLE Transaction (
    TransactionID INT AUTO_INCREMENT PRIMARY KEY,
    TransactionAmount DECIMAL(15, 2) NOT NULL,
    TransactionType ENUM('出金', '入金', "振込") NOT NULL,
    TransactionDate TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    SourceAccountID VARCHAR(10),
    DestinationAccountID VARCHAR(10),
    FOREIGN KEY (SourceAccountID) REFERENCES Account(AccountID),
    FOREIGN KEY (DestinationAccountID) REFERENCES Account(AccountID)
);
