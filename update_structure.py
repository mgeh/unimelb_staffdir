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

    def add_management_line(self, manager, person):
        '''Add management line'''
        self.batch_cypher(self.write_batch, ('MATCH (a:People), (b:People) WHERE a.name="%s" and b.name="%s" SET b.position="%s" CREATE (a)-[:MANAGES]->(b)' % (manager[0], person[0], person[1])))



if __name__ == '__main__':

    dbpath = 'http://weapon.its.unimelb.edu.au:7474/db/data'
    report_date = '2013-05-17'

    data = [[
                ['Steven Manos','Research Services Manager'],
                [['Neil Kileen', 'eResearch Analyst (Neuroscience)'],
                ['David Flanders','Community Manager'],
                ["Owen O'Neil",'Data Archiving Analyst'],
                ['Andy Tseng','Data Infrastructure Architect'],
                ['Dirk Van Der Knijff','HPC Specialist'],
                ['Terry Brennan','Senior Project Manager']]
            ],

            [
                ['David Flanders','Community Manager'],
                [['Jeff Tyler', 'Communications and Engagement Officer'],
                 ['Bernard Meade', 'Outreach & Innovation Officer']
                ]
            ],

            [
                ['Terry Brennan','Senior Project Manager'],
                [['Martin Krahnert','Senior Business Analyst'],
                ['Ursula Soulsby','Senior Business Analyst'],
                ['Greg Sauter','Project Manager'],
                ['Lennie Au', 'Project Manager'],
                ['Nick Golovachenko', 'Project Manager'],
                ['Sara Ogston', 'Senior Project Coordinator']]
            ],

            [
                ['Greg Sauter','Project Manager'],
                [['Clinton Walsh','Cloud User Support Lead'],
                ['Phil Hoenig', 'Technical Writer'],
                ['Marcus Furlong', 'Cloud Development and Operations Engineer'],
                ['Sam Morrison', 'Cloud Development and Operations Engineer']]
            ],

            [
                ['Lennie Au', 'Project Manager'],
                [['Devendran Jagadisan', 'Cloud Development and Operations Engineer'],
                ['Craig Sanders', 'Cloud Development and Operations Engineer'],
                ['Damien Mannix', 'NSP Systems Administrator']]
            ],

            [
                ['Michael Carolan', 'Application Services Director'],
                [['Phil Brown', 'Application Services Manager (Engagement)'],
                ['Priya Ravindra','Application Services Manager (Library, Learning & Teaching)'],
                ['Marlena Axel','Application Services Manager (Admin & Enablement)'],
                ['Simon Abbott','Application Services Manager (Melbourne Students & Learning)'],
                ['Sophia Lagastes','Senior Project Manager'],
                ['Brigitte Whitley','Senior Business Analyst']]
            ],

            [
                ['Phil Brown', 'Application Services Manager (Engagement)'],
                [['Ben Edward', 'Application Services Manager'],
                ['Freddy Navas', 'Business Analyst'],
                ['Koula Tsiaplias', 'Associate Project Manager'],
                ['Lene Cortsen', 'Business Program Lead']]
            ],

            [
                ['Lene Cortsen', 'Business Program Lead'],
                [['Jonathan Wright', 'Project Officer'],
                ['Supriya Gurusprasad','Data Analyst']]
            ],

            [
                ['Sophia Lagastes','Senior Project Manager'],
                [['Paul Beaumont', 'Project Manager'],
                ['Michael Gehling', 'Project Manager'],
                ['Kristy Cross', 'Associate Project Manager'],
                ['Tamara Wishart', 'Associate Project Manager'],
                ['Tony Goodhram', 'Senior Business Analyst'],
                ['Michelle Ely', 'Business Analyst'],
                ['Christine Priestley', 'Business Analyst'],
                ['Kate Taylor', 'Business Analyst'],
                ['Tania Finn Angelo', 'Senior Change & Communications Manager'],
                ['Hilary Sissons', 'Change & Communications Manager']]
            ],

            [
                ['Mark Brodsky', 'IT Strategy & Planning Director'],
                [['Yope Vagenas', 'Organisational Development Manager'],
                ['Wayne Tufek', 'IT Security and Risk Manager'],
                ['Stephen Young', 'Senior Policy Officer'],
                ['John Cain', 'IT Financial Manager'],
                ['Murray Parsons', 'Service & Project Management Office Manager']]
            ]
        ]

    loader = uploader(dbpath, data, report_date)
    
    def add_people(temp):
        for person in temp[1]:
            loader.add_management_line(temp[0], person)
        loader.batch_flush(loader.write_batch)


    for each in data:
        add_people(each)

    # add_people()


