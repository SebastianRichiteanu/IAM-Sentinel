# IAM-Sentinel 🔐📊

IAM-Sentinel is a tool that transforms AWS IAM (Identity and Access Management) configuration data into a graph structure, enabling security analysts and developers to visualize, query, and audit access relationships using [Neo4j](https://neo4j.com/).

This tool was developed as part of my Master’s thesis in Cyber Security at the University of Bucharest. The full thesis can be accessed [here](Thesis.pdf), which covers the technical model, architecture, algorithms, and security insights this project is built upon.

---

## 🚀 Features

- 🔁 **Graph-Based IAM Mapping**: Converts IAM entities—users, groups, roles, policies, and actions—into a graph model.
- 📊 **Neo4j Integration**: Uses Neo4j for storage, visualization, and graph analytics.
- 🛡️ **Security Audit Ready**: Run centrality analysis, community detection, and permission traceability for security and compliance.
- 🔍 **Custom Queries Support**: Preloaded with example queries for detecting excessive permissions and misconfigurations.
- 🧱 **Modular Architecture**: Cleanly separated components using Go’s idioms for maintainability and flexibility.
- 🐳 **Docker Support**: Includes a `docker-compose` file for easy Neo4j deployment.

---

## 📚 Background

IAM-Sentinel was created in response to the increasing complexity of managing AWS IAM configurations at scale. Traditional IAM analysis methods often fall short when it comes to identifying deeply nested relationships, indirect trust paths, or redundant permissions.

By leveraging graph theory and Neo4j’s Graph Data Science (GDS) library, IAM-Sentinel makes it easy to:

- Visualize AWS IAM structure
- Detect overly permissive policies
- Identify critical users/roles based on centrality
- Explore permission propagation paths

📖 *For an in-depth explanation of the concepts and research behind this tool, see the [Thesis.pdf](Thesis.pdf).*

---

## 🛠️ Getting Started

### 1. Clone the Repository

```bash
git clone https://github.com/SebastianRichiteanu/IAM-Sentinel.git
cd IAM-Sentinel
```

### 2. Start Neo4j with Docker

```bash
docker-compose up -d
```

This will start Neo4j locally on `http://localhost:7474`. The default login is:
- Username: `neo4j`
- Password: `password` (can be changed in `docker-compose.yml`)

### 3. Prepare IAM Data

Export IAM details using the following AWS CLI command:

```bash
aws iam get-account-authorization-details > iam-data.json
```

Place the exported file in the `examples/` folder. (You can modify the path and credentials in the source code.)

### 4. Run IAM-Sentinel

Make sure you have Go installed. Then, simply run:

```bash
go run .
```

This will parse the IAM data, map it into a graph structure, and ingest it into Neo4j.

---

## 💡 Example Queries

Once data is loaded, use the Neo4j UI (`localhost:7474`) or Cypher to run analyses:

```cypher
// 1. List groups with their users
MATCH (g:Group)-[:HAS_MEMBER]->(u:User)
RETURN g, u

// 2. Users with allowed S3 actions
MATCH (u:User)-[r]->(a:Action)
WHERE r.effect = "Allow" AND a.action STARTS WITH "s3:"
RETURN u, a

// 3. Most influential policy by betweenness centrality
CALL gds.betweenness.stream('fullGraph')
YIELD nodeId, score
WITH gds.util.asNode(nodeId) AS node, score
WHERE 'Policy' IN labels(node)
RETURN node.PolicyName, score
ORDER BY score DESC
LIMIT 1
```

---

## 📐 Architecture Overview

IAM-Sentinel is modular and follows the **Single Responsibility Principle**:

- `Sentinel`: Core orchestrator
- `Neo4jConnector`: Handles database connections
- `Parser`: Parses JSON IAM data
- `ResourceMapper`: Transforms data into graph nodes/edges
- `Analyzer`: Runs example queries (optional, mostly for demos)
- `Exporter`: Handles logging/output

---

## 🔍 Use Cases

- **Security Audits**: Spot over-permissive or unused roles/policies
- **Compliance Checks**: Visual proof of access flows and permission hierarchies
- **IAM Optimization**: Identify and remove redundant access paths
- **Versioned Analysis**: Compare IAM states over time using static exports

---

## 📈 Graph Analytics

IAM-Sentinel supports advanced GDS algorithms for deeper insights:

- Centrality Metrics: Degree, Betweenness, Closeness, PageRank
- Community Detection: Louvain, Label Propagation
- Graph Projections
- Triangle Count & Clustering Coefficient

---

## 🧠 Related Reading

- [Neo4j Graph Data Science](https://neo4j.com/product/graph-data-science/)
- [AWS IAM Docs](https://docs.aws.amazon.com/iam/)
- [Go Programming Language](https://golang.org/)
