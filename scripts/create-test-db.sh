#!/bin/bash

# Script to create an isolated MariaDB test database for yFramework integration tests

TIMESTAMP=$(date +"%Y%m%d_%H%M%S")
TEST_DB_NAME="devframework_test_${TIMESTAMP}"

echo "Creating test database: ${TEST_DB_NAME}"

# Create the database
mariadb -h 127.0.0.1 -u devframework -pMJeltULtaPx2agep -e "CREATE DATABASE IF NOT EXISTS ${TEST_DB_NAME};"

if [ $? -eq 0 ]; then
    echo "Successfully created test database: ${TEST_DB_NAME}"
    echo "Use this database for integration tests"
    echo "To drop this database after testing, run:"
    echo "mariadb -h 127.0.0.1 -u devframework -pMJeltULtaPx2agep -e \"DROP DATABASE IF EXISTS ${TEST_DB_NAME};\""
else
    echo "Failed to create test database: ${TEST_DB_NAME}"
    exit 1
fi