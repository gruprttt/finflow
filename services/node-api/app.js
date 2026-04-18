const express = require("express");
const axios = require("axios");

const app = express();

app.get("/", (req, res) => {
  res.send("Node API running 🚀");
});

// Call Go service
app.get("/go-test", async (req, res) => {
  try {
    const response = await axios.get("http://go-service:8080/hello");

    res.json({
      message: "Got response from Go service",
      data: response.data
    });
  } catch (err) {
    res.status(500).json({ error: err.message });
  }
});

app.listen(3000, () => {
  console.log("Node API running on port 3000");
});
