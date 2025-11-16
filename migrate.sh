#!/bin/sh
set -e

echo "Waiting for Postgres..."
until pg_isready -h db -p 5432; do
  sleep 1
done

echo "Applying migrations..."
migrate -path ./migrations -database "$DATABASE_URL" up || true

echo "Migrations applied"
