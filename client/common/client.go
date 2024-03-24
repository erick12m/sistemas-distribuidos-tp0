package common

import (
	"bufio"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"
)

// ClientConfig Configuration used by the client
type ClientConfig struct {
	ID            string
	ServerAddress string
	LoopLapse     time.Duration
	LoopPeriod    time.Duration
	BatchSize     int
	BetsFile      string
}

// Client Entity that encapsulates how
type Client struct {
	config    ClientConfig
	conn      ConnectionHandler
	data      ClientData
	completed bool
}

type ClientData struct {
	Name      string
	LastName  string
	Document  string
	Birthdate string
	Number    string
}

// NewClient Initializes a new client receiving the configuration and data
// as a parameter
func NewClient(config ClientConfig, data ClientData) *Client {
	client := &Client{
		config:    config,
		data:      data,
		completed: false,
	}
	return client
}

// CreateClientSocket Initializes client socket. In case of
// failure, error is printed in stdout/stderr and exit 1
// is returned
func (c *Client) createClientSocket() error {
	conn, err := net.Dial("tcp", c.config.ServerAddress)
	if err != nil {
		log.Fatalf(
			"action: connect | result: fail | client_id: %v | error: %v",
			c.config.ID,
			err,
		)
	}
	stream := initStream(conn, c.config.ID)
	c.conn = *InitConnectionHandler(stream, c.config.ID)
	return nil
}

// StartClientLoop Send messages to the client until some time threshold is met
func (c *Client) StartClientLoop() {

	// Set up a signal handler to gracefully shutdown the client
	signalHandlerChannel := make(chan os.Signal, 1)
	signal.Notify(signalHandlerChannel, syscall.SIGTERM)

	select {
	case <-signalHandlerChannel:
		c.handleShutdown(signalHandlerChannel)
	default:
	}

	// Create the connection the server in every loop iteration. Send an
	c.createClientSocket()
	betsFile, err := openBetFile(c.config.BetsFile)
	if err != nil {
		log.Fatalf("action: open_file | result: fail | client_id: %v | error: %v",
			c.config.ID,
			err,
		)
	}
	defer betsFile.Close()

	scanner := bufio.NewScanner(betsFile)

	for !c.completed {
		batch := ""
		for i := 0; i < c.config.BatchSize; i++ {
			if !scanner.Scan() {
				c.completed = true
				break
			}
			batch += c.config.ID + "," + scanner.Text() + "\n"
		}
		if err != nil {
			log.Fatalf("action: get_bets_batch | result: fail | client_id: %v | error: %v",
				c.config.ID,
				err,
			)
		}
		c.conn.Write(batch)

		response, err := c.conn.Read()
		if err != nil {
			log.Errorf("action: read | result: fail | client_id: %v | error: %v",
				c.config.ID,
				err,
			)
		}
		log.Infof("action: read | result: success | client_id: %v | response: %v", c.config.ID, response)
	}
	c.conn.Write("Finished")
	// Send messages if the loopLapse threshold has not been surpassed
	askForWinners(c)
}

func askForWinners(c *Client) {
	log.Infof("action: ask_for_winners | result: in_progress | client_id: %v",
		c.config.ID,
	)
loop:
	for timeout := time.After(c.config.LoopLapse); ; {
		select {
		case <-timeout:
			log.Infof("action: timeout_detected | result: success | client_id: %v",
				c.config.ID,
			)
			break loop
		default:
		}
		c.conn.Write(c.config.ID)
		response, err := c.conn.Read()
		if err != nil {
			log.Errorf("action: read | result: fail | client_id: %v | error: %v",
				c.config.ID,
				err,
			)
		}
		log.Infof("action: read | result: success | client_id: %v | response: %v", c.config.ID, response)
		if response == "Error: not all clients finished yet" {
			time.Sleep(c.config.LoopPeriod)
		} else {
			winners := getWinnersQuantity(response)
			log.Infof("action: consulta_ganadores | result: success | cant_ganadores: %v",
				winners,
			)
			break loop
		}
	}
	c.conn.Close()
}

func (c *Client) handleShutdown(signalHandlerChannel chan os.Signal) {
	log.Infof("action: graceful_shutdown | result: in_progress | client_id: %v",
		c.config.ID,
	)
	c.completed = true
	c.conn.Close()
	log.Infof("action: socket_shutdown | result: success | client_id: %v",
		c.config.ID,
	)
	close(signalHandlerChannel)
	log.Infof("action: signal_handler_channel_shutdown | result: success | client_id: %v",
		c.config.ID,
	)
	log.Infof("action: graceful_shutdown | result: success | client_id: %v",
		c.config.ID,
	)
}
