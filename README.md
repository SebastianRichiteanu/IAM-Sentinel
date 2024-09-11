# IAM-Sentinel

IAM-Sentinel facilitates the integration of AWS IAM data into neo4j.

This project was developed as part of my masters's thesis in Cyber Security. For a detailed exploration of the concepts and methodologies behind this project, you can find the full thesis paper [here](Thesis.pdf).

## How to use?

The tool uses a neo4j database to store the information. One can be spinned up using docker: ```docker-compose up -d```.

The ingestion files have to be placed inside the `examples` folder, you can change this folder inside the code.

To run the tool simply execute `go run .`.

One can query and visualize the data in Neo4j's UI, the default address is `localhost:7474`. 
