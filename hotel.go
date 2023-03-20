package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	//"io/ioutil"
	"log"
	"net/http"

	//"strconv"
	//"strings"

	_ "github.com/go-sql-driver/mysql"
)

// Structure du client

type Client struct {
	ID        int    `json:"id"`
	Nom       string `json:"nom"`
	Prenom    string `json:"prenom"`
	Adresse   string `json:"adresse"`
	Telephone string `json:"telephone"`
	Email     string `json:"email"`
}
type Reservation struct {
	ID         int     `json:"id"`
	DateEntree string  `json:"dateEntree"`
	DateSortie string  `json:"dateSortie"`
	PrixTotal  float32 `json:"prixTotal"`
}

////////////////////////////////////////////////////////////////////////////////////////////////////////
//Connexion Mysql

func dbConnect() (*sql.DB, error) {

	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/geshotel")
	if err != nil {
		return nil, err
	}
	return db, nil
}

// //////////////////////////////////////////////////////////////////////////////////////////////////

func getClients(w http.ResponseWriter, r *http.Request) {
	// Connexion à la base de données
	db, err := sql.Open("mysql", "root:@tcp(localhost:3306)/geshotel")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Erreur lors de la connexion à la base de données")
		return
	}
	defer db.Close()

	// Récupération de tous les clients depuis la table "clients"
	rows, err := db.Query("SELECT id, nom, prenom, adresse, telephone, email FROM client")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Erreur lors de la récupération des clients depuis la base de données")
		return
	}
	defer rows.Close()

	// Boucle sur les résultats et stockage dans une slice de clients
	clients := make([]Client, 0)
	for rows.Next() {
		var client Client
		err := rows.Scan(&client.ID, &client.Nom, &client.Prenom, &client.Adresse, &client.Telephone, &client.Email)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintln(w, "Erreur lors de la lecture des données du client depuis la base de données")
			return
		}
		clients = append(clients, client)
	}
	if err := rows.Err(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Erreur lors de la récupération des clients depuis la base de données")
		return
	}

	// Envoi des clients en format JSON dans la réponse HTTP
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(clients)
}

// //////////
func getReservations(w http.ResponseWriter, r *http.Request) {
	// Connexion à la base de données
	db, err := sql.Open("mysql", "root:@tcp(localhost:3306)/geshotel")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Erreur lors de la connexion à la base de données")
		return
	}
	defer db.Close()

	// Récupération de tous les clients depuis la table "reservation"
	rows, err := db.Query("SELECT * FROM reservation")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Erreur lors de la récupération des reservations depuis la base de données")
		return
	}
	defer rows.Close()

	// Boucle sur les résultats et stockage dans une slice de clients
	reservations := make([]Reservation, 0)
	for rows.Next() {
		var reservation Reservation
		err := rows.Scan(&reservation.ID, &reservation.DateEntree, &reservation.DateSortie, &reservation.PrixTotal)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintln(w, "Erreur lors de la lecture des données de la reservation depuis la base de données")
			return
		}
		reservations = append(reservations, reservation)
	}
	if err := rows.Err(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Erreur lors de la récupération des reservations depuis la base de données")
		return
	}

	// Envoi des clients en format JSON dans la réponse HTTP
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(reservations)
}

// /////////////////////////////////////////////////////////////////////////////////////////////
// un client
func getClient(clientID int) (Client, error) {
	// Ouvre une connexion à la base de données MySQL
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/geshotel")
	if err != nil {
		return Client{}, err
	}
	defer db.Close()

	// Prépare et exécute la requête de récupération du client depuis la base de données
	query := "SELECT id, nom, prenom, adresse, telephone, email FROM client WHERE id=?"
	row := db.QueryRow(query, clientID)

	// Récupère les données du client depuis la ligne de résultat
	var client Client
	err = row.Scan(&client.ID, &client.Nom, &client.Prenom, &client.Adresse, &client.Telephone, &client.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return Client{}, fmt.Errorf("client with ID %d not found", clientID)
		}
		return Client{}, err
	}

	return client, nil
}

