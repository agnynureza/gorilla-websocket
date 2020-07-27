# gorilla-websocket

testing 1 
1. open 2 terminal 
2. first terminal, go run main.go 
3. second terminal, live-server

simple websocket receive data from client and show in server console

testing 2
1. go run main.go 
2. open index.html in browser (appear header template "Saham/Price")
3. hit localhost:8080/data in postman (POST body raw {name:"BBCA", price:"30.000"})

see how we successfully push data to the client without needing to reaload the page 