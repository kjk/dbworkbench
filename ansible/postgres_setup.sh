#!/bin/bash
cd ~
psql -c "CREATE ROLE dbworkbench_admin WITH LOGIN PASSWORD 'f1r34nd1c3';" postgres
psql -c "CREATE DATABASE dbworkbench_admin OWNER dbworkbench_admin;" postgres

psql -c "CREATE ROLE dbworkbench_demo WITH LOGIN PASSWORD 's0ng4ndd4nc3';" postgres
psql -c "CREATE DATABASE dbworkbench_demo OWNER dbworkbench_demo;" postgres

touch db_inited
