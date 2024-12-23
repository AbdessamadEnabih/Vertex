package net

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"reflect"
	"strings"
	"time"

	"github.com/AbdessamadEnabih/Vertex/internal/persistance"
	"github.com/AbdessamadEnabih/Vertex/pkg/config"
	"github.com/AbdessamadEnabih/Vertex/pkg/state"
)

type Server struct {
	state *state.State
}

func NewServer(state *state.State) *Server {
	return &Server{state: state}
}

func getServerConfiguration() (string, int) {
	server_config, err := config.GetConfigByField("Server")
	if err != nil {
		log.Printf("Error while loading Server configuration: %s", err)
		return "0.0.0.0", 6380
	}
	return reflect.ValueOf(server_config).FieldByName("Adress").String(), int(reflect.ValueOf(server_config).FieldByName("Port").Int())
}

func (s *Server) Start() error {
	address, port := getServerConfiguration()

	fmt.Print(port)
	ln, err := net.Listen("tcp", fmt.Sprintf("%s:%d", address, port))
	if err != nil {
		return err
	}
	defer ln.Close()

	log.Printf("TCP server listening on %s:%d\n", address, port)

	ticker := time.NewTicker(30 * time.Second)

	defer ticker.Stop()

	go func() {
		for range ticker.C {
			persistance.Save(s.state)
			log.Printf("State Saved")
		}
	}()

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		go s.handleConnection(conn)
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				log.Printf("Error reading message: %v\n", err)
			}
			break
		}

		log.Printf("Received message: %s\n", msg)
		msg = strings.TrimSpace(msg)

		switch {
		case strings.HasPrefix(msg, "SET "):
			key, value, err := parseKeyValue(msg)
			if err == nil {
				err = s.state.Set(key, value)
				if err != nil {
					writer.WriteString(formatErrorString(msg, err.Error()))
				} else {
					writer.WriteString("OK\r\n")
				}
			} else {
				writer.WriteString(formatErrorString(msg, err.Error()))
			}

		case strings.HasPrefix(msg, "GET "):
			key := strings.TrimSpace(msg[4:]) // Remove "GET " prefix
			value, err := s.state.Get(key)
			if err != nil {
				writer.WriteString("NODATA\r\n")
			} else {
				writer.WriteString(fmt.Sprintf("VALUE %s\r\n%s\r\n", key, value))
			}

		case strings.HasPrefix(msg, "DELETE "):
			key := strings.TrimSpace(msg[8:]) // Remove "DELETE " prefix
			err := s.state.Delete(key)
			if err != nil {
				writer.WriteString("NODATA\r\n")
			} else {
				writer.WriteString("OK\r\n")
			}

		case msg == "ALL":
			values := s.state.GetAll()
			writer.WriteString("ALLVALUES\r\n")
			for key, value := range values {
				writer.WriteString(fmt.Sprintf("%s=%s\r\n", key, value))
			}
			writer.WriteString("ENDOFALLVALUES\r\n")

		case msg == "FLUSH":
			err := s.state.FlushAll()
			if err != nil {
				writer.WriteString(formatErrorString(msg, err.Error()))
			} else {
				writer.WriteString("OK\r\n")
			}

		default:
			writer.WriteString("UNKNOWN\r\n")
		}

		writer.Flush()
	}
}

func parseKeyValue(msg string) (string, string, error) {
	parts := strings.SplitN(msg, " ", 3)
	if len(parts) != 3 {
		return "", "", errors.New("invalid format")
	}
	return parts[1], parts[2], nil
}

func formatErrorString(command string, err string) string {
	return fmt.Sprintf("Error in %s : %s", command, err)
}
