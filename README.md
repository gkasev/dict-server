###Redis:

+ Server should be configured to listen on the default port -- 6379.
+ Clients should be able to connect to redis without password (default setting)
Redis Data:

+ The sorted set, used for auto complete, should be called "myzset".
+ Use the following command to create a set for test purposes:
  
  ZADD myzset 1 "Anna" 1 "Brittany" 1 "Cinderella" 1 "Diana" 1 "Eva" 1 "Fiona" 1 "Gunda" 1 "Hege" 1 "Inga" 1 "Johanna" 1 "Kitty" 1 "Linda" 1 "Nina" 1 "Ophelia" 1 "Petunia" 1 "Amanda" 1 "Raquel" 1 "Cindy" 1 "Doris" 1 "Eve" 1 "Evita" 1 "Sunniva" 1 "Tove" 1 "Unni" 1 "Violet" 1 "Liza" 1 "Elizabeth" 1 "Ellen" 1 "Wenche" 1 "Vicky"

###Dictionary server:

+ Start the dictionary server from the command line. Specify the static path as argument:
  
  localhost:~ georgi$ dict-server /Users/georgi/go/bin/static

###Dictionary server behavior:

+ The server listens on "http://localhost:8080/"
+ For auto complete use "http://localhost:8080/autocomp/query", for example:

  localhost:~ georgi$ curl -i localhost:8080/autocomp/E
  HTTP/1.1 200 OK
  Content-Type: application/json
  Date: Fri, 01 May 2015 11:13:45 GMT
  Content-Length: 41

  ["Elizabeth","Ellen","Eva","Eve","Evita"]
