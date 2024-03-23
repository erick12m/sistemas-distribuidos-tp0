# Parte 1: Introducción a Docker

### Ejercicio N°1:

Se modificó el archivo docker-compose-dev.yaml y se agregò un segundo cliente llamado client2 utilizando la misma estructura que tenia client1, cambiando su id

Para comprobar basta con ejecutar nuevamente

```bash
make docker-compose-up
make docker-compose-logs
```

Y se podrá observar en los logs el nuevo cliente y sus mensajes.

### Ejercicio N°1.1:

Se creó un nuevo archivo docker-compose-script-with-n-clients.sh para poder generar dinámicamente clientes basado en comandos por terminal; además se modificó el comando del makefile `docker-compose-up` para utilizar ahora dicho script para inicializar los contenedores y el comando `docker-compose-logs` para usar el nuevo docker compose.

Para usar y comprobar

```bash
make docker-compose-up <Number of clients>
make docker-compose-logs
```

### Ejercicio N°2:

Para este ejercicio se definieron volumenes en el server y client; donde se montan sus respectivos archivos de configuración. Para esto se removió del dockerfile el COPY del archivo de config y se agregaron al .dockerignore dichos archisvos para que no se detecten como cambios al hacer `docker compose up`.

Se agregó al script generador del docker compose el volumen con el archivo de configuración que correspondía según el caso.

Para probar se cambió el puerto del servidor en ambos archivos al 8080 y se observó que los cambios se reflejaron correctamente.

```bash
server | 2024-03-16 02:52:12 DEBUG action: config | result: success | port: 8080 | listen_backlog: 5 | logging_level: DEBUG
server | 2024-03-16 02:52:12 INFO action: accept_connections | result: in_progress
server | 2024-03-16 02:52:13 INFO action: accept_connections | result: success | ip: 172.25.125.4
server | 2024-03-16 02:52:13 INFO action: receive_message | result: success | ip: 172.25.125.4 | msg: [CLIENT 2] Message N°1
server | 2024-03-16 02:52:13 INFO action: accept_connections | result: in_progress
server | 2024-03-16 02:52:13 INFO action: accept_connections | result: success | ip: 172.25.125.3
server | 2024-03-16 02:52:13 INFO action: receive_message | result: success | ip: 172.25.125.3 | msg: [CLIENT 1] Message N°1
server | 2024-03-16 02:52:13 INFO action: accept_connections | result: in_progress
client2 | time="2024-03-16 02:52:13" level=info msg="action: config | result: success | client_id: 2 | server_address: server:8080 | loop_lapse: 20s | loop_period: 5s | log_level: DEBUG"
```

Se comprobó que utilizó la cache al hacer `make docker-compose-down` seguido de `make docker-compose-up clients=2` queriendo decir que no se montó de nuevo toda la imagen por el cambio de configuración.

### Ejercicio N°3:

Se creó un nuevo contenedor `netcat-client` con su propio dockerfile en el que se instala netcat y se utiliza un script para comprobar que la respuesta al enviarle un mensaje al server sea la misma que dicho mensaje.

La ip y el puerto del serrver se pueden configurar en `netcat-client/config.env`

Para realizar el test basta con ejecutar `make run-netcat-test` y observar los resultados. En caso de éxito se verá por terminal

```bash
docker compose -f docker-compose-n-clients.yaml up -d --build --remove-orphans
[+] Running 1/0
 ✔ Container server  Running                                                                                                                                                         0.0s
make[1]: se sale del directorio '/home/erick/distribuidos/sistemas-distribuidos-tp0'
docker run --rm --network sistemas-distribuidos-tp0_testing_net --env-file ./netcat-client/config.env --name netcat-client netcat-client:latest
OK: Server response is correct: message: Test message to server, response: Test message to server
```

En caso de error; por ejemplo al cambiar el puerto del server a uno invalido:

```bash
docker compose -f docker-compose-n-clients.yaml up -d --build --remove-orphans
[+] Running 1/0
 ✔ Container server  Running                                                                                                                                                         0.0s
make[1]: se sale del directorio '/home/erick/distribuidos/sistemas-distribuidos-tp0'
docker run --rm --network sistemas-distribuidos-tp0_testing_net --env-file ./netcat-client/config.env --name netcat-client netcat-client:latest
ERROR: Server response is incorrect: message: Test message to server, response:
```

