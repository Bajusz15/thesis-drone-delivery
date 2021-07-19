FROM library/postgres:9.5.18
ENV POSTGRES_PASSWORD=postgres
ENV POSTGRES_HOST_AUTH_METHOD=trust
ADD init.sql /docker-entrypoint-initdb.d/
