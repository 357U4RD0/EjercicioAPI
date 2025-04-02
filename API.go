// ------------------------------------------------------------------------------------------------------------------------------------
//
// He de admitir que para algunos aspectos me apoyé de la inteligencia artificial para que me brinde ideas o pequeños bloques de código
//
// ------------------------------------------------------------------------------------------------------------------------------------

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Incident define la estructura de un incidente
type Incident struct {
	ID          primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Reporter    string             `json:"reporter" bson:"reporter"`
	Description string             `json:"description" bson:"description"`
	Status      string             `json:"status" bson:"status"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
}

// ErrorResponse estructura para respuestas de error
type ErrorResponse struct {
	Message string `json:"message"`
}

var client *mongo.Client
var collection *mongo.Collection

func main() {
	// Conectar a MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Obtener la URL de MongoDB de una variable de entorno o usar una por defecto
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}

	var err error
	client, err = mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal("Error al conectar a MongoDB:", err)
	}

	// Comprobar la conexión
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("Error al hacer ping a MongoDB:", err)
	}
	fmt.Println("Conectado a MongoDB!")

	// Inicializar colección
	collection = client.Database("support").Collection("incidents")

	// Inicializar router Chi
	r := chi.NewRouter()

	// Middleware global
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	// Rutas
	r.Route("/incidents", func(r chi.Router) {
		r.Post("/", createIncident)
		r.Get("/", getIncidents)
		
		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", getIncident)
			r.Put("/", updateIncident)
			r.Delete("/", deleteIncident)
		})
	})

	// Iniciar servidor
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Printf("Servidor iniciado en el puerto %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}

// createIncident maneja la creación de un nuevo incidente
func createIncident(w http.ResponseWriter, r *http.Request) {
	var incident Incident
	err := json.NewDecoder(r.Body).Decode(&incident)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Formato inválido de la solicitud")
		return
	}

	// Validaciones
	if incident.Reporter == "" {
		respondWithError(w, http.StatusBadRequest, "El campo reporter es obligatorio")
		return
	}

	if len(incident.Description) < 10 {
		respondWithError(w, http.StatusBadRequest, "La descripción debe tener al menos 10 caracteres")
		return
	}

	// Establecer valores por defecto
	incident.Status = "pendiente"
	incident.CreatedAt = time.Now()

	// Insertar en la base de datos
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := collection.InsertOne(ctx, incident)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error al crear el incidente")
		return
	}

	incident.ID = result.InsertedID.(primitive.ObjectID)

	respondWithJSON(w, http.StatusCreated, incident)
}

// getIncidents obtiene todos los incidentes
func getIncidents(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error al obtener los incidentes")
		return
	}
	defer cursor.Close(ctx)

	var incidents []Incident
	if err = cursor.All(ctx, &incidents); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error al procesar los incidentes")
		return
	}

	respondWithJSON(w, http.StatusOK, incidents)
}

// getIncident obtiene un incidente por su ID
func getIncident(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "ID de incidente inválido")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var incident Incident
	err = collection.FindOne(ctx, bson.M{"_id": id}).Decode(&incident)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			respondWithError(w, http.StatusNotFound, "Incidente no encontrado")
		} else {
			respondWithError(w, http.StatusInternalServerError, "Error al buscar el incidente")
		}
		return
	}

	respondWithJSON(w, http.StatusOK, incident)
}

// updateIncident actualiza el estado de un incidente
func updateIncident(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "ID de incidente inválido")
		return
	}

	// Solo permitimos actualizar el status
	var updateData struct {
		Status string `json:"status"`
	}

	err = json.NewDecoder(r.Body).Decode(&updateData)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Formato inválido de la solicitud")
		return
	}

	// Validar que el status sea uno de los permitidos
	if updateData.Status != "pendiente" && updateData.Status != "en proceso" && updateData.Status != "resuelto" {
		respondWithError(w, http.StatusBadRequest, "Status inválido. Debe ser 'pendiente', 'en proceso' o 'resuelto'")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Comprobar que el incidente existe
	var existingIncident Incident
	err = collection.FindOne(ctx, bson.M{"_id": id}).Decode(&existingIncident)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			respondWithError(w, http.StatusNotFound, "Incidente no encontrado")
		} else {
			respondWithError(w, http.StatusInternalServerError, "Error al buscar el incidente")
		}
		return
	}

	// Actualizar únicamente el status
	update := bson.M{
		"$set": bson.M{
			"status": updateData.Status,
		},
	}

	_, err = collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error al actualizar el incidente")
		return
	}

	// Obtener el incidente actualizado
	err = collection.FindOne(ctx, bson.M{"_id": id}).Decode(&existingIncident)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error al obtener el incidente actualizado")
		return
	}

	respondWithJSON(w, http.StatusOK, existingIncident)
}

// deleteIncident elimina un incidente por su ID
func deleteIncident(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "ID de incidente inválido")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Comprobar que el incidente existe
	var existingIncident Incident
	err = collection.FindOne(ctx, bson.M{"_id": id}).Decode(&existingIncident)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			respondWithError(w, http.StatusNotFound, "Incidente no encontrado")
		} else {
			respondWithError(w, http.StatusInternalServerError, "Error al buscar el incidente")
		}
		return
	}

	result, err := collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error al eliminar el incidente")
		return
	}

	if result.DeletedCount == 0 {
		respondWithError(w, http.StatusNotFound, "Incidente no encontrado")
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Incidente eliminado correctamente"})
}

// respondWithError envía una respuesta de error en formato JSON
func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, ErrorResponse{Message: message})
}

// respondWithJSON envía una respuesta en formato JSON
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}