   package main

   import (
   	"context"
   	"encoding/json"
   	"fmt"
   	"log"
   	"net/http"
   	"os"

   	"go.mongodb.org/mongo-driver/bson/primitive"
   	"go.mongodb.org/mongo-driver/mongo"
   	"go.mongodb.org/mongo-driver/mongo/options"
   )

   var collection *mongo.Collection

   type Payment struct {
   	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
   	ReservationID string             `bson:"reservation_id" json:"reservation_id"`
   	Amount        float64            `bson:"amount" json:"amount"`
   	Status        string             `bson:"status" json:"status"`
   }

   func initDB() {
   	uri := os.Getenv("MONGO_URI")
   	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
   	if err != nil {
   		log.Fatal(err)
   	}
   	collection = client.Database("payments_db").Collection("transactions")
   	fmt.Println("Conectado a MongoDB")
   }

   func pingHandler(w http.ResponseWriter, r *http.Request) {
   	w.Header().Set("Content-Type", "application/json")
   	w.Write([]byte(`{"status":"pass"}`))
   }

   func createPaymentHandler(w http.ResponseWriter, r *http.Request) {
   	var payment Payment
   	if err := json.NewDecoder(r.Body).Decode(&payment); err != nil {
   		http.Error(w, err.Error(), http.StatusBadRequest)
   		return
   	}
   	payment.Status = "COMPLETED"
   	payment.ID = primitive.NewObjectID()
   	_, err := collection.InsertOne(context.TODO(), payment)
   	if err != nil {
   		http.Error(w, "Error guardando el pago", http.StatusInternalServerError)
   		return
   	}
   	w.Header().Set("Content-Type", "application/json")
   	json.NewEncoder(w).Encode(payment)
   }

   func main() {
   	initDB()
   	http.HandleFunc("/payments", createPaymentHandler)
   	http.HandleFunc("/payments/ping", pingHandler)
   	fmt.Println("ms-payments corriendo en el puerto 8080")
   	log.Fatal(http.ListenAndServe(":8080", nil))
   }