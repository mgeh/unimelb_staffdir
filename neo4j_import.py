import csv
from py2neo import neo4j, rel
import sys
from time import sleep
from database import neo4j_db as database
import logging

logging.basicConfig(level=logging.DEBUG)

class uploader():
    """
    Function for uploading crawl data into Neo4j
    """
    def __init__(self, dbpath, data, report_date):
        self.report_date = report_date
        self.data = data
        self.hosts = []
        self.pages = []
        self.batch_limit = 1000

        self.db = database(dbpath)
        self.db.connect()
        logging.debug("init: database connection completed")

        self.write_batch = self.db.write_batch

        self.db = database(dbpath)
        self.db.connect()
        logging.debug("init: database connection completed")

    def batch_common(self, batch):
        if len(batch) == self.batch_limit:
            try:
                batch.submit()
            except:
                sleep(1)
                self.db.connect()
                sleep(1)
                batch.submit()
            logging.debug("batch_common: flushing batch")
            batch.clear()
        return

    def batch_flush(self, batch):
        output = batch.submit()
        logging.debug("batch_flush: flushing batch")
        return output

    def batch_add(self, data_type, batch, index, key, value, item):
        '''Create a new node in index'''
        batch.get_or_create_in_index(data_type, index, key, value, item)
        return self.batch_common(batch)

    def batch_cypher(self, batch, request):
        batch.append_cypher(request)
        return self.batch_common(batch)

    def add_person(self, kwargs):
        '''Insert new record into db'''
        self.batch_cypher(self.write_batch, ('CREATE (person:People {name: "%s", email: "%s", department: "%s", department_number: "%s", phone: "%s", mobile: "%s"})' % (kwargs[0],kwargs[4], kwargs[5], kwargs[6], kwargs[9], kwargs[10])))

    def add_fa(self, fa):
        '''Insert new record into db'''
        self.batch_cypher(self.write_batch, ('CREATE (area:Areas {name: "%s"})' % (fa)))



if __name__ == '__main__':

    dbpath = 'http://localhost:7474/db/data'
    report_date = '2013-05-17'

    temp = []
    with open('contacts.csv', newline='') as csvfile:
        data = csv.reader(csvfile)
        for row in data:
            temp.append(row)
    temp = temp[1:]

    loader = uploader(dbpath, data, report_date)
    
    def add_people():
        for person in temp:
            loader.add_person(person)
        loader.batch_flush(loader.write_batch)

    def add_areas():
        for fa in set([x[5] for x in temp]):
            if len(fa) > 1:
                loader.add_fa(fa)
        loader.batch_flush(loader.write_batch)

    # add_people()


