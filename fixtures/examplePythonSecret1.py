#!/usr/bin/python
import MySQLdb

# Deliberately insecure Python source code to verify the detection
# rule is working.

db = MySQLdb.connect('db.company.com','admin','test123','production')

cur = db.cursor()
cur.execute('SELECT * FROM SOMETHING_SECURE')
for row in cur.fetchall():
    print(row[0])

db.close()