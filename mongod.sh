#!/bin/bash
mkdir ./mongo
mkdir ./mongo/log
mkdir ./mongo/db

mongod --config mongod.conf
