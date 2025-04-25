const API_URL = 'http://localhost:3000/incidents';
const contenido = document.getElementById('contenido');

async function getAllIncidents() {
  try {
    const res = await fetch(API_URL);
    const data = await res.json();
    contenido.innerHTML = data.map(incidente => `
      <div class="card">
        <h3>${incidente.title}</h3>
        <p>${incidente.description}</p>
        <span><strong>Estado:</strong> ${incidente.status}</span>
      </div>
    `).join('');
  } catch (error) {
    console.error('Error al obtener incidentes:', error);
  }
}

function mostrarFormularioCrear() {
  contenido.innerHTML = `
    <h2>Crear incidente</h2>
    <input type="text" id="titulo" placeholder="Título"><br>
    <input type="text" id="descripcion" placeholder="Descripción"><br>
    <label for="Estado">Estado:</label>
    <select id="Estado">
      <option value="pendiente">Pendiente</option>
      <option value="en progreso">En progreso</option>
      <option value="resuelto">Resuelto</option>
    <button onclick="crearIncidente()">Enviar</button>
  `;
}

async function crearIncidente() {
  const titulo = document.getElementById('titulo').value;
  const descripcion = document.getElementById('descripcion').value;
  const estado = document.getElementById('estado').value;

  try {
    await fetch(API_URL, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ title: titulo, description: descripcion, status: estado })
    });
    getAllIncidents();
  } catch (error) {
    console.error('Error al crear incidente:', error);
  }
}

function mostrarBuscarPorId() {
  contenido.innerHTML = `
    <h2>Buscar incidente por ID</h2>
    <input type="number" id="idBuscar" placeholder="ID"><br>
    <button onclick="buscarPorId()">Buscar</button>
  `;
}

async function buscarPorId() {
  const id = document.getElementById('idBuscar').value;
  try {
    const res = await fetch(`${API_URL}/${id}`);
    const data = await res.json();
    contenido.innerHTML = `
      <div class="card">
        <h3>${data.title}</h3>
        <p>${data.description}</p>
        <span><strong>Estado:</strong> ${data.status}</span>
      </div>
    `;
  } catch (error) {
    console.error('Error al buscar incidente:', error);
  }
}

function mostrarFormularioActualizar() {
  contenido.innerHTML = `
    <h2>Actualizar estado de incidente</h2>
    <input type="number" id="idActualizar" placeholder="ID del incidente"><br>
    <label for="estadoActualizar">Nuevo estado:</label>
    <select id="estadoActualizar">
      <option value="pendiente">Pendiente</option>
      <option value="en progreso">En progreso</option>
      <option value="resuelto">Resuelto</option>
    </select><br>
    <button onclick="actualizarIncidente()">Actualizar</button>
  `;
}

async function actualizarIncidente() {
  const id = document.getElementById('idActualizar').value;
  const estado = document.getElementById('estadoActualizar').value;

  try {
    await fetch(`${API_URL}/${id}`, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ status: estado })
    });
    getAllIncidents();
  } catch (error) {
    console.error('Error al actualizar incidente:', error);
  }
}


function mostrarFormularioEliminar() {
  contenido.innerHTML = `
    <h2>Eliminar incidente</h2>
    <input type="number" id="idEliminar" placeholder="ID del incidente"><br>
    <button onclick="eliminarIncidente()">Eliminar</button>
  `;
}

async function eliminarIncidente() {
  const id = document.getElementById('idEliminar').value;
  try {
    await fetch(`${API_URL}/${id}`, { method: 'DELETE' });
    getAllIncidents();
  } catch (error) {
    console.error('Error al eliminar incidente:', error);
  }
}

function mostrarFormularioBuscar() {
    const contenido = document.getElementById('contenido');
    contenido.innerHTML = `
      <h2>Buscar incidente por ID</h2>
      <input type="number" id="buscarId" placeholder="ID del incidente">
      <button onclick="buscarPorId()">Buscar</button>
      <div id="resultadoBuscar"></div>
    `;
  }
  
  async function buscarPorId() {
    const id = document.getElementById('buscarId').value;
    const resultado = document.getElementById('resultadoBuscar');
  
    if (!id) {
      resultado.innerHTML = `<p style="color: red;">Ingresa un ID válido</p>`;
      return;
    }
  
    try {
      const response = await fetch(`http://localhost:3000/incidents/${id}`);
      if (!response.ok) {
        resultado.innerHTML = `<p style="color: red;">Incidente no encontrado</p>`;
        return;
      }
  
      const data = await response.json();
      resultado.innerHTML = `
        <div class="incidente">
          <h3>${data.title}</h3>
          <p><strong>Descripción:</strong> ${data.description}</p>
          <p><strong>Estado:</strong> ${data.status}</p>
        </div>
      `;
    } catch (error) {
      console.error("Error al buscar incidente:", error);
      resultado.innerHTML = `<p style="color: red;">Error al conectar con la API</p>`;
    }
  }  