import csv
from py2neo import neo4j, rel
import sys
from time import sleep
from database import neo4j_db as database
import logging
import re

logging.basicConfig(level=logging.DEBUG)

class uploader():
    """
    Function for uploading crawl data into Neo4j
    """
    def __init__(self, dbpath, report_date):
        self.report_date = report_date
        self.batch_limit = 5000

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
        batch.clear()
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
        first_name = kwargs[6]
        if len(kwargs[13]) > 1 and ' ' not in kwargs[13]:
            first_name = kwargs[13]
        name = " ".join([x for x in [first_name, kwargs[7], kwargs[8]] if len(x) > 1]).title()
        self.batch_cypher(self.write_batch, ('CREATE (person:Person {\
            name: "%s",\
             f_name: "%s",\
             m_name: "%s",\
             l_name: "%s",\
             pref_name: "%s",\
             title: "%s",\
             dob: "%s",\
             gender: "%s",\
             email: "%s",\
             department: "", \
             position: "%s",\
             position_group: "%s",\
             department_number: "%s", \
             loc_building: "%s",\
             loc_campus: "%s",\
             loc_floor: "%s",\
             loc_room: "%s",\
             phone: "%s", \
             mobile: "%s",\
             employee_number: "%s",\
             billing_code: "%s",\
             start_date: "%s"})' % (name, kwargs[6].title(), kwargs[7].title(), kwargs[8].title(), kwargs[13].title(), kwargs[4].title(), kwargs[5], kwargs[9], kwargs[3], kwargs[21].title(), kwargs[24], 
             kwargs[18], kwargs[34], kwargs[35], kwargs[36], kwargs[37], kwargs[31], kwargs[33], kwargs[1], kwargs[38], kwargs[11])))

    def add_fa(self, fa):
        '''Insert new record into db'''
        self.batch_cypher(self.write_batch, ('CREATE (area:Area {name: "%s"})' % (fa)))

    def add_reportingline(self, p_from, p_to):
        '''Insert reporting line relationship'''
        self.batch_cypher(self.write_batch, ('MATCH (a:Person), (b:Person) WHERE a.employee_number="%s" and b.employee_number="%s" CREATE (a)-[:MANAGES]->(b)' % (p_from, p_to)))

    def add_index(self, value):
        '''Create index'''
        self.batch_cypher(self.write_batch, ('CREATE INDEX ON :Person(%s)' % (value)))
        self.batch_flush(self.write_batch)

    def add_depcode(self, depcode, depname):
        '''Update all nodes with the given depcode, and set department name'''
        self.batch_cypher(self.write_batch, ('MATCH (a:Person) WHERE a.department_number="%s" SET a.department="%s"' % (depcode, depname)))

    def update_phone(self, email, phone, mobile=None):
        '''Update record with new phone and mobile numbers'''
        self.batch_cypher(self.write_batch, ('MATCH (a:Person) WHERE a.email="%s" SET a.phone="%s"' % (email, phone)))

if __name__ == '__main__':

    dbpath = 'http://localhost:7474/db/data'
    report_date = '2014-05-22'
    loader = uploader(dbpath, report_date)

    def feed_data():
        the_list = {}
        temp = []
        with open('may_2014.csv', newline='', encoding="ISO-8859-1") as csvfile:
            data = csv.reader(csvfile)
            try:
                for row in data:
                    print(row)
                    if row[1] in the_list:
                        temp = the_list[row[1]][0]
                    the_list[row[1]] = row
                    if temp:
                        the_list[row[1]].append(temp)
                    temp = []

                    if row[19] not in the_list:
                        the_list[row[19]] = [[],]
                    if type(the_list[row[19]][-1]) == str:
                        the_list[row[19]].append([])
                    the_list[row[19]][-1].append(row[1])
            except UnicodeDecodeError as err:
                print("UnicodeDecodeError {0}".format(err)) 
                sys.exit(1)

        add_people(the_list)
        add_relationship(the_list)
        add_departments()

    def update_phones():
        depcodes = []
        pattern = re.compile('[\W_]+')
        with open('activedir_2013.csv', newline='') as csvfile:
                data = csv.reader(csvfile)
                try:
                    for row in data:
                        if row[0] and row[1] and 'E' not in row[1] and  row[1][0] != '1': 
                            row[1] = pattern.sub('', row[1])
                            row[2] = pattern.sub('', row[2])
                            if '61' in row[1]:
                                print(row[1])
                                row[1] = row[1][2:]
                            if len(row[1]) < 9:
                                print(row)
                                depcodes.append(row)
                except UnicodeDecodeError as err:
                    print("UnicodeDecodeError {0}".format(err))
                    sys.exit(1)
        for each in depcodes:
            if len(each)> 1:
                loader.update_phone(*each)
        loader.batch_flush(loader.write_batch)
    
    def add_people(the_list):
        for person in the_list:
            if len(person)> 1 and len(the_list[person]) > 1:
                loader.add_person(the_list[person])
        loader.batch_flush(loader.write_batch)

    def add_relationship(the_list):
        loader.add_index('employee_number')
        for person in the_list:
            if type(the_list[person][-1]) is list and len(the_list[person]) > 1:
                for each in the_list[person][-1]:
                    loader.add_reportingline(the_list[person][1], each)
        loader.batch_flush(loader.write_batch)

    def add_areas():
        for fa in set([x[5] for x in temp]):
            if len(fa) > 1:
                loader.add_fa(fa)
        loader.batch_flush(loader.write_batch)

    def add_departments():
        depcodes = []
        with open('department_codes.csv', newline='') as csvfile:
                data = csv.reader(csvfile)
                try:
                    for row in data:
                       depcodes.append(row)
                except UnicodeDecodeError as err:
                    print("UnicodeDecodeError {0}".format(err))
                    sys.exit(1)
        for each in depcodes:
            loader.add_depcode(each[0], each[1])
        loader.batch_flush(loader.write_batch)


    feed_data()
    update_phones()


