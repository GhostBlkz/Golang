package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type User struct {
	ID       int
	Username string
	Email    string
}

type bicicleta struct {
	ID     int
	Modelo string
	Marca  string
	Preco  float64
}

var stringConexao = "user=postgres dbname=golang password=123 host=localhost sslmode=disable"

// album represents data about a record album.
type album struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

// albums slice to seed record album data.
var albums = []album{
	{ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
	{ID: "2", Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
	{ID: "3", Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
}

func main() {
	router := gin.Default()
	// relaciona a rota /albums com a função getAlbums
	router.GET("/albums", getAlbums)
	router.GET("/albums/:id", getAlbumById)
	router.DELETE("/albums/:id", deleteAlbumById)
	router.POST("/albums", postAlbums)

	router.GET("/users", getAllUsers)
	router.GET("/users/:id", getUserById)
	router.POST("/users", addUser)
	router.DELETE("/users/:id", deleteUserById)
	router.PATCH("/users/:id", updateUserById)

	// rotas da bicicleta
	router.GET("/bicicletas", getBicicletas)
	router.POST("/bicicletas", postBicicletas)
	router.GET("/bicicletas/:id", getBicicletaById)
	router.DELETE("/bicicletas/:id", deleteBicicletaById)
	router.PATCH("/bicicletas/:id", updateBicicletaById)

	// sobe o servidor
	router.Run("localhost:8080")
}

func updateUserById(c *gin.Context) {
	id := c.Param("id")
	var userUpdated User
	if err := c.BindJSON(&userUpdated); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid request data"})
		return
	}
	db, err := sql.Open("postgres", stringConexao)
	if err != nil {
		log.Fatal(err)
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Não conectou"})
		return
	}
	defer db.Close()
	query := "update users set username = $1, email = $2 where id = $3"
	result, err := db.Exec(query, userUpdated.Username, userUpdated.Email, id)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "erro servidor"})
		return
	}
	rowsAffected, err := result.RowsAffected()
	if rowsAffected == 0 {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "user not found"})
		return
	}
	c.IndentedJSON(http.StatusOK, userUpdated)
}
func getUserById(c *gin.Context) {
	id := c.Param("id")
	db, err := sql.Open("postgres", stringConexao)
	if err != nil {
		log.Fatal(err)
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Não conectou"})
		return
	}
	defer db.Close()
	var user User
	query := "select id, username, email from users where id = $1"
	err = db.QueryRow(query, id).Scan(&user.ID, &user.Username, &user.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			c.IndentedJSON(http.StatusNotFound, gin.H{"message": "User not found"})
			return
		}
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Erro servidor"})
		return
	}
	c.IndentedJSON(http.StatusOK, user)
}

func deleteUserById(c *gin.Context) {
	id := c.Param("id")
	db, err := sql.Open("postgres", stringConexao)
	if err != nil {
		log.Fatal(err)
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Could not connect"})
		return
	}
	defer db.Close()
	result, err := db.Exec("DELETE FROM users where id = $1", id)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Could not delete"})
		return
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Error checking rows"})
		return
	}
	if rowsAffected == 0 {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "User not found"})
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"messagem": "User removed sucessfully"})
}
func addUser(c *gin.Context) {
	var newUser User
	if err := c.BindJSON(&newUser); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid request data"})
		return
	}
	db, err := sql.Open("postgres", stringConexao)
	if err != nil {
		log.Fatal(err)
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Could not connect"})
		return
	}
	defer db.Close()
	query := "INSERT INTO users (username, email) VALUES ($1, $2) returning id"
	var userID int
	err = db.QueryRow(query, newUser.Username, newUser.Email).Scan(&userID)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Could not insert"})
		return
	}
	newUser.ID = userID
	c.IndentedJSON(http.StatusCreated, newUser)
}

func getAllUsers(c *gin.Context) {
	// abri a conexão com banco de dados
	db, err := sql.Open("postgres", stringConexao)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	rows, err := db.Query("SELECT id, username, email FROM users")
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, err)
	}
	var users []User
	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Username, &user.Email)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
		}
		users = append(users, user)
	}
	c.IndentedJSON(http.StatusOK, users)
}

// getAlbums responds with the list of all albums as JSON.
func getAlbums(c *gin.Context) {
	// retorna ao frontend a lista de album no formato JSON
	c.IndentedJSON(http.StatusOK, albums)
}

func postAlbums(frontend *gin.Context) {
	var vinil album
	// vinil vai guardar os dados do frontend
	if err := frontend.BindJSON(&vinil); err != nil {
		return // tem erro
	}
	// adiciona vinil no album
	albums = append(albums, vinil)
	frontend.IndentedJSON(http.StatusCreated, vinil)
}

