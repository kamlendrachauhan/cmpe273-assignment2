# cmpe273-assignment2

#To Test

POST http://localhost:8880/locations
{   "name" : "Kamlendra",   "address" : "190 Ryland St",   "city" : "San Jose",   "state" : "CA",   "zip" : "94110"}
Response
{"id":"562c6f6315ffb717c04dd47a","name":"Kamlendra","address":"190 Ryland St","city":"San Jose","state":"CA","zip":"94110","coordinate":{"lat":37.3408482,"lng":-121.8984085}}

GET http://localhost:8880/locations/562c6f6315ffb717c04dd47a
Response
{"id":"562c6f6315ffb717c04dd47a","name":"Kamlendra","address":"190 Ryland St","city":"San Jose","state":"CA","zip":"94110","coordinate":{"lat":37.3408482,"lng":-121.8984085}}
