#!flask/bin/python

'''
Neo4j staff directory
'''

from flask import Flask, jsonify, abort, request, make_response, url_for, send_file
from flask.views import MethodView
from flask.ext.restful import Api, Resource, reqparse, fields, marshal
import re
from py2neo import neo4j, rel
import logging
from time import sleep
import sys

from werkzeug.contrib.fixers import ProxyFix

app = Flask(__name__,)
api = Api(app)

person_fields = {
    'id': fields.String,
    'name': fields.String,
    'email': fields.String,
    'department': fields.String,
    'mobile': fields.String,
    'phone': fields.String,
    'position': fields.String
}

area_fields = {
    'name': fields.String
}

class neo4j_db(object):
    """Connect to neo4j database"""
    def __init__(self, dbpath='http://localhost:7470/db/data'):
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
            sleep(2)

    def get_person(self, query):
        ''' figure out the query structure '''
        def process_loc(loc):
            its = re.compile(r'(?i)(ITS|Information Technology Services)')
            eng = re.compile(r'(?i)(Engineering|.*eng.*)')
            if its.search(loc) is not None:
                return 'ITS|Information Technology Services'
            elif eng.search(loc) is not None:
                return 'Engineering'
            return '(?i).*%s.*' % '.*'.join(loc.lower().split())

        def process_name(name, loc=None):
            nquery = 'name'
            if '@' in name:
                nquery = '%s="%s"' % ("email",name)
            elif re.match('^[0-9 ()+]{4,15}$',name):
                nquery = '%s=~".*%s.*" OR a.%s =~".*%s.*"' % ("phone",name, "mobile",name)

            else:
                if ' ' in name:
                   name = '.* '.join(name.split())
                nquery = '%s=~"(?i)^%s.*"' % ('name',name)
            if loc:
                nquery += ' AND a.%s=~"%s"' % ("department", process_loc(loc))
            return nquery

        def get_data(query):
            temp = []
            for x in query.stream():
                print(x)
                
                a = x[1].get_properties()
                a['id'] = x[0]
                print(a)
                temp.append(a)
            return temp

        if "'s " not in query: 
            if ' in ' in query:
                temp = query.split(' in ')
                nquery = process_name(*temp)
            else:
                nquery = process_name(query)
            print(nquery)
            cypher = neo4j.CypherQuery(self.db, 'MATCH (a:People) WHERE a.%s RETURN id(a), a LIMIT 500' % nquery )
        else:
            temp = query.split("'s ")
            
            if temp[1] == 'manager':
                print("%s , %s" % (temp[0], process_name(temp[0])))
                cypher = neo4j.CypherQuery(self.db, 'MATCH (a:People)-[:MANAGES]->(b:People) WHERE b.%s RETURN id(a), a LIMIT 500' % process_name(temp[0]))
            elif temp[1] == 'colleagues':
                cypher = neo4j.CypherQuery(self.db, 'MATCH (b:People)<-[:MANAGES]-(a:People)-[:MANAGES]->(c:People) WHERE b.%s RETURN id(c), c' % process_name(temp[0]))
        
        temp = get_data(cypher)
        if len(temp) == 0:
            cypher = neo4j.CypherQuery(self.db, 'MATCH (a:People) WHERE a.position=~"(?i).*%s.*" RETURN id(a), a LIMIT 1000' % query) 
            temp = get_data(cypher)

        return temp

    def get_node(self, node):
        def get_data(query):
            temp = []
            for x in query.stream():
                a = x[0].get_properties()
                a['id'] = node
                temp.append(a)
            return temp

        cypher = neo4j.CypherQuery(self.db, 'START n=node(%s) RETURN n' % node )
        temp = get_data(cypher)
        if len(temp):
            temp[0]['manager'] = get_data(neo4j.CypherQuery(self.db, 'MATCH (a)-[:MANAGES]->(b) WHERE id(b) = %s RETURN a' % node ))
            temp[0]['colleagues'] = get_data(neo4j.CypherQuery(self.db, 'MATCH (c)<-[:MANAGES]-(a)-[:MANAGES]->(b) WHERE id(b) = %s RETURN c' % node ))
            temp[0]['reports'] = get_data(neo4j.CypherQuery(self.db, 'MATCH (a)-[:MANAGES]->(b) WHERE id(a) = %s RETURN b' % node ))
        print(temp)
        return temp


    def get_similarpages(self, url):
        url = url[:-1]
        node_type = 'Pages'
        if '.pdf' in url or '.doc' in url or '.ppt' in url:
            node_type = 'Documents'
        query = neo4j.CypherQuery(self.db, 'MATCH (x:Pages)-[r:LINKS_TO]-(y:%s) WHERE y.url=~"^%s.*" AND x.url=~"^%s.*" RETURN y, count(r) ORDER BY count(r) DESC LIMIT 50' % (node_type, url,url.split('/',1)[0]) )
        temp = []
        for x in query.stream():
            a = x.get_properties()
            temp.append({"url": a['url']})
        return temp

