const express = require("express");
const axios = require("axios");

const app = express();
app.use(express.json());

app.get("/", (req, res) => {
  res.send("Node API running 🚀");
});

app.post("/users", async (req, res) => {
  const response = await axios.post("http://go-service:8080/users", req.body);
  res.json(response.data);
});

app.post("/orders", async (req, res) => {
  const response = await axios.post("http://go-service:8080/orders", req.body);
  res.json(response.data);
});

app.get("/orders/:id", async (req, res) => {
  const response = await axios.get(`http://go-service:8080/orders/${req.params.id}`);
  res.send(response.data);
});

app.listen(3000, () => console.log("Node API running"));
