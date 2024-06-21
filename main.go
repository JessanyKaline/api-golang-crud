package main

//Importa pacotes necessários para a aplicação
import (
	"encoding/json" //Pacote para codificação e decodificação JSON
	"fmt" //Pacote para formatação de entrada e saída
	"log" //Pacote para registro de mensagens de log
	"net/http" //Pacote para criar servidores HTTP
	"os" // Pacote para interagir com o sistema operacional, como leitura de arquivos.
	"strconv" //Pacote para conversões de string para tipos básicos.

	"github.com/gorilla/mux" //Pacote externo para roteamento HTTP
)

// struct Define a estrutura de dados
// tag Json Define como são mapeados para e de json
type Event struct {
	ID           int      `json:"id"`
	Name         string   `json:"name"`
	Organization string   `json:"organization"`
	Date         string   `json:"date"`
	Price        int      `json:"price"`
	Rating       string   `json:"rating"`
	ImageURL     string   `json:"image_url"`
	CreatedAt    string   `json:"created_at"`
	Location     string   `json:"location"`
	Spots        []string `json:"spots"`
	ReservedSpots []string `json:"reserved_spots"`
}

// Variável Global que armazenará a lista de eventos
var events []Event

//Função que carrega dados, em go usa-se 'func'
func loadData() {
	//Le o arquivo json
	data, err := os.ReadFile("data.json")
	//Quando erro não for nulo registra uma msg de erro e encerra
	if err != nil {
		log.Fatalf("Failed to read data file: %s", err)
	}
    //Aqui cria a estrutura usando events e acrescentando spots
	var rawData struct {
		Events []Event `json:"events"`
		Spots  []struct {
			ID      int    `json:"id"`
			Name    string `json:"name"`
			Status  string `json:"status"`
			EventID int    `json:"event_id"`
		} `json:"spots"`
	}
    
	//Deserializa json em estruturas Go do arquivo data.json para estrutura rawData
	err = json.Unmarshal(data, &rawData)
	//Caso de algum erro registra e finaliza
	if err != nil {
		log.Fatalf("Failed to unmarshal data: %s", err)
	}

	// Populando os spots dos eventos

	//Cria um mapa chamado mapeia IDs de eventos para listas de nomes de spots
	eventSpots := make(map[int][]string)
	//Itera sobre todos os spots deserializados
	for _, spot := range rawData.Spots {
		//Adiciona o nome do spot atual à lista de spots do evento
		eventSpots[spot.EventID] = append(eventSpots[spot.EventID], spot.Name)
	}
    
	//Itera sobre todos os eventos deserializados
	for i, event := range rawData.Events {
		//Atribui a lista de spots correspondentes ao evento atual.
		rawData.Events[i].Spots = eventSpots[event.ID]
		//Inicializa a lista de spots reservados como vazia para cada evento.
		rawData.Events[i].ReservedSpots = []string{}
	}
    //Atribui a lista de eventos deserializados e atualizados à variável global events.
	events = rawData.Events
}

//Lista todos os eventos
//http.ResponseWriter - interface usada para construir a resposta http
func getEvents(w http.ResponseWriter, r *http.Request) {
	//Codifica a lista de eventos em JSON e escreve na resposta
	json.NewEncoder(w).Encode(events)
}

//w http.ResponseWriter: Interface para escrever a resposta HTTP
//*r http.Request: Estrutura que representa a requisição HTTP recebida
func getEventByID(w http.ResponseWriter, r *http.Request) {
	//Extrai variaveis da rota
	params := mux.Vars(r)
	//Obtém o valor do parâmetro eventId da URL, que é uma string
	eventID, err := strconv.Atoi(params["eventId"])
	//Trata o erro e envia msg, encerra função
	if err != nil {
		http.Error(w, "Invalid event ID", http.StatusBadRequest)
		return
	}
    //Itera sobre a lista global events
	for _, event := range events {
		//Compara o ID do evento atual com o eventID extraído da URL
		if event.ID == eventID {
			//Se o ID corresponder, codifica o evento encontrado em formato JSON e o escreve na resposta HTTP
			json.NewEncoder(w).Encode(event)
			return
		}
	}

	http.Error(w, "Event not found", http.StatusNotFound)
}

func getEventSpots(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	eventID, err := strconv.Atoi(params["eventId"])
	if err != nil {
		http.Error(w, "Invalid event ID", http.StatusBadRequest)
		return
	}

	for _, event := range events {
		if event.ID == eventID {
			//Aqui é a linha que muda se comparado a função getEventByID
			json.NewEncoder(w).Encode(event.Spots)
			return
		}
	}

	http.Error(w, "Event not found", http.StatusNotFound)
}

func reserveSpot(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	eventID, err := strconv.Atoi(params["eventId"])
	if err != nil {
		http.Error(w, "Invalid event ID", http.StatusBadRequest)
		return
	}
    // Define uma estrutura para armazenar os dados do corpo da requisição
	var requestData struct {
		Spot string `json:"spot"`
	}

    //Decodifica o corpo da requisição JSON e armazena os dados na estrutura requestData.
	err = json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	for i, event := range events {
		if event.ID == eventID {
			// Itera sobre a lista de spots reservados do evento.
			for _, reservedSpot := range event.ReservedSpots {
				// Checa se o spot esta reservado
				if reservedSpot == requestData.Spot {
					http.Error(w, "Spot already reserved", http.StatusBadRequest)
					return
				}
			}

			// Adiciona o spot solicitado à lista de spots reservados do evento
			events[i].ReservedSpots = append(events[i].ReservedSpots, requestData.Spot)
			//Codifica uma mensagem de sucesso em JSON e a escreve na resposta HTTP
			json.NewEncoder(w).Encode(map[string]string{"message": "Spot reserved successfully"})
			return
		}
	}

	http.Error(w, "Event not found", http.StatusNotFound)
}

func main() {
	loadData()
    //Cria um novo roteador HTTP
	r := mux.NewRouter()
	//Define as rotas e associa as funções
	r.HandleFunc("/events", getEvents).Methods("GET")
	r.HandleFunc("/events/{eventId}", getEventByID).Methods("GET")
	r.HandleFunc("/events/{eventId}/spots", getEventSpots).Methods("GET")
	r.HandleFunc("/event/{eventId}/reserve", reserveSpot).Methods("POST")

	fmt.Println("Server running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
