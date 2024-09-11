# IAM-Sentinel

IAM-Sentinel facilitates the integration of AWS IAM data into neo4j.

This project was developed as part of my masters's thesis in Cyber Security. For a detailed exploration of the concepts and methodologies behind this project, you can find the full thesis paper [here](Thesis.pdf).

## Usage

The tool uses a neo4j database to store the information. One can be spinned up using docker and the provided docker-compose file: ```docker-compose up -d```.

The ingestion files have to be placed inside the `examples` folder, you can change this folder inside the code, alongside with the database credentials.

To run the tool simply execute `go run .`

Afterwards, the data can be queried and visualized in Neo4j's UI, the default address is `localhost:7474`. 
