//This file is used to create a default user when running game server on local
//Reference - https://faun.pub/managing-mongodb-on-docker-with-docker-compose-26bf8a0bbae3
db.createUser(
    {
        user: "developer",
        pwd: "developer",
        roles: [
            {
                role: "readWrite",
                db: "demo"
            }
        ]
    }
)