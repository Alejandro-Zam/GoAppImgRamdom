package main

import (
	"encoding/base64"
	"fmt"
	"html/template"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type ImageData struct {
	Name string
	Data string
}

func handler(w http.ResponseWriter, r *http.Request) {

	dir := ruta

	hostname, err := os.Hostname()
	if err != nil {
		fmt.Println("Error al obtener el nombre de host:", err)

	}

	rand.Seed(time.Now().UnixNano())

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Selecciona 3 índices aleatorios para las imágenes
	randomIndices := rand.Perm(len(files))[:3]
	var images []ImageData

	// Itera sobre los índices seleccionados y codifica las imágenes en base64
	for _, index := range randomIndices {
		fileName := files[index].Name()
		imagePath := filepath.Join(dir, fileName)

		// Leer los datos de la imagen
		imageBytes, err := ioutil.ReadFile(imagePath)
		if err != nil {
			fmt.Println(err)
			continue
		}

		imageBase64 := base64.StdEncoding.EncodeToString(imageBytes)

		imageInfo := ImageData{

			Name: fileName,
			Data: imageBase64,
		}

		images = append(images, imageInfo)

	}

	data := struct {
		Hostname string
		Imagenes []ImageData
	}{
		Hostname: hostname,
		Imagenes: images,
	}

	// Parsear la plantilla
	tmpl, err := template.ParseFiles("template.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Renderizar la plantilla con los datos
	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

var argumentos = os.Args[1:]

var puerto = argumentos[0]

var ruta = argumentos[1]

func main() {

	http.HandleFunc("/", handler)

	fmt.Println("Servidor corriendo en el puerto " + puerto)

	http.ListenAndServe(":"+puerto, nil)

	if err := http.ListenAndServe(":"+puerto, nil); err != nil {
		fmt.Println("Error al iniciar el servidor:", err)
	}
}
