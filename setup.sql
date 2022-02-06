CREATE database activities_groups;

USE activities_groups;

CREATE TABLE Conversations (ConversationID varchar (10) NOT NULL PRIMARY KEY, PosterName varchar(30) NOT NULL, PostContent varchar(50) NOT NULL, PostDate date NOT NULL, PostTime time NOT NULL); 

CREATE TABLE Replies (ReplyID int NOT NULL PRIMARY KEY, PosterName varchar(30) NOT NULL, PostContent varchar(50) NOT NULL, PostDate date NOT NULL, PostTime time NOT NULL, ConversationID varchar (5) NOT NULL); 

