# Transactional Key Value Store
This project is an interactive command line interface to a transaction key value store. 
Users are able to compile and run the program and get a shell with prompt where they can type commands. Users can enter
commands to set/get/delete key/value pairs and count values. All values are treated as strings, and the key/value data are
 stored in-memory for a session, it does not write to disk.

The cli interface allows users to perform operations in transactions, which gives users the flexibility to commit or 
rollback their changes to the key-value store. It includes the ability to nest transactions, rollback, and commit within nested transactions. 
The program is written purely in Go, it does not depend on any third party library.

The CLI supports the following commands: <br>
`SET <key> <value> // store the value for key` <br>
`GET <key>         // return the current value for key` <br />
`DELETE <key>      // remove the entry for key` <br />
`COUNT <value>     // return the number of keys that have the given value` <br />
`BEGIN             // start a new transaction` <br />
`COMMIT            // complete the current transaction` <br />
`ROLLBACK          // revert to state prior to BEGIN call` <br />


### Requirements
- Docker


### How to run
- Clone or download the project
- Navigate to the root directory on your terminal
- Run `make docker-build` to build the docker image
- Run `make docker-run` to run the generated image in a container.
