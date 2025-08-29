# Create a local Postgres 16 container, persist data, expose 5432
docker run -d \
  --name pg-modjot-local \
  -e POSTGRES_USER=appuser \
  -e POSTGRES_PASSWORD=apppass \
  -e POSTGRES_DB=modjot \
  -p 5433:5432 \
  -v pgdata:/var/lib/postgresql/data \
  postgres:16
