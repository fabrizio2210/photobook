from pymongo import MongoClient
import uuid
import logging
import inspect

class db():
  client = None
  db_name = None


  @classmethod
  def create_db(cls):
    cls.client[cls.db_name]


  @classmethod
  def set_db_url(cls, db_url):
    logging.info("Setup connection string")
    cls.client = MongoClient(db_url, connect=False)


  @classmethod
  def set_db_name(cls, db_name):
    cls.db_name = db_name


  @classmethod
  def present(cls):
    logging.info("Verifying DB presence")
    return cls.db_name in cls.client.list_database_names()


  @classmethod
  def delete_all(cls):
    logging.info("Client: %s", cls.client)
    if cls.client is not None:
      for db_name in cls.client.list_database_names():
        if db_name not in ['admin', 'config', 'local']:
          logging.info("Dropping %s", db_name)
          cls.client.drop_database(db_name)

  
  class Model():
    def __init__(self):
      pass


    @classmethod
    def get_class_attrs(cls):
      attributes = inspect.getmembers(cls, lambda a:not(inspect.isroutine(a))) 
      return [a for a in attributes if not(a[0].startswith('__') and a[0].endswith('__'))]


    def get_attrs(self):
      attributes = inspect.getmembers(self, lambda a:not(inspect.isroutine(a))) 
      return [a for a in attributes if not(a[0].startswith('__') and a[0].endswith('__'))]


    def __repr__(self):
      ans = ""
      for attr in self.get_attrs():
        ans += '%s => "%s"\n' % (str(attr[0]), str(attr[1]))
      return ans


    def json(self):
      attrs = self.get_attrs()
      jsn = {}
      for a in attrs:
        if a[0].startswith('_'):
          continue
        jsn[a[0]] = a[1]
      return jsn


    @classmethod
    def from_json(cls, jsn):
      cls_attrs = cls.get_class_attrs()
      return cls(**jsn)


    @classmethod
    def create_table(cls):
      db_instance = db.client[db.db_name]
      db_instance[cls.__tablename__]


    @classmethod
    def find(cls, **kwargs):
      db_instance = db.client[db.db_name]
      table = db_instance[cls.__tablename__]
      logging.debug("====== Start Query ======")
      logging.debug("Find query: %s", kwargs)
      cursor = table.find(kwargs)
      if cls.__sort__:
        cursor.sort(*cls.__sort__)
      rows = list(cursor)
      logging.debug("In the function: %s", rows)
      logging.debug("======== End Query ======")
      res = []
      if len(rows) > 0:
        for row in rows:
          res.append(cls.from_json(row))
      return res


    @classmethod
    def greaterThan(cls, **kwargs):
      db_instance = db.client[db.db_name]
      table = db_instance[cls.__tablename__]
      logging.debug("====== Start Query ======")
      logging.debug("Find gt query: %s", kwargs)
      query = {}
      for arg in kwargs.keys():
        query[arg] = { "$gt": kwargs[arg] }
      logging.debug("Find gt actual query: %s", query)
      cursor = table.find(query)
      if cls.__sort__:
        cursor.sort(*cls.__sort__)
      logging.debug("======== End Query ======")
      attrs = cls.get_class_attrs()
      attrs.sort()
      rows = list(cursor)
      res = []
      if len(rows) > 0:
        for row in rows:
          res.append(cls.from_json(row))
      return res


    def save(self):
      db_instance = db.client[db.db_name]
      table = db_instance[self.__tablename__]
      if self.id is not None:
        table.replace_one({'id': self.id}, self.json(), True)
      else:
        self.id = str(uuid.uuid4())
        table.insert_one(self.json())


    def delete(self):
      db_instance = db.client[db.db_name]
      table = db_instance[self.__tablename__]
      query = { 'id': self.id }
      logging.debug("====== Start Query ======")
      logging.debug("Delete query: %s", query)
      cursor = table.delete_one(query)
      logging.debug("======== End Query ======")
