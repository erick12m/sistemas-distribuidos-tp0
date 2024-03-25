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
