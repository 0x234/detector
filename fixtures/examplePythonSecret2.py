#!/usr/bin/python
import MySQLdb

# Deliberately insecure Python source code to verify the detection
# rule is working. This time in a class and split over multiple lines.

class DBConnection:

    def database(self):
        self.db = MySQLdb.connect('db.corp.com',
                                    'admin',
                                    'hunter2','production')

    def display_data(self):
        cur = self.db.cursor()
        cur.execute('SELECT * FROM SOMETHING_SECURE')
        for row in cur.fetchall():
            print(row[0])

        self.db.close()