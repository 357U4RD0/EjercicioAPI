const express = require('express');
const cors = require('cors');
const path = require('path');
const app = express();

const pool = require('./db');
const incidentsRouter = require('./routes/incidents');

app.use(cors());
app.use(express.json());

// Servir archivos estáticos del frontend
app.use(express.static(path.join(__dirname, 'frontend')));

// Rutas API
app.use('/incidents', incidentsRouter);

const PORT = process.env.PORT || 3000;
app.listen(PORT, () => {
  console.log(`Servidor corriendo en puerto ${PORT}`);
});