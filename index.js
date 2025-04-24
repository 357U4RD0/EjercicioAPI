const express = require('express');
const app = express();
app.use(express.json());

const incidentsRoutes = require('./routes/incidents');
app.use('/incidents', incidentsRoutes);

app.listen(3000, () => {
  console.log('Servidor corriendo en http://localhost:3000');
});