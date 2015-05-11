#!/bin/bash
cd ~
# superuser needed to create hstore extension on the database
psql -c "CREATE ROLE dbworkbench_admin WITH LOGIN SUPERUSER PASSWORD 'f1r34nd1c3';" postgres

psql -c "CREATE ROLE dbworkbench_demo WITH LOGIN PASSWORD 's0ng4ndd4nc3';" postgres
psql -c "CREATE DATABASE dbworkbench_demo OWNER dbworkbench_demo;" postgres

touch db_inited