class searchPeople(Resource):
    # decorators = [auth.login_required]

    def __init__(self):
        self.db = neo4j_db()
        self.db.connect()
        self.reqparse = reqparse.RequestParser()
        self.reqparse.add_argument('query', type=str)
        super(searchPeople, self).__init__()

    def get(self):
        if 'q' not in request.args:
            abort(404)
        data = self.db.get_person(request.args['q'])
        print('here', data, request.args['q'])
        length = len(data)
        if length == 0:
            return {"data" : [], "status": "sorry, no results found.", "size": length}

        return {"data" : [marshal(x, person_fields) for x in data ], "status": "success", "size": length}

class GetPerson(Resource):
    # decorators = [auth.login_required]

    def __init__(self):
        self.db = neo4j_db()
        self.db.connect()
        self.reqparse = reqparse.RequestParser()
        super(GetPerson, self).__init__()

    def get(self, person):
        data = self.db.get_node(person)
        length = len(data)
        if length == 0:
            return {"data" : [], "status": "sorry, no results found.", "size": length}

        return {"data" : [marshal(x, person_fields) for x in data ], "status": "success", "size": length}

class GetDepartment(Resource):
    # decorators = [auth.login_required]

    def __init__(self):
        self.db = neo4j_db()
        self.db.connect()
        self.reqparse = reqparse.RequestParser()
        self.reqparse.add_argument('hostname', type=str)
        super(GetDepartment, self).__init__()

    def get(self):
        if 'q' not in request.args:
            abort(404)
        data = self.db.get_similarpages(request.args['q'])
        print('here', data, request.args['q'])
        length = len(data)
        if length == 0:
            return {"data" : [], "status": "sorry, no results found.", "size": length}

        return {"data" : [marshal(x, page_fields) for x in data ], "status": "success", "size": length}


class GetComplex(Resource):
    # decorators = [auth.login_required]

    def __init__(self):
        self.db = neo4j_db()
        self.db.connect()
        self.reqparse = reqparse.RequestParser()
        self.reqparse.add_argument('hostname', type=str)
        super(GetDepartment, self).__init__()

    def get(self):
        if 'q' not in request.args:
            abort(404)
        data = self.db.get_similarpages(request.args['q'])
        print('here', data, request.args['q'])
        length = len(data)
        if length == 0:
            return {"data" : [], "status": "sorry, no results found.", "size": length}

        return {"data" : [marshal(x, page_fields) for x in data ], "status": "success", "size": length}


api.add_resource(searchPeople, '/webstructure/api/v1.0/search', endpoint = 'search')
api.add_resource(GetPerson, '/webstructure/api/v1.0/person/<string:person>')
api.add_resource(GetDepartment, '/webstructure/api/v1.0/department', endpoint = 'department')
api.add_resource(GetComplex, '/webstructure/api/v1.0/complex', endpoint = 'complex')


app.wsgi_app = ProxyFix(app.wsgi_app)

if __name__ == '__main__':
    app.run(debug = True)