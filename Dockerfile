FROM postgres:14.8-alpine

ENV POSTGRES_USER=zkp_user
ENV POSTGRES_PASSWORD=zkp_password
ENV POSTGRES_DB=zkp_db

# COPY internal/db/schema.sql /docker-entrypoint-initdb.d/
EXPOSE 5432