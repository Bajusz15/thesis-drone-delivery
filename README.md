#Máté Bajusz's bsc thesis, drone delivery
My BSc thesis. I wrote my thesis on the data management problems of parcel delivery with drones. This included optimizing routes for delivery with a combinatorial optimization algorithm and comparing the effectiveness of communication protocols and databases. The system design required that the protocols and databases could be exchanged for comparison, which was achieved with ports and adapters architecture. The system is made up of 2 containerized applications that can scale independently, both coded in Go. Used technologies such as Docker, gRPC, MongoDB, PostgreSQL.

This file contains instructions on setting up the drone-deliver simulation.
## Usage

**1) Check if your directory structure looks like the following:**
+ thesis-drone-delivery
    - backend
        - databases
        - server
        - drone-swarm
    - benchmark
    - web-client
    - docker-compose.yml
    -   README.md

**2) Go into the root of the project and start the project.**
+ On the very first time you start this

  Make sure you have PostgreSQL installed, and running on your machine.
  Make sure you have MongoDB installed, and running on your machine.
```bash
$ docker-compose up --build
```
+ If you have already started it, the application needs no rebuilding, if there was no code change
```bash
$ docker-compose up 
```

**3) Start using the application in the browser with /web-client/index.html page**


