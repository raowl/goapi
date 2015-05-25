# goapi

A work in progress to make practical golang api with good practices that serves as a base for my own projects.
It needs mongodb running on the server.

for creating your own key files inside config directory:
$ openssl genrsa -out demo.rsa 1024 # the 1024 is the size of the key we are generating
$ openssl rsa -in demo.rsa -pubout > demo.rsa.pub 

* Api examples:

//register username
curl -H "Accept: application/vnd.api+json" -H 'Content-Type: application/vnd.api+json' -X POST -d \
'{"data":{"username":"example1","password":"example1"}}' http://localhost:8080/api/v1/user

// get token
curl -H "Accept: application/vnd.api+json" -H 'Content-Type: application/vnd.api+json' -X POST -d \
'{"data":{"username":"example1","password":"example1"}}' http://localhost:8080/api/v1/user/auth

//insert marker
curl -H \
"Authorization: Bearer \
REPLACEWITHTOKEN" \
-H "Accept: application/vnd.api+json" -H 'Content-Type: application/vnd.api+json' \
-X POST -d '{"data":{"coordinates":[10.23,2.2344],"name":"example", "address": "some street", "website":"www.google.com"}}' \
http://localhost:8080/markers

//get markers 

curl -H \
"Authorization: Bearer REPLACEWITHTOKEN" \
-H "Accept: application/vnd.api+json" -H 'Content-Type: application/vnd.api+json' \
http://localhost:8080/markers

* some resource that have been very helpful:

http://nicolasmerouze.com/build-web-framework-golang/
https://sendgrid.com/blog/tokens-tokens-intro-json-web-tokens-jwt-go/

I will be glad if you can join me in the development, thanks for reading.
