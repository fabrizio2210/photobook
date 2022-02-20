import unittest
import logging
import time
import json
import sys
import warnings
from os import path
sys.path.append( path.dirname( path.dirname( path.abspath(__file__) ) ) )
import app
from utility.data import bootstrap


def initialize_test():
  bootstrap(force=True, dev=True, quiet=True)

class TestAPI_without_auth(unittest.TestCase):

  def setUp(self):
    app.app.testing = True
    self.app = app.app.test_client()
    initialize_test()

  def test_010_root(self):
    rv = self.app.get('/')
    self.assertEqual(rv.status, '404 NOT FOUND')


if __name__ == '__main__':
  logging.basicConfig(level=logging.INFO)
  unittest.main()
