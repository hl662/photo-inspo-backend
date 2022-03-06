# Photo-Inspo Backend

## Description

This Go module serves as a HTTP server for my [Photo-Inspo](https://github.com/hl662/photo-inspo) web application.

## Prerequisites
1. Go 1.17 ([Install Here](https://go.dev/doc/install))
2. A [MongoDB](https://www.mongodb.com) Account
3. A MongoDB [Atlas](https://docs.atlas.mongodb.com/getting-started/) cluster deployed, and a database created within the cluster. 
A database username and password set up.

## Instructions

1. Open a bash terminal
2. In the project directory, run `go build`. After running, make sure there is an executable file like `photo-inspo-backend`
3. We need to set 4 environment variables in the same terminal: <br>
`export mongoPwd=<YOUR_ATLAS_DATABASE_PASSWORD>`<br>
`export mongoUsr=<YOUR_ATLAS_DATABASE_USERNAME>`<br>
`export mongoDBName=<YOUR_ATLAS_DATABASE_NAME>`<br>
`export encryptKey=<A_32_BIT_LONG_KEY>`<br>
4. Run the executable file. In your terminal, you can run `./photo-inspo-backend`
5. You can run curl requests, one example would be:
`curl --location --request POST 'https://localhost:8080/signup' \
   --header 'Content-Type: application/json' \
   --data-raw '{
   "username": "test",
   "password": "foobar"
   }'`
<br> If you go back to your mongoDB Atlas database, you can see a new document created.