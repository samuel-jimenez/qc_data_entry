#!/usr/bin/fish

sed -i -e '1s/^/PRAGMA foreign_keys = ON;\n/' -e 's/CREATE TABLE IF NOT EXISTS/create table/g' -e 's/INSERT INTO/insert into/' -e 's/"//g' -e s/VALUES/values/ -e s/NULL/null/g qc.0.sqlite3.sql