func getAlbumById(frontend *gin.Context) {
	id := frontend.Param("id")
	for _, a := range albums {
		if a.ID == id {
			frontend.IndentedJSON(http.StatusOK, a)
			return
		}
	}
	frontend.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
}

func deleteAlbumById(frontend *gin.Context) {
	id := frontend.Param("id")
	// cria um vetor auxiliar
	aux := []album{}

	for _, a := range albums {
		if a.ID != id {
			aux = append(aux, a)
		}
	}
	if len(aux) == len(albums) {
		frontend.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
		return
	}
	albums = aux
	frontend.IndentedJSON(http.StatusOK, gin.H{"message": "album deleted"})
}







// Métodos CRUDS da tabela bicicleta

func getBicicletas(c *gin.Context) {
	db, err := sql.Open("postgres", stringConexao)
	if err != nil {
		log.Fatal(err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Erro ao conectar no banco de dados"})
		return
	}
	defer db.Close()

	rows, err := db.Query("SELECT id, modelo, marca, preco FROM bicicletas")
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Erro ao buscar bicicletas"})
		return
	}

	defer rows.Close()

	var bicicletas []bicicleta
	for rows.Next() {
		var bike bicicleta
		err := rows.Scan(&bike.ID, &bike.Modelo, &bike.Marca, &bike.Preco)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Erro ao ler dados da bicicleta"})
			return
		}
		bicicletas = append(bicicletas, bike)
	}

	if err = rows.Err(); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Erro no resultado das bicicletas"})
		return
	}

	c.IndentedJSON(http.StatusOK, bicicletas)
}

func postBicicletas(frontend *gin.Context) {
	var bike bicicleta
	if err := frontend.BindJSON(&bike); err != nil {
		return
	}
	db, err := sql.Open("postgres", stringConexao)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	query := "INSERT INTO bicicletas (modelo, marca, preco) VALUES ($1, $2, $3) returning id"
	var bikeID int
	err = db.QueryRow(query, bike.Modelo, bike.Marca, bike.Preco).Scan(&bikeID)
	if err != nil {
		frontend.IndentedJSON(http.StatusInternalServerError, err)
		return
	}
	bike.ID = bikeID
	frontend.IndentedJSON(http.StatusCreated, bike)
}

func getBicicletaById(frontend *gin.Context) {
	id := frontend.Param("id")
	db, err := sql.Open("postgres", stringConexao)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	var bike bicicleta
	query := "SELECT id, modelo, marca, preco FROM bicicletas WHERE id = $1"
	err = db.QueryRow(query, id).Scan(&bike.ID, &bike.Modelo, &bike.Marca, &bike.Preco)
	if err != nil {
		if err == sql.ErrNoRows {
			frontend.IndentedJSON(http.StatusNotFound, gin.H{"message": "bicicleta não encontrada"})
			return
		}
		frontend.IndentedJSON(http.StatusInternalServerError, err)
		return
	}
	frontend.IndentedJSON(http.StatusOK, bike)
}

func deleteBicicletaById(frontend *gin.Context) {
	id := frontend.Param("id")
	db, err := sql.Open("postgres", stringConexao)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	result, err := db.Exec("DELETE FROM bicicletas WHERE id = $1", id)
	if err != nil {
		frontend.IndentedJSON(http.StatusInternalServerError, err)
		return
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		frontend.IndentedJSON(http.StatusInternalServerError, err)
		return
	}
	if rowsAffected == 0 {
		frontend.IndentedJSON(http.StatusNotFound, gin.H{"message": "bicicleta não encontrada"})
		return
	}
	frontend.IndentedJSON(http.StatusOK, gin.H{"message": "bicicleta deletada com sucesso"})
}

func updateBicicletaById(frontend *gin.Context) {
	id := frontend.Param("id")
	var bikeUpdated bicicleta
	if err := frontend.BindJSON(&bikeUpdated); err != nil {
		frontend.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid request data"})
		return
	}
	db, err := sql.Open("postgres", stringConexao)
	if err != nil {
		log.Fatal(err)
		frontend.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Could not connect"})
		return
	}
	defer db.Close()
	query := "UPDATE bicicletas SET modelo = $1, marca = $2, preco = $3 WHERE id = $4"
	result, err := db.Exec(query, bikeUpdated.Modelo, bikeUpdated.Marca, bikeUpdated.Preco, id)
	if err != nil {
		frontend.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Could not update"})
		return
	}
	rowsAffected, err := result.RowsAffected()
	if rowsAffected == 0 {
		frontend.IndentedJSON(http.StatusNotFound, gin.H{"message": "bicicleta nao encontrada"})
		return
	}
	frontend.IndentedJSON(http.StatusOK, bikeUpdated)
}