func getClientHandler(w http.ResponseWriter, r *http.Request) {
	// Récupère l'ID du client depuis l'URL
	clientID, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/clients/"))
	if err != nil {
		http.Error(w, "Invalid client ID", http.StatusBadRequest)
		return
	}

	// Récupère le client correspondant depuis la base de données
	client, err := getClient(clientID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Encode le client en JSON et renvoie la réponse HTTP
	json.NewEncoder(w).Encode(client)
}

// ////////
func getReservation(reservationID int) (Reservation, error) {
	// Ouvre une connexion à la base de données MySQL
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/geshotel")
	if err != nil {
		return Reservation{}, err
	}
	defer db.Close()

	// Prépare et exécute la requête de récupération du reservation depuis la base de données
	query := "SELECT * FROM reservation WHERE id=?"
	row := db.QueryRow(query, reservationID)

	// Récupère les données du reservation depuis la ligne de résultat
	var reservation Reservation
	err = row.Scan(&reservation.DateEntree, &reservation.DateSortie, &reservation.PrixTotal)
	if err != nil {
		if err == sql.ErrNoRows {
			return Reservation{}, fmt.Errorf("reservation with ID %d not found", reservationID)
		}
		return Reservation{}, err
	}

	return reservation, nil
}

func getReservationHandler(w http.ResponseWriter, r *http.Request) {
	// Récupère l'ID du reservation depuis l'URL
	reservationID, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/reservations/"))
	if err != nil {
		http.Error(w, "Invalid reservation ID", http.StatusBadRequest)
		return
	}

	// Récupère le reservation correspondant depuis la base de données
	reservation, err := getReservation(reservationID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Encode le reservation en JSON et renvoie la réponse HTTP
	json.NewEncoder(w).Encode(reservation)
}

/////////////////////////////////////////////////////////////////////////////////////////////
//Création d'un nouveau client

func createClientHandler(w http.ResponseWriter, r *http.Request) {
	// Lecture des données du client à partir du corps de la requête
	var client Client
	err := json.NewDecoder(r.Body).Decode(&client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Connexion à la base de données
	db, err := dbConnect()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Exécution de la requête d'insertion
	result, err := db.Exec("INSERT INTO client (nom, prenom, adresse, telephone, email ) VALUES (?,?, ?, ?, ?)", client.Nom, client.Prenom, client.Adresse, client.Telephone, client.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Récupération de l'ID généré automatiquement
	id, err := result.LastInsertId()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Création de la réponse
	client.ID = int(id)
	response, err := json.Marshal(client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Envoi de la réponse
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(response)
}

// ///////
func createReservationHandler(w http.ResponseWriter, r *http.Request) {
	// Lecture des données du reservation à partir du corps de la requête
	var reservation Reservation
	err := json.NewDecoder(r.Body).Decode(&reservation)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Connexion à la base de données
	db, err := dbConnect()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Exécution de la requête d'insertion
	result, err := db.Exec("INSERT INTO reservation (id, dateEntree, dateSortie, prixTotal) VALUES (?,?, ?, ?, ?)", reservation.ID, reservation.DateEntree, reservation.DateSortie, reservation.PrixTotal)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Récupération de l'ID généré automatiquement
	id, err := result.LastInsertId()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Création de la réponse
	reservation.ID = int(id)
	response, err := json.Marshal(reservation)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Envoi de la réponse
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(response)
}

// //////////////////////////////////////////////////////////////////////////////////////////////
func updateClientHandler(w http.ResponseWriter, r *http.Request) {
	// Récupère l'ID du client à mettre à jour depuis les paramètres de la requête
	clientIDStr := r.URL.Query().Get("id")
	clientID, err := strconv.Atoi(clientIDStr)
	if err != nil {
		http.Error(w, "Invalid client ID", http.StatusBadRequest)
		return
	}

	// Récupère les données du client à partir du corps de la requête
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	var updatedClient Client
	err = json.Unmarshal(reqBody, &updatedClient)
	if err != nil {
		http.Error(w, "Failed to parse request body", http.StatusBadRequest)
		return
	}

	// Ouvre une connexion à la base de données MySQL
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/geshotel")
	if err != nil {
		http.Error(w, "Failed to connect to database", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Prépare et exécute la requête de mise à jour du client dans la base de données
	query := "UPDATE client SET nom=?, prenom=?, adresse=?, telephone=?, email=? WHERE id=?"
	_, err = db.Exec(query, updatedClient.Nom, updatedClient.Prenom, updatedClient.Adresse, updatedClient.Telephone, updatedClient.Email, clientID)
	if err != nil {
		http.Error(w, "Failed to update client in database", http.StatusInternalServerError)
		return
	}

	// Retourne une réponse HTTP 200 OK
	w.WriteHeader(http.StatusOK)
}

// /////////////////////////////////////////////////////////////////////////////////////////
func deleteClientHandler(w http.ResponseWriter, r *http.Request) {
	// Récupérer l'id du client à supprimer à partir de la requête URL
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid client ID", http.StatusBadRequest)
		return
	}

	// Ouvrir la connexion à la base de données
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/geshotel")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Préparer la requête SQL pour supprimer le client
	stmt, err := db.Prepare("DELETE FROM client WHERE id = ?")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	// Exécuter la requête SQL pour supprimer le client
	res, err := stmt.Exec(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Vérifier si le client a été supprimé avec succès
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if rowsAffected == 0 {
		http.Error(w, "Client not found", http.StatusNotFound)
		return
	}

	// Renvoyer une réponse JSON pour indiquer que le client a été supprimé avec succès
	response := map[string]interface{}{
		"status":  "success",
		"message": "Client deleted successfully",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func updateReservationHandler(w http.ResponseWriter, r *http.Request) {
	// Récupère l'ID du reservation à mettre à jour depuis les paramètres de la requête
	reservationIDStr := r.URL.Query().Get("id")
	reservationID, err := strconv.Atoi(reservationIDStr)
	if err != nil {
		http.Error(w, "Invalid reservation ID", http.StatusBadRequest)
		return
	}

	// Récupère les données du reservation à partir du corps de la requête
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	var updatedReservation Reservation
	err = json.Unmarshal(reqBody, &updatedReservation)
	if err != nil {
		http.Error(w, "Failed to parse request body", http.StatusBadRequest)
		return
	}

	// Ouvre une connexion à la base de données MySQL
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/geshotel")
	if err != nil {
		http.Error(w, "Failed to connect to database", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Prépare et exécute la requête de mise à jour du reservation dans la base de données
	query := "UPDATE reservation SET nom=?, prenom=?, adresse=?, telephone=?, email=? WHERE id=?"
	_, err = db.Exec(query, updatedReservation.DateEntree, updatedReservation.DateSortie, updatedReservation.PrixTotal, reservationID)
	if err != nil {
		http.Error(w, "Failed to update reservation in database", http.StatusInternalServerError)
		return
	}

	// Retourne une réponse HTTP 200 OK
	w.WriteHeader(http.StatusOK)
}

// /////////////////////////////////////////////////////////////////////////////////////////
func deleteReservationHandler(w http.ResponseWriter, r *http.Request) {
	// Récupérer l'id du reservation à supprimer à partir de la requête URL
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid reservation ID", http.StatusBadRequest)
		return
	}

	// Ouvrir la connexion à la base de données
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/geshotel")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Préparer la requête SQL pour supprimer le reservation
	stmt, err := db.Prepare("DELETE FROM reservation WHERE id = ?")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	// Exécuter la requête SQL pour supprimer le reservation
	res, err := stmt.Exec(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Vérifier si le reservation a été supprimé avec succès
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if rowsAffected == 0 {
		http.Error(w, "Reservation not found", http.StatusNotFound)
		return
	}

	// Renvoyer une réponse JSON pour indiquer que le reservation a été supprimé avec succès
	response := map[string]interface{}{
		"status":  "success",
		"message": "Reservation deleted successfully",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Structure de la chambre

type Chambre struct {
	ID     int    `json:"id"`
	Etat   string `json:"etat"`
	Classe string `json:"classe"`
}

// //////////////////////////////////////////////////////////////////////////////////////////////////

func getChambres(w http.ResponseWriter, r *http.Request) {
	// Connexion à la base de données
	db, err := sql.Open("mysql", "root:@tcp(localhost:3306)/geshotel")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Erreur lors de la connexion à la base de données")
		return
	}
	defer db.Close()

	// Récupération de tous les chambres depuis la table "chambres"
	rows, err := db.Query("SELECT id, etat, classe FROM chambre")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Erreur lors de la récupération des chambres depuis la base de données")
		return
	}
	defer rows.Close()

	// Boucle sur les résultats et stockage dans une slice de chambres
	chambres := make([]Chambre, 0)
	for rows.Next() {
		var chambre Chambre
		err := rows.Scan(&chambre.ID, &chambre.Etat, &chambre.Classe)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintln(w, "Erreur lors de la lecture des données de la chambre depuis la base de données")
			return
		}
		chambres = append(chambres, chambre)
	}
	if err := rows.Err(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Erreur lors de la récupération des chambres depuis la base de données")
		return
	}

	// Envoi des chambres en format JSON dans la réponse HTTP
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(chambres)
}

// /////////////////////////////////////////////////////////////////////////////////////////////
// un chambre
func getChambre(chambreID int) (Chambre, error) {
	// Ouvre une connexion à la base de données MySQL
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/geshotel")
	if err != nil {
		return Chambre{}, err
	}
	defer db.Close()

	// Prépare et exécute la requête de récupération du chambre depuis la base de données
	query := "SELECT id, nom, prenom, adresse, telephone, email FROM chambre WHERE id=?"
	row := db.QueryRow(query, chambreID)

	// Récupère les données du chambre depuis la ligne de résultat
	var chambre Chambre
	err = row.Scan(&chambre.ID, &chambre.Etat, &chambre.Classe)
	if err != nil {
		if err == sql.ErrNoRows {
			return Chambre{}, fmt.Errorf("chambre with ID %d not found", chambreID)
		}
		return Chambre{}, err
	}

	return chambre, nil
}

func getChambreHandler(w http.ResponseWriter, r *http.Request) {
	// Récupère l'ID du chambre depuis l'URL
	chambreID, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/chambres/"))
	if err != nil {
		http.Error(w, "Invalid chambre ID", http.StatusBadRequest)
		return
	}

	// Récupère le chambre correspondant depuis la base de données
	chambre, err := getChambre(chambreID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Encode le chambre en JSON et renvoie la réponse HTTP
	json.NewEncoder(w).Encode(chambre)
}

/////////////////////////////////////////////////////////////////////////////////////////////
//Création d'un nouveau chambre

func createChambreHandler(w http.ResponseWriter, r *http.Request) {
	// Lecture des données du chambre à partir du corps de la requête
	var chambre Chambre
	err := json.NewDecoder(r.Body).Decode(&chambre)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Connexion à la base de données
	db, err := dbConnect()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Exécution de la requête d'insertion
	result, err := db.Exec("INSERT INTO chambre ( etat, classe ) VALUES (?,?)", chambre.Etat, chambre.Classe)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Récupération de l'ID généré automatiquement
	id, err := result.LastInsertId()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Création de la réponse
	chambre.ID = int(id)
	response, err := json.Marshal(chambre)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Envoi de la réponse
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(response)
}

// //////////////////////////////////////////////////////////////////////////////////////////////
func updateChambreHandler(w http.ResponseWriter, r *http.Request) {
	// Récupère l'ID du chambre à mettre à jour depuis les paramètres de la requête
	chambreIDStr := r.URL.Query().Get("id")
	chambreID, err := strconv.Atoi(chambreIDStr)
	if err != nil {
		http.Error(w, "Invalid chambre ID", http.StatusBadRequest)
		return
	}

	// Récupère les données du chambre à partir du corps de la requête
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	var updatedChambre Chambre
	err = json.Unmarshal(reqBody, &updatedChambre)
	if err != nil {
		http.Error(w, "Failed to parse request body", http.StatusBadRequest)
		return
	}

	// Ouvre une connexion à la base de données MySQL
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/geshotel")
	if err != nil {
		http.Error(w, "Failed to connect to database", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Prépare et exécute la requête de mise à jour du chambre dans la base de données
	query := "UPDATE chambre SET nom=?, prenom=?, adresse=?, telephone=?, email=? WHERE id=?"
	_, err = db.Exec(query, updatedChambre.Etat, updatedChambre.Classe, chambreID)
	if err != nil {
		http.Error(w, "Failed to update chambre in database", http.StatusInternalServerError)
		return
	}

	// Retourne une réponse HTTP 200 OK
	w.WriteHeader(http.StatusOK)
}

// /////////////////////////////////////////////////////////////////////////////////////////
func deleteChambreHandler(w http.ResponseWriter, r *http.Request) {
	// Récupérer l'id du chambre à supprimer à partir de la requête URL
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid chambre ID", http.StatusBadRequest)
		return
	}

	// Ouvrir la connexion à la base de données
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/geshotel")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Préparer la requête SQL pour supprimer le chambre
	stmt, err := db.Prepare("DELETE FROM chambre WHERE id = ?")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	// Exécuter la requête SQL pour supprimer le chambre
	res, err := stmt.Exec(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Vérifier si le chambre a été supprimé avec succès
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if rowsAffected == 0 {
		http.Error(w, "Chambre not found", http.StatusNotFound)
		return
	}

	// Renvoyer une réponse JSON pour indiquer que le chambre a été supprimé avec succès
	response := map[string]interface{}{
		"status":  "success",
		"message": "Chambre deleted successfully",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// //////////////////////////////////////////////////////////////////////////////

// Structure du hotel

type Hotel struct {
	ID      int    `json:"id"`
	Nom     string `json:"nom"`
	Adresse string `json:"adresse"`
}

// //////////////////////////////////////////////////////////////////////////////////////////////////

func getHotels(w http.ResponseWriter, r *http.Request) {
	// Connexion à la base de données
	db, err := sql.Open("mysql", "root:@tcp(localhost:3306)/geshotel")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Erreur lors de la connexion à la base de données")
		return
	}
	defer db.Close()

	// Récupération de tous les hotels depuis la table "hotels"
	rows, err := db.Query("SELECT id, nom, prenom, adresse, telephone, email FROM hotel")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Erreur lors de la récupération des hotels depuis la base de données")
		return
	}
	defer rows.Close()

	// Boucle sur les résultats et stockage dans une slice de hotels
	hotels := make([]Hotel, 0)
	for rows.Next() {
		var hotel Hotel
		err := rows.Scan(&hotel.ID, &hotel.Nom, &hotel.Adresse)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintln(w, "Erreur lors de la lecture des données du hotel depuis la base de données")
			return
		}
		hotels = append(hotels, hotel)
	}
	if err := rows.Err(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Erreur lors de la récupération des hotels depuis la base de données")
		return
	}

	// Envoi des hotels en format JSON dans la réponse HTTP
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(hotels)
}

// /////////////////////////////////////////////////////////////////////////////////////////////
// un hotel
func getHotel(hotelID int) (Hotel, error) {
	// Ouvre une connexion à la base de données MySQL
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/geshotel")
	if err != nil {
		return Hotel{}, err
	}
	defer db.Close()

	// Prépare et exécute la requête de récupération du hotel depuis la base de données
	query := "SELECT * FROM hotel WHERE id=?"
	row := db.QueryRow(query, hotelID)

	// Récupère les données du hotel depuis la ligne de résultat
	var hotel Hotel
	err = row.Scan(&hotel.ID, &hotel.Nom, &hotel.Adresse)
	if err != nil {
		if err == sql.ErrNoRows {
			return Hotel{}, fmt.Errorf("hotel with ID %d not found", hotelID)
		}
		return Hotel{}, err
	}

	return hotel, nil
}

func getHotelHandler(w http.ResponseWriter, r *http.Request) {
	// Récupère l'ID du hotel depuis l'URL
	hotelID, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/hotels/"))
	if err != nil {
		http.Error(w, "Invalid hotel ID", http.StatusBadRequest)
		return
	}

	// Récupère le hotel correspondant depuis la base de données
	hotel, err := getHotel(hotelID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Encode le hotel en JSON et renvoie la réponse HTTP
	json.NewEncoder(w).Encode(hotel)
}

/////////////////////////////////////////////////////////////////////////////////////////////
//Création d'un nouveau hotel

func createHotelHandler(w http.ResponseWriter, r *http.Request) {
	// Lecture des données du hotel à partir du corps de la requête
	var hotel Hotel
	err := json.NewDecoder(r.Body).Decode(&hotel)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Connexion à la base de données
	db, err := dbConnect()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Exécution de la requête d'insertion
	result, err := db.Exec("INSERT INTO hotel (nom, adresse ) VALUES (?,?)", hotel.Nom, hotel.Adresse)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Récupération de l'ID généré automatiquement
	id, err := result.LastInsertId()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Création de la réponse
	hotel.ID = int(id)
	response, err := json.Marshal(hotel)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Envoi de la réponse
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(response)
}

// //////////////////////////////////////////////////////////////////////////////////////////////
func updateHotelHandler(w http.ResponseWriter, r *http.Request) {
	// Récupère l'ID du hotel à mettre à jour depuis les paramètres de la requête
	hotelIDStr := r.URL.Query().Get("id")
	hotelID, err := strconv.Atoi(hotelIDStr)
	if err != nil {
		http.Error(w, "Invalid hotel ID", http.StatusBadRequest)
		return
	}

	// Récupère les données du hotel à partir du corps de la requête
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	var updatedHotel Hotel
	err = json.Unmarshal(reqBody, &updatedHotel)
	if err != nil {
		http.Error(w, "Failed to parse request body", http.StatusBadRequest)
		return
	}

	// Ouvre une connexion à la base de données MySQL
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/geshotel")
	if err != nil {
		http.Error(w, "Failed to connect to database", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Prépare et exécute la requête de mise à jour du hotel dans la base de données
	query := "UPDATE hotel SET nom=?,  adresse=? WHERE id=?"
	_, err = db.Exec(query, updatedHotel.Nom, updatedHotel.Adresse, hotelID)
	if err != nil {
		http.Error(w, "Failed to update hotel in database", http.StatusInternalServerError)
		return
	}

	// Retourne une réponse HTTP 200 OK
	w.WriteHeader(http.StatusOK)
}

// /////////////////////////////////////////////////////////////////////////////////////////
func deleteHotelHandler(w http.ResponseWriter, r *http.Request) {
	// Récupérer l'id du hotel à supprimer à partir de la requête URL
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid hotel ID", http.StatusBadRequest)
		return
	}

	// Ouvrir la connexion à la base de données
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/geshotel")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Préparer la requête SQL pour supprimer le hotel
	stmt, err := db.Prepare("DELETE FROM hotel WHERE id = ?")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	// Exécuter la requête SQL pour supprimer le hotel
	res, err := stmt.Exec(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Vérifier si le hotel a été supprimé avec succès
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if rowsAffected == 0 {
		http.Error(w, "Hotel not found", http.StatusNotFound)
		return
	}

	// Renvoyer une réponse JSON pour indiquer que le hotel a été supprimé avec succès
	response := map[string]interface{}{
		"status":  "success",
		"message": "Hotel deleted successfully",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Structure du service

type Service struct {
	ID        int    `json:"id"`
	PtitDej   bool   `json:"ptitdej"`
	Telephone string `json:"telephone"`
	Bar       bool   `json:"bar"`
}

// //////////////////////////////////////////////////////////////////////////////////////////////////

func getServices(w http.ResponseWriter, r *http.Request) {
	// Connexion à la base de données
	db, err := sql.Open("mysql", "root:@tcp(localhost:3306)/gesservice")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Erreur lors de la connexion à la base de données")
		return
	}
	defer db.Close()

	// Récupération de tous les services depuis la table "services"
	rows, err := db.Query("SELECT * FROM service")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Erreur lors de la récupération des services depuis la base de données")
		return
	}
	defer rows.Close()

	// Boucle sur les résultats et stockage dans une slice de services
	services := make([]Service, 0)
	for rows.Next() {
		var service Service
		err := rows.Scan(&service.ID, &service.PtitDej, &service.Telephone, &service.Bar)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintln(w, "Erreur lors de la lecture des données du service depuis la base de données")
			return
		}
		services = append(services, service)
	}
	if err := rows.Err(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Erreur lors de la récupération des services depuis la base de données")
		return
	}

	// Envoi des services en format JSON dans la réponse HTTP
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(services)
}

// /////////////////////////////////////////////////////////////////////////////////////////////
// un service
func getService(serviceID int) (Service, error) {
	// Ouvre une connexion à la base de données MySQL
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/gesservice")
	if err != nil {
		return Service{}, err
	}
	defer db.Close()

	// Prépare et exécute la requête de récupération du service depuis la base de données
	query := "SELECT id, nom, prenom, adresse, telephone, email FROM service WHERE id=?"
	row := db.QueryRow(query, serviceID)

	// Récupère les données du service depuis la ligne de résultat
	var service Service
	err = row.Scan(&service.ID, &service.PtitDej, &service.Telephone, &service.Bar)
	if err != nil {
		if err == sql.ErrNoRows {
			return Service{}, fmt.Errorf("service with ID %d not found", serviceID)
		}
		return Service{}, err
	}

	return service, nil
}

func getServiceHandler(w http.ResponseWriter, r *http.Request) {
	// Récupère l'ID du service depuis l'URL
	serviceID, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/services/"))
	if err != nil {
		http.Error(w, "Invalid service ID", http.StatusBadRequest)
		return
	}

	// Récupère le service correspondant depuis la base de données
	service, err := getService(serviceID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Encode le service en JSON et renvoie la réponse HTTP
	json.NewEncoder(w).Encode(service)
}

/////////////////////////////////////////////////////////////////////////////////////////////
//Création d'un nouveau service

func createServiceHandler(w http.ResponseWriter, r *http.Request) {
	// Lecture des données du service à partir du corps de la requête
	var service Service
	err := json.NewDecoder(r.Body).Decode(&service)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Connexion à la base de données
	db, err := dbConnect()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Exécution de la requête d'insertion
	result, err := db.Exec("INSERT INTO service (ptitDej, telephone, bar ) VALUES (?,?, ?)", service.PtitDej, service.Telephone, service.Bar)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Récupération de l'ID généré automatiquement
	id, err := result.LastInsertId()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Création de la réponse
	service.ID = int(id)
	response, err := json.Marshal(service)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Envoi de la réponse
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(response)
}

/* //////////////////////////////////////////////////////////////////////////////////////////////
func updateServiceHandler(w http.ResponseWriter, r *http.Request) {
	// Récupère l'ID du service à mettre à jour depuis les paramètres de la requête
	serviceIDStr := r.URL.Query().Get("id")
	serviceID, err := strconv.Atoi(serviceIDStr)
	if err != nil {
		http.Error(w, "Invalid service ID", http.StatusBadRequest)
		return
	}

	// Récupère les données du service à partir du corps de la requête
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	var updatedService Service
	err = json.Unmarshal(reqBody, &updatedService)
	if err != nil {
		http.Error(w, "Failed to parse request body", http.StatusBadRequest)
		return
	}

	// Ouvre une connexion à la base de données MySQL
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/gesservice")
	if err != nil {
		http.Error(w, "Failed to connect to database", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Prépare et exécute la requête de mise à jour du service dans la base de données
	query := "UPDATE service SET ptitDej=?, telephone=?, bar=? WHERE id=?"
	_, err = db.Exec(query, updatedService.PtitDej, updatedService.Telephone, updatedService.Bar)
	if err != nil {
		http.Error(w, "Failed to update service in database", http.StatusInternalServerError)
		return
	}

	// Retourne une réponse HTTP 200 OK
	w.WriteHeader(http.StatusOK)
}

*/
/////////////////////////////////////////////////////////////////////////////////////////
func deleteServiceHandler(w http.ResponseWriter, r *http.Request) {
	// Récupérer l'id du service à supprimer à partir de la requête URL
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid service ID", http.StatusBadRequest)
		return
	}

	// Ouvrir la connexion à la base de données
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/gesservice")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Préparer la requête SQL pour supprimer le service
	stmt, err := db.Prepare("DELETE FROM service WHERE id = ?")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	// Exécuter la requête SQL pour supprimer le service
	res, err := stmt.Exec(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Vérifier si le service a été supprimé avec succès
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if rowsAffected == 0 {
		http.Error(w, "Service not found", http.StatusNotFound)
		return
	}

	// Renvoyer une réponse JSON pour indiquer que le service a été supprimé avec succès
	response := map[string]interface{}{
		"status":  "success",
		"message": "Service deleted successfully",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// /////////////////////////////////////////////////////////////////////////////////////////////
func main() {

	http.HandleFunc("/clients", getClients)
	http.HandleFunc("/clients/create", createClientHandler)
	http.HandleFunc("/clients/update/", updateClientHandler)
	http.HandleFunc("/clients/delete/", deleteClientHandler)
	http.HandleFunc("/clients/", getClientHandler)

	http.HandleFunc("/reservations", getReservations)
	http.HandleFunc("/reservations/", getReservationHandler)
	http.HandleFunc("/reservations/update/", updateReservationHandler)
	http.HandleFunc("/reservations/delete/", deleteReservationHandler)
	http.HandleFunc("/reservations/create", createReservationHandler)

	http.HandleFunc("/chambres", getChambres)
	http.HandleFunc("/chambres/", getChambreHandler)
	http.HandleFunc("/chambres/update/", updateChambreHandler)
	http.HandleFunc("/chambres/delete/", deleteChambreHandler)
	http.HandleFunc("/chambres/create", createChambreHandler)

	http.HandleFunc("/hotels", getHotels)
	http.HandleFunc("/hotels/", getHotelHandler)
	http.HandleFunc("/hotels/update/", updateHotelHandler)
	http.HandleFunc("/hotels/delete/", deleteHotelHandler)
	http.HandleFunc("/hotels/create", createHotelHandler)

	http.HandleFunc("/services", getServices)
	http.HandleFunc("/services/", getServiceHandler)
	//http.HandleFunc("/services/update/", updateServiceHandler)
	http.HandleFunc("/services/delete/", deleteServiceHandler)
	http.HandleFunc("/services/create", createServiceHandler)

	// Lancement du serveur sur le port 8080
	log.Fatal(http.ListenAndServe(":8080", nil))

}
