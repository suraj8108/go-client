package client

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"sync"

	"github.com/suraj8108/clientApp/model"
	"github.com/suraj8108/clientApp/redisutils"
	"github.com/suraj8108/clientApp/utils"
)


type ConnectionDetails struct {
	HttpClient   net.Conn
	Active       bool
	NumberOfExec int
	Port         int
}

type ConnHandler struct {
	HttpConnection chan ConnectionDetails
	StudentDetails chan model.Student
	redisHandler   redisutils.RedisClient
}

func creteHTTPClient(port int) (net.Conn, error) {
	conn, err := net.Dial("tcp", fmt.Sprintf("127.0.0.1:%v", port))
	if err != nil {
		log.Printf("Can not create a TCP connection with Port %v \n", port)
		return nil, err
	}

	return conn, nil
}

func NewConnHandler(redisHandler redisutils.RedisClient) ConnHandler {
	allConnections := make(chan ConnectionDetails, 10)
	port := 8081
	for i := 0; i < 10; i++ {
		conn, err := creteHTTPClient(port)
		if err != nil {
			log.Printf("Can not create a TCP connection with Port %v \n", port)
			continue
		}

		connectionDetails := ConnectionDetails{HttpClient: conn, Active: false, NumberOfExec: 0, Port: port}
		allConnections <- connectionDetails

		port++
	}

	return ConnHandler{
		HttpConnection: allConnections,
		StudentDetails: make(chan model.Student, 1000),
		redisHandler:   redisHandler,
	}
}

func (ch *ConnHandler) ClientOperation() {

	sc := sync.WaitGroup{}

	// Reading TPS files
	txns := utils.ReadTPSFile("tps.txt")

	// My Redis key
	startKey := 0

	sc.Add(1)
	go ch.BuildTCPConnection(&sc)

	for index, txn := range txns {
		fmt.Printf("Operation for %v sec started", index)
		redisKeys := utils.CreateRedisKeys(startKey, txn)
		startKey += txn
		rawData, _ := ch.redisHandler.FetchBulkDataFromRedis(redisKeys)
		studentDetails := utils.UnMarshalRedisBulkData(rawData)
		for _, student := range studentDetails {
			ch.StudentDetails <- student
		}
	}

	sc.Wait()
	close(ch.StudentDetails)
}

func (ch *ConnHandler) BuildTCPConnection(sc *sync.WaitGroup) {
	fmt.Println("Start the TCP calls")

	for students := range ch.StudentDetails {
		fmt.Println("Hit the HTTP Connection")
		fmt.Println("Process Student", students)

		conection := <-ch.HttpConnection
		sc.Add(1)
		go ch.SendHTTPRequest(conection, students, sc)
	}
	fmt.Println("All Student Completed")
	sc.Done()
}

func (ch *ConnHandler) SendHTTPRequest(connection ConnectionDetails, student model.Student, sc *sync.WaitGroup) {
	defer sc.Done()

	connection.Active = true

	defer func() {
		connection.Active = false
	}()

	log.Printf("User %v is handled by Server running on port %v \n", student.StudentName, connection.Port)

	// Increment the execution count
	connection.NumberOfExec = connection.NumberOfExec + 1

	requestBody, _ := json.Marshal(student)
	connection.HttpClient.Write(requestBody)
	buf := make([]byte, 1024)
	n, err := connection.HttpClient.Read(buf)

	if err != nil {
		log.Fatalf("Error while sending the request %v for user %v \n", err, student.StudentName)
	}
	var studentResp model.Student
	json.Unmarshal(buf[:n], &studentResp)

	fmt.Printf("Student Server Response %v for student %v \n", studentResp, student.StudentName)

	if connection.NumberOfExec <= 100 {
		ch.HttpConnection <- connection
	} else {
		port := connection.Port
		conn, err := creteHTTPClient(port)
		if err != nil {
			log.Fatalf("Error while creating HTTP client on Port %v \n", port)
		} else {
			connectionDetails := ConnectionDetails{HttpClient: conn, Active: true, NumberOfExec: 0, Port: port}
			ch.HttpConnection <- connectionDetails
		}
		connection.HttpClient.Close()
	}

	// Goroutine to save the userId
	sc.Add(1)
	go ch.redisHandler.InsertDataInRedisBySelfKey(studentResp.StudentRedisKey, studentResp, sc)
}
