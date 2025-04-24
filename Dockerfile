# Imagen base de Node.js
FROM node:18

# Directorio de trabajo dentro del contenedor
WORKDIR /app

# Copia los archivos de dependencias
COPY package*.json ./

# Instala dependencias
RUN npm install

# Copia el resto de la app
COPY . .

# Expone el puerto de tu app
EXPOSE 3000

# Comando por defecto
CMD ["node", "index.js"]
