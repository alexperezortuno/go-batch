# Stage 1: Builder para preparar assets
FROM postgres:16-alpine AS builder

ARG DB_USER
ARG DB_PASS
ENV PGOPTIONS="--client-min-messages=warning"

WORKDIR /sql-scripts

# Copia todos los scripts SQL (asegúrate de que estén en el mismo directorio que el Dockerfile)
COPY 03_init.sql 04_schemas.sql 05_tables.sql ./

# Script maestro que ejecuta los archivos en orden
RUN echo "CREATE USER ${DB_USER} WITH PASSWORD '${DB_PASS}';" >> 01_create-user.sql && \
    echo "CREATE DATABASE app_db OWNER ${DB_USER};" >> 02_create-db.sql && \
    echo "--- initialize db\n\\i 01-create-user.sql\n\\i 02-create-db.sql\n\\c app_db\n\\i 03_init.sql" >> 00_init-all.sql

# Stage 2: Imagen final
FROM postgres:16-alpine

ARG DB_USER
ARG DB_PASS
ENV PGOPTIONS="--client-min-messages=warning"

# Copia los scripts de inicialización
COPY --from=builder /sql-scripts/ /docker-entrypoint-initdb.d/

# Salud del contenedor
HEALTHCHECK --interval=30s --timeout=5s \
  CMD pg_isready -U $DB_USER -d $DB_PASS

VOLUME /var/lib/postgresql/data
EXPOSE 5432
