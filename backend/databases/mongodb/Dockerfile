FROM mongo:4.4.5
ENV MONGO_INITDB_DATABASE=drone_delivery
ENV MONGO_INITDB_ROOT_USERNAME=admin-user
ENV MONGO_INITDB_ROOT_PASSWORD=admin-password
ADD init-mongo.js /docker-entrypoint-initdb.d/
