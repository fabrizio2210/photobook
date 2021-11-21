import sqlite3
import logging
import inspect
import os

class db():
  db_file = None

  @classmethod
  def create_tables(cls):
    pass

  @classmethod
  def set_db_filename(cls, db_filename):
    cls.db_file = db_filename

  @classmethod
  def present(cls):
    if os.path.isfile(cls.db_file):
      return True
    return False

  @classmethod
  def delete_all(cls):
    if os.path.isfile(cls.db_file):
      os.unlink(cls.db_file)
  
  class Field():
    def __init__(self, _type = None, primary_key = None):
      if primary_key == None:
          primary_key = False
      if _type == None:
        _type = "string"
      self.primary_key = primary_key
      if _type == "string":
        self._type = "text"
      if _type == "integer":
        _type = "INTEGER"
      self._type = _type

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

    @classmethod
    def create_table(cls):
      attrs = cls.get_class_attrs()
      attrs.sort()
      ph = "("
      for col in attrs:
        ph += "{name} {_type} ".format(name=col[0],_type=col[1]._type)
        if col[1].primary_key:
          ph += " PRIMARY KEY"
        ph += ","
      ph = ph.rstrip(',')
      ph += ")"
      create_table = "CREATE TABLE IF NOT EXISTS {table} {placeholder}".format(table=cls.__tablename__, placeholder=ph)
      logging.debug("====== Start Query ======")
      logging.debug(create_table) 
      logging.debug("======== End Query ======")
      connection = sqlite3.connect(db.db_file)
      cursor = connection.cursor()
      cursor.execute(create_table)
      connection.commit()
      connection.close()

    @classmethod
    def find(cls, **kwargs):
      list_args = []
      where_args = ""
      query = "SELECT * FROM {table} ".format(table=cls.__tablename__)
      query += "WHERE "
      for arg in kwargs:
        query += " %s=? AND " % arg
        list_args.append(kwargs[arg])
      query += "1 = 1"
      logging.debug("====== Start Query ======")
      logging.debug(query)
      logging.debug(list_args)
      logging.debug("======== End Query ======")
      connection = sqlite3.connect(db.db_file)
      cursor = connection.cursor()
      result = cursor.execute(query, (*list_args,))
      rows = result.fetchall()
      connection.close()
      attrs = cls.get_class_attrs()
      attrs.sort()
      res = []
      if len(rows)>0:
        for row in rows:
          # mapping correctly 
          res.append(cls(**dict(zip([a[0] for a in attrs],[*row]))))
      return res

    @classmethod
    def greaterThan(cls, **kwargs):
      list_args = []
      where_args = ""
      query = "SELECT * FROM {table} ".format(table=cls.__tablename__)
      query += "WHERE "
      for arg in kwargs:
        query += " %s>? AND " % arg
        list_args.append(kwargs[arg])
      query += "1 = 1"
      logging.debug("====== Start Query ======")
      logging.debug(query)
      logging.debug(list_args)
      logging.debug("======== End Query ======")
      connection = sqlite3.connect(db.db_file)
      cursor = connection.cursor()
      result = cursor.execute(query, (*list_args,))
      rows = result.fetchall()
      connection.close()
      attrs = cls.get_class_attrs()
      attrs.sort()
      res = []
      if len(rows)>0:
        for row in rows:
          # mapping correctly 
          res.append(cls(**dict(zip([a[0] for a in attrs],[*row]))))
      return res

    def save(self):
      # Check if it is an update
      use_id = False
      class_attrs = self.__class__.get_class_attrs()
      for col in class_attrs:
        if col[1].primary_key and col[1]._type == "INTEGER":
          id_value = getattr(self, col[0])
          if id_value is not None:
            id_col = col[0]
            use_id = True
      
      attrs = self.get_attrs()
      attrs.sort()

      # Build query for insertion
      rowcount = 0
      if use_id:
        # if the row to update exists
        query = "SELECT * FROM {table} WHERE {id_col} = {id_value}".format(table=self.__tablename__, id_col=id_col, id_value=id_value)
        # Execution of the preselect
        logging.debug("====== Start Preselect Query ======")
        logging.debug(query)
        connection = sqlite3.connect(db.db_file)
        cursor = connection.cursor()
        result = cursor.execute(query)
        connection.commit()
        rowcount = len(result.fetchall())
        connection.close()
        logging.debug("====== Stop Preselect Query ======")
        logging.debug("rowcount: %d, rowid: %d" % (rowcount, cursor.lastrowid))

      if use_id and rowcount > 0:
        # remove the primary_key column if using UPDATE
        attrs = [ a for a in attrs if not a[0] == id_col]
        ph = ""
        for attr in attrs:
          ph +=  " %s = ?," % attr[0]
        ph = ph.rstrip(',')
        query = "UPDATE {table} SET {placeholder} WHERE {id_col} = {id_value}".format(table=self.__tablename__, placeholder=ph, id_col=id_col, id_value=id_value)
      else:
        ph = "(" + ",".join(['?']*len(attrs)) + ")"
        query = "INSERT INTO {table} VALUES {placeholder}".format(table=self.__tablename__, placeholder=ph)

      row = (*[a[1] for a in attrs],)
      # Execution of the query
      logging.debug("====== Start Query ======")
      logging.debug(query)
      logging.debug(row)
      connection = sqlite3.connect(db.db_file)
      cursor = connection.cursor()
      cursor.execute(query, row)
      connection.commit()
      connection.close()

      # Update instance with new ID
      # if it is a new insertion
      if not use_id:
        class_attrs = self.__class__.get_class_attrs()
        for col in class_attrs:
          if col[1].primary_key and col[1]._type == "INTEGER":
            setattr(self, col[0], cursor.lastrowid)
            self.id = cursor.lastrowid
            logging.debug("rowid: %d" % self.id)
      logging.debug("======== End Query ======")

    def delete(self):
      attrs = self.get_attrs()
      list_args = []
      where_args = ""
      query = "DELETE FROM {table} ".format(table=self.__tablename__)
      query += "WHERE "
      if len(attrs) == 0:
        logging.warning("Error: trying to drop table")
        return

      # Check if an ID exists
      use_id = False
      class_attrs = self.__class__.get_class_attrs()
      for col in class_attrs:
        if col[1].primary_key and col[1]._type == "INTEGER":
          list_args.append(getattr(self, col[0]))
          query += " %s=? " % col[0]
          use_id = True

      if not use_id:
        for arg in attrs:
          query += " %s=? AND " % arg[0]
          list_args.append(arg[1])
        query += "1 = 1"

      logging.debug("====== Start Query ======")
      logging.debug(query)
      logging.debug(list_args)
      logging.debug("======== End Query ======")
      connection = sqlite3.connect(db.db_file)
      cursor = connection.cursor()
      cursor.execute(query, (*list_args,))
      connection.commit()
      connection.close()
      logging.debug("rowcount: %d, rowid: %d" % (cursor.rowcount, cursor.lastrowid))
