import sqlite3
import logging
from time import sleep
import sys
from py2neo import neo4j, rel
from io import StringIO


logging.basicConfig(level=logging.INFO)

class neo4j_db(object):
    """Connect to neo4j database"""
    def __init__(self, dbpath):
        self.dbpath = dbpath

    def _db_init(self):
        logging.debug("Trying to connect to db.")
        try:
            self.db = neo4j.GraphDatabaseService(self.dbpath)
            self.write_batch = neo4j.WriteBatch(self.db)
            return True
        except Exception as e:
            logging.debug("db_init: %s" % e)
            return False

    def connect(self):
        i = 0
        while True:
            a = self._db_init()
            if a:
                break
            if i == 5:
                logging.error("db_connect: DB connection failed.")
                sys.exit(1)
            i += 1
            sleep(3)


class sqlite(object):
    '''Use sqlite3 for intermediate storage'''
    def __init__(self, dbfile, opt=False):
        if opt:
            self._db = self._connect_opt(dbfile)
        else:
            self._db = self._connect(dbfile)
        self._cursor = self._db.cursor()
        self._create_tables()
        self._counter = 100
        self._queue_size = 100
        self._temp_queue = []
        pass

    def _connect(self, dbfile):
        try:
            return sqlite3.connect(dbfile)
        except:
            logging.error("Failed to connect/open db file.")
        return

    def _connect_opt(self, dbfile):
        con = sqlite3.connect(dbfile)
        tempfile = StringIO()
        for line in con.iterdump():
            tempfile.write('%s\n' % line)
        con.close()
        tempfile.seek(0)

        # Create a database in memory and import from tempfile
        newdb = sqlite3.connect(":memory:")
        newdb.cursor().executescript(tempfile.read())
        newdb.commit()
        newdb.row_factory = sqlite3.Row
        return newdb

    def _create_tables(self):
        self._cursor.execute('create table if not exists hostnames (hostname text primary key, node int)')
        self._cursor.execute('create table if not exists pages (url text primary key, path text, node int)')
        # self._cursor.execute('create table if not exists contains (hostname references hostnames(hostname), page references pages(path))')
        self._cursor.execute('create table if not exists contains (hostname references hostnames(hostname), page references pages(url))')
        self._cursor.execute('create table if not exists links (from_page references pages(url), to_page references pages(url), anchor text)')
        self._cursor.execute('PRAGMA journal_mode=MEMORY')
        #self._cursor.execute('PRAGMA synchronous = OFF')
        self._cursor.execute('PRAGMA default_cache_size=10000')
        self._db.commit()

    def _commit(self):
        if not self._counter:
            self._counter = 100
            self._db.commit()
            return
        self._counter -= 1

    def _queue(self, items, flush=False):
        self._temp_queue.append(items)
        if self._queue_size == len(self._temp_queue) or flush:
            return self._temp_queue
        return

    def get_pages_size(self):
        self._cursor.execute('select count(*) from pages')
        return [x for x in self._cursor][0][0]


    def add_host(self, items):
        self._cursor.execute('insert into hostnames values(?, ?)', items)
        self._commit()

    def add_page(self, items):
        self._cursor.execute('insert into pages values(?, ?, ?)', items)
        self._commit()

    def add_contains(self, items):
        self._cursor.execute('insert into contains values(?, ?)', items)
        self._commit()

    def add_link(self, items):
        self._cursor.execute('insert into links values(?, ?, ?)', items)
        self._commit()

    def get_hosts(self, search=False):
        if not search:
            self._cursor.execute('select hostname from hostnames')
        else:
            self._cursor.execute('select * from hostnames where hostname like ?', ('%' + search + '%',))
        return [x for x in self._cursor]

    def get_pages(self, host=False, offset=0, limit=-1, search=False):
        if not search:
            self._cursor.execute('select * from pages LIMIT ? OFFSET ?', (limit, offset))
        else:
            self._cursor.execute('select * from pages where url=?', (host + search,))
        return [x for x in self._cursor]

    def get_links(self, hostname):
        self._cursor.execute('select links.from_page, pages.url, links.anchor from pages, links where pages.url = links.to_page and links.from_page in (select page from contains where hostname=?)', (hostname,))
        return [x for x in self._cursor]

    def update_node(self, host, node):
        self._cursor.execute('update hostnames set node =? where hostname=?', (node, host,))
        self._commit()

    def update_page(self, url, node):
        self._cursor.execute('update pages set node =? where url=?', (node, url,))
        self._commit()

    def close(self):
        self._db.commit()
        self._db.close()


