<h1>ChatRoom (Alan Sikarov, Keith Kwan, Dawin Ye) </h1>

<h2> Setup </h2>
Our code is split up into two parts, a server file and a client file.
Servers must be initialized first for a client to connect, so create 
a server by running <code><b>go run server.go {port number}</b></code>. To 
create a client, run <code><b>go run client.go {host address} {port number} {username}</b></code>.
For the purposes of this assignment, we have decided to limit the host address to 
127.0.0.1 so clients can only connect to servers on the same machine. 

<h2> Usage </h2>
After running the server in a terminal, the only command that the server can take is
"EXIT", which will exit the server and close all clients and their connections. Other 
inputs into the server terminal will not be supported. 

To send a message from user1 to user2 go to user1's client and enter 

"<b>user2,, {message to deliver}</b>"

We decided to use two commas as a separator between the username and content because two commas 
wouldn't normally be part of a message, whereas one comma or any other punctuation would be. Messages 
that are not formatted this way will generate errors, such as not having a separator or having multiple 
sets of separators.

<h2> Documentation </h2>


<h2> Justification</h2>

