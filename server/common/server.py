import socket
import logging
import signal
from common.connection_handler import ConnectionHandler
from common.utils import  get_winners_for_agency, store_bets, deserialize_bets


class Server:
    def __init__(self, port, listen_backlog):
        # Initialize server socket
        self._server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self._server_socket.bind(('', port))
        self._server_socket.listen(listen_backlog)
        self._server_shutdown = False
        self._connections_accepted = 0
        self._clients_finished = 0
        
        # Register signal handler for graceful shutdown
        signal.signal(signal.SIGTERM, self.__signal_handler)
     
    def __signal_handler(self, signum, frame):
        """
        Signal handler for graceful shutdown
        """
        if signum == signal.SIGTERM:
            self.graceful_shutdown()

    def graceful_shutdown(self):
        """
        Graceful shutdown of the server
        """
        logging.info("action: graceful_shutdown | result: in_progress")
        self._server_shutdown = True
        self._server_socket.shutdown(socket.SHUT_RDWR)
        self._server_socket.close()
        logging.info("action: socket_shutdown | result: success")
        logging.info("action: graceful_shutdown | result: success")
           
    def run(self):

        while not self._server_shutdown:
            logging.debug(f"connections_accepted: {self._connections_accepted}")
            client_sock = self.__accept_new_connection()
            if client_sock:
                self.__handle_client_connection(client_sock)

    def __handle_client_connection(self, client_sock):
        """
        Read message from a specific client socket and closes the socket

        If a problem arises in the communication with the client, the
        client socket will also be closed
        """
        if self._server_shutdown:
            return
        try:
            while True:
                #Loop with a single client sending bets
                logging.info("action: receive_message | result: in_progress")
                connection_handler = ConnectionHandler(client_sock)
                message = connection_handler.read_message()
                logging.info(f"action: receive_message | result: success | ip: {client_sock.getpeername()[0]}")
                if message == "Finished":
                    self._clients_finished += 1
                    break
                bets = deserialize_bets(message)
                store_bets(bets)
                connection_handler.send_message("Bets stored successfully")
            if self.__can_return_winners():
                logging.info("action: sorteo | result: success")
            logging.info(f"action: store_bets | result: finished | ip: {client_sock.getpeername()[0]}")
            while True:
                #Loop with the client requesting the winners
                logging.info("action: receive_winners_message | result: in_progress")
                connection_handler = ConnectionHandler(client_sock)
                message = connection_handler.read_message()
                logging.info(f"action: receive_winners_message | result: success | ip: {client_sock.getpeername()[0]} | message: {message}")
                if self.__can_return_winners():
                    winners_message = get_winners_for_agency(message)
                    connection_handler.send_message(winners_message)
                    logging.info(f"action: send_winners_message | result: success | ip: {client_sock.getpeername()[0]} | message: {winners_message}")
                    break
                connection_handler.send_message("Error: not all clients finished yet")
                
        except OSError as e:
            connection_handler.send_message("Error storing bet")
            logging.error("action: receive_message | result: fail | error: {e}")
        finally:
            client_sock.close()
            
    def __can_return_winners(self):
        return self._clients_finished == self._connections_accepted

    def __accept_new_connection(self):
        """
        Accept new connections

        Function blocks until a connection to a client is made.
        Then connection created is printed and returned
        """

        # Connection arrived
        if self._server_shutdown:
            return None
        try:
            logging.info('action: accept_connections | result: in_progress')
            c, addr = self._server_socket.accept()
            logging.info(f'action: accept_connections | result: success | ip: {addr[0]}')
            self._connections_accepted += 1
            return c
        except OSError as e:
            logging.error(f'action: accept_connections | result: fail | error: {e}')
            return None
