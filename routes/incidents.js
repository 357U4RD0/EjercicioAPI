const express = require('express');
const router = express.Router();
const pool = require('../db');

// Crear un nuevo incidente
router.post('/', async (req, res) => {
  const { title, description, status } = req.body;
  try {
    const result = await pool.query(
      'INSERT INTO incidents (title, description, status) VALUES ($1, $2, $3) RETURNING *',
      [title, description, status]
    );
    res.status(201).json(result.rows[0]);
  } catch (err) {
    res.status(500).json({ error: err.message });
  }
});

// Listar Incidentes
router.get('/', async (req, res) => {
    try {
      const result = await pool.query('SELECT * FROM incidents');
      res.json(result.rows);
    } catch (err) {
      res.status(500).json({ error: err.message });
    }
  });
  
  // Obtener por ID
router.get('/:id', async (req, res) => {
    const { id } = req.params;
    try {
      const result = await pool.query('SELECT * FROM incidents WHERE id = $1', [id]);
      if (result.rows.length === 0) return res.status(404).json({ error: 'Incidente no encontrado' });
      res.json(result.rows[0]);
    } catch (err) {
      res.status(500).json({ error: err.message });
    }
  });
  
// Actualizar estado
router.put('/:id', async (req, res) => {
    const { id } = req.params;
    const { status } = req.body;
    try {
      const result = await pool.query(
        'UPDATE incidents SET status = $1 WHERE id = $2 RETURNING *',
        [status, id]
      );
      if (result.rows.length === 0) return res.status(404).json({ error: 'Incidente no encontrado' });
      res.json(result.rows[0]);
    } catch (err) {
      res.status(500).json({ error: err.message });
    }
  });
  
module.exports = router;