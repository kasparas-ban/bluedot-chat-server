User A attempts to send a message to user B.
1. User A opens a WebSocket connection to the chat server via http://localhost:3002
2. On successfull connection, the first message user A sends is its username
    { username: userA }
3. The server saves the connection to the connectionsDB
    chat_server_id  username
    localhost:3002  userA