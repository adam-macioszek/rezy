# Rezy the Reservation system
This is the backend for a restaurant reservation system. it supports two endpoints, /reservation and /book
with which a user can look at reservations and book them.
## Requirements
Docker

## Usage
In order to setup up the API server you will need to have Docker installed on your.
Running 

```
docker compose up
```
From the top-level package directory, this will create 2 docker containers, one hosting the API server, and the other hosting the database.

By default the API server runs on 8080, this can be modified in the app.env file. To change what port docker exposes for the api use the docker-compose file.

The program will currently output a jwt token in the terminal that is used in the auth bearer field
of request to create new reservations.

### Example API usage
Currently, the /reservation supports Get and Post.
Get a request with send back a list of non-booked reservations in the database, with a
default page size of 5. This page size can be modified up to 15. Individual reservations
can be retrieved by supplying the id.
## /reservation
- GET: returns a list of unbooked reservations, the default page size of 5, max of 15, and default page offset of 1.
- POST: upload a new reservation, must have a valid JWT token in auth Bearer Header

## /reservation/:id  
- GET: returns a reservation if one exists matching the ID
```
GET
/reservation/10                              //retrieves reservation with id of 10 if present
/reservation?page_size=10                   //returns 10 reservations in database
/reservation?page_size=5&page_id=2          //returns 5 reservations with a page offset of 1 
```

## /book
- POST: used to book a reservation, reservation must start after the current time, be unbooked, and have table size larger than party size.
Our extremely optimized backend system will automatically find a reservation at the same time with a smaller table size that matches your desired size if it exists, so the reservation you return might be different from what you request.
```
POST
/book
{
    "id": 62,
    "table_size": 4
}
```

## Testing
There is some basic testing done for the CRUD database operations and for the MakeReservation transaction. In a production 
environment, this would need to be increased to look for things like possible deadlocks. The testing as it exists now is simply meant to validate the basic use cases and as a sanity check. 
There is also basic validation done with mocking for the API, but the coverage here is not complete and only really covers
the Get reservation endpoint, and not extensively. In a production environment, this would need to be improved.

To run the test for the entire project use the following command from the top-level directory:
```
go test -v ./...
```
To run individual tests navigate to the directory with the test and run
```
go test file-to-test.go
```

## Known Limitations/ Future improvements
- The database does not use ssl and the credentials are set to root:secret which should be changed.
- The database might be returning to much information to the api with errors, exposing the implementation details.
- More comprehensive Database transaction test to check for deadlocks.
- A users database and authorization should be added so each user can login, get their own tokens and only modify reservations they created.
- The JWT is currently create with HS256, a asymetric key should probably be used.
- More extensive testing on the API needs to be done to ensure it is functionally correct.
- The API should use HTTPS not HTTP
