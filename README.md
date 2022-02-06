# ETI Assignment 2 - 3.23 Activities within Groups
# Design consideration of your microservices
The activities have been split into Conversations Service and Replies Service. There are two APIs, one database containing two tables, and one console. 

In the database, there are two tables named Conversations and Replies.

<table>
  <tr><th colspan="3">Conversations Table</th></tr>
  <tr>
    <th>Name</th>
    <th>Date Type</th>
    <th>Description</th>
  </tr>
  <tr>
    <td>ConversationID</td>
    <td>varchar(10)</td>
    <td>Conversation ID specified by the user, PRIMARY KEY</td>
  </tr>
  <tr>
    <td>PosterName</td>
    <td>varchar(30)</td>
    <td>Name of the user posting</td>
  </tr>
  <tr>
    <td>PostContent</td>
    <td>varchar(50)</td>
    <td>Post Content that the user writes in the conversation</td>
  </tr>
    <tr>
    <td>PostDate</td>
    <td>date</td>
    <td>Date of created conversation</td>
  </tr>
    <tr>
    <td>PostTime</td>
    <td>time</td>
    <td>Time of created conversation</td>
  </tr>
</table>

<table>
  <tr><th colspan="3">Replies Table</th></tr>
  <tr>
    <th>Name</th>
    <th>Date Type</th>
    <th>Description</th>
  </tr>
  <tr>
    <td>ReplyID</td>
    <td>int</td>
    <td>Reply ID automatically generated, PRIMARY KEY</td>
  </tr>
  <tr>
    <td>PosterName</td>
    <td>varchar(30)</td>
    <td>Name of the user posting</td>
  </tr>
  <tr>
    <td>PostContent</td>
    <td>varchar(50)</td>
    <td>Post Content that the user writes in the reply</td>
  </tr>
    <tr>
    <td>PostDate</td>
    <td>date</td>
    <td>Date of created reply</td>
  </tr>
    <tr>
    <td>PostTime</td>
    <td>time</td>
    <td>Time of created reply</td>
  </tr>
    <tr>
    <td>ConversationID</td>
    <td>varchar(10)</td>
    <td>Conversation ID that the user is replying to</td>
  </tr>
</table>

<h2>Conversations Service</h2>
<h3>View</h3>
The conversations are retrieved and displayed for the user to view. The user can view all conversations available, or view a specific conversation by typing in the Conversation ID. The view conversations function displays all information regarding the conversation including the Conversation ID, Poster Name, Post Content, Post Date, and Post Time.

<h3>Create</h3>
A conversation is created by the user after inputting all necessary information. The user can create a conversation by inputting information like Conversation ID which can be named anything as long as it is within the character limit, Poster Name, and Post Content.

<h3>Update</h3>
An existing conversation created by the user is updated. The user can update a conversation by inputting in a new Conversation ID, Poster Name, and/or Post Content. Any new information entered in by the user will replace the existing information. 

<h3>Delete</h3>
An existing conversation created by the user is deleted. The user can delete a conversation by inputting the Conversation ID.

<h2>Replies Service</h2>
<h3>View</h3>
The replies are retrieved and displayed for the user to view. The user can view all replies available, or view a specific reply by typing in the Reply ID. The view replies function displays all information regarding the reply including the Reply ID, Poster Name, Post Content, Post Date, Post Time, and Conversation ID.

<h3>Create</h3>
A reply is created by the user after inputting all necessary information. The user can create a reply by inputting information like Conversation ID, Poster Name, and Post Content. The Reply ID is automatically generated when the reply is created. 

<h3>Update</h3>
An existing reply created by the user is updated. The user can update a reply by inputting in a Reply ID, and a new Poster Name, and/or Post Content. Any new information entered in by the user (aside from the Reply ID) will replace the existing information. 

<h3>Delete</h3>
An existing reply created by the user is deleted. The user can delete a reply by inputting the Reply ID.

<h2>Console</h2>
The console is the frontend which allows the user to input and call the services and functions implemented inside.

# Architecture diagram
![Untitled Diagram (1)](https://user-images.githubusercontent.com/64128624/152691608-3b850362-1d6b-4fce-9519-a906f5b9a02b.jpg)

# Link to your container image
https://hub.docker.com/repository/docker/s10175863/assignment2 <br/>
https://hub.docker.com/repository/docker/s10175863/assignment2_replies <br/>
https://hub.docker.com/repository/docker/s10175863/assignment2_console

# Instructions for setting up and running your microservices
1. Clone the Repository
2. Create the database

docker pull s10175863/assignment2:latest<br>
docker pull s10175863/assignment2_replies:latest<br/>
docker pull s10175863/assignment2_console:latest<br/>