### Ejercicio N°4:

Se modificó `client/common/client.go` agregando un channel para manejar señales, haciendo que cuando llegue una SIGTERM se cierre el socket; dicho channel y termine el cliente de forma gracefull

Para el server se modificó `server/common/server.py` para agregar un handler de señales; de igual forma que en el cliente se encarga de cerrar primero el socket; ademàs de agregar modificaciones para el momento de aceptar conexiones que no acepte una nueva si se recibe el SIGTERM.

Además se cambió el tiempo de timeout utilizado `make docker-compose-down` a 10 segundos, para que se pueda apreciar en el cliente el cierre gracefull, ya que este contaba con un sleep de 5segundos por iteración.

Se puede comprobar analizando los logs al momento de ejecutar `make docker-compose-down`

```bash
client2  | time="2024-03-16 22:34:52" level=info msg="action: graceful_shutdown | result: in_progress | client_id: 2"
client2  | time="2024-03-16 22:34:52" level=info msg="action: socket_shutdown | result: success | client_id: 2"
client2  | time="2024-03-16 22:34:52" level=info msg="action: signal_handler_channel_shutdown | result: success | client_id: 2"
client2  | time="2024-03-16 22:34:52" level=info msg="action: graceful_shutdown | result: success | client_id: 2"
client2  | time="2024-03-16 22:34:52" level=info msg="action: loop_finished | result: success | client_id: 2"
client1  | time="2024-03-16 22:34:52" level=info msg="action: graceful_shutdown | result: in_progress | client_id: 1"
client1  | time="2024-03-16 22:34:52" level=info msg="action: socket_shutdown | result: success | client_id: 1"
client1  | time="2024-03-16 22:34:52" level=info msg="action: signal_handler_channel_shutdown | result: success | client_id: 1"
client1  | time="2024-03-16 22:34:52" level=info msg="action: graceful_shutdown | result: success | client_id: 1"
client1  | time="2024-03-16 22:34:52" level=info msg="action: loop_finished | result: success | client_id: 1"
client2 exited with code 0
client2 exited with code 0
client1 exited with code 0
client1 exited with code 0
server   | 2024-03-16 22:34:52 INFO     action: graceful_shutdown | result: in_progress
server   | 2024-03-16 22:34:52 INFO     action: socket_shutdown | result: success
server   | 2024-03-16 22:34:52 INFO     action: graceful_shutdown | result: success
server   | 2024-03-16 22:34:52 ERROR    action: accept_connections | result: fail | error: [Errno 22] Invalid argument
server exited with code 0
```

### Ejercicio N°5:

Para este ejercicio se cambió en general tanto cliente como servidor, ambos con la misma estructura; tienen un stream encargado de realizar las escrituras y lecturas del socket para evitar short-read y short-write; una estructura Connection Handler encargada de usar el stream y servir de puente entre el servidor/cliente y el stream, esta estructura se encarga de el armado de los mensajes y su protocolo, por ejemplo se encarga de enviar la longitod del mensaje antes de la data en si; asì el otro extremo sabe cuántos bytes leer.

El protocolo consiste en el envío de mensajes usando TCP, informando primero en 4 bytes la longitud de la data a enviar seguido de la data en sí; que debe cumplir con la longitud informada.

Ademàs se modificó el script generador de docker compose para agregar las nuevas variables de entorno y las apuestas enviadas son separadas por coma; en este momento los clientes envian de a una apuesta, pero el server ya puede deserializar de a mas de una si se separan con un salto de línea.

Para la ejecución basta con correr `make docker-compose-up clients=5` y observer en los logs los mensajes de apuesta_almacenada del server.

### Ejercicio N°6:

Se modificó la lógica del cliente, para ahora usar el archivo de apuestas correspondiente, el cual se va leyendo de a batches configurables con la nueva variable de entorno agregada, y se envian dichos batches al server, cuando se obtiene una respuesta del server se procede a enviar el siguiente batch, así hasta que se completa de leer el archivo, y se envía un mensaje Finished, que cuando el server lo lee da por finalizada la conexión

Para la ejecución basta con correr `make docker-compose-up clients=5` Se puede observer en los archivos del container del server como se escriben las 78697 de los diferentes clientes en el archivo bets.csv. Además de los logs.
