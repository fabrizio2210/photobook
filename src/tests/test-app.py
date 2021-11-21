import unittest
import logging
import time
import json
import sys
import warnings
from os import path
sys.path.append( path.dirname( path.dirname( path.abspath(__file__) ) ) )
import app
from utils.data import bootstrap

admin_user = "fabrizio"
admin_password = "pwd"

normal_user = "fabrizio2"
normal_password = "pwd2"

second_user = "fabrizio3"
second_password = "pwd2"

def initialize_test():
  bootstrap(force=True, dev=True, quiet=True)
#########
# General without authentication

class TestAPI_without_auth(unittest.TestCase):

  def setUp(self):
    app.app.testing = True
    self.app = app.app.test_client()
    initialize_test()

  def test_010_root(self):
    rv = self.app.get('/')
    self.assertEqual(rv.status, '404 NOT FOUND')

  def test_020_get_token(self):
    # An invalid credential raises an error in the log
    logging.disable(logging.CRITICAL)
    logins = [ {"cred": { "username": admin_user, "password": admin_password}, "res": "200 OK"},
               { "cred": { "username": normal_user, "password": normal_password}, "res": "200 OK"},
               { "cred": { "username": "not_exist_s345ffsdg", "password": "not_exist"}, "res": "401 UNAUTHORIZED"},
               { "cred": { "username": admin_user, "password": "not_exist"}, "res": "401 UNAUTHORIZED"} ]
    for login in logins:
      rv = self.app.post('/api/auth', json = login['cred'])
      self.assertEqual(rv.status, login['res'])
    logging.disable(logging.NOTSET)


##############
# Main Setting

class TestAPI_MainSettingAsAdmin(unittest.TestCase):

  def setUp(self):
    app.app.testing = True
    self.app = app.app.test_client()
    initialize_test()
    rv = self.app.post('/api/auth', json = { "username": admin_user, "password": admin_password})
    self.headers = {'Content-Type': 'application/json', 'Authorization': "JWT " + json.loads(rv.data.decode("utf-8"))['access_token']}

  def test_020_get_settings(self):
    rv = self.app.get('/api/v1/settings', headers = self.headers)
    self.assertEqual(rv.status, '200 OK')

  def test_021_get_setting(self):
    rv = self.app.get('/api/v1/setting/max_project_run_number', headers = self.headers)
    self.assertEqual(rv.status, '200 OK')
    self.assertEqual(json.loads(rv.data.decode("utf-8")), {"id": 1,
      "name": "max_project_run_number",
      "description": "Number of executions per project that are stored",
      "value": None,
      "default_value": 20})

  def test_027_put_setting(self):
    rv = self.app.put('/api/v1/setting/max_project_run_number', json = {"value": 10}, headers = self.headers)
    self.assertEqual(rv.status, '200 OK')
    rv = self.app.get('/api/v1/setting/max_project_run_number', headers = self.headers)
    self.assertEqual(rv.status, '200 OK')
    self.assertEqual(json.loads(rv.data.decode("utf-8")), {"id": 1,
      "name": "max_project_run_number",
      "description": "Number of executions per project that are stored",
      "value": 10,
      "default_value": 20})


class TestAPI_MainSettingasUser(unittest.TestCase):

  def setUp(self):
    app.app.testing = True
    self.app = app.app.test_client()
    initialize_test()
    rv = self.app.post('/api/auth', json = { "username": normal_user, "password": normal_password})
    self.headers = {'Content-Type': 'application/json', 'Authorization': "JWT " + json.loads(rv.data.decode("utf-8"))['access_token']}

  def test_020_get_settings(self):
    rv = self.app.get('/api/v1/settings', headers = self.headers)
    self.assertEqual(rv.status, '200 OK')

  def test_021_get_setting(self):
    rv = self.app.get('/api/v1/setting/max_project_run_number', headers = self.headers)
    self.assertEqual(rv.status, '200 OK')
    self.assertEqual(json.loads(rv.data.decode("utf-8")), {"id": 1,
      "name": "max_project_run_number",
      "description": "Number of executions per project that are stored",
      "value": None,
      "default_value": 20})

  def test_025_put_setting_forbidden(self):
    rv = self.app.put('/api/v1/setting/max_project_run_number', json = {"value": 1}, headers = self.headers)
    self.assertEqual(rv.status, '403 FORBIDDEN')
    rv = self.app.get('/api/v1/setting/max_project_run_number', headers = self.headers)
    self.assertEqual(rv.status, '200 OK')
    self.assertEqual(json.loads(rv.data.decode("utf-8")), {"id": 1,
      "name": "max_project_run_number",
      "description": "Number of executions per project that are stored",
      "value": None,
      "default_value": 20})


##############################
# Project and Project Settings

class TestAPI_ProjectAsUser(unittest.TestCase):

  def setUp(self):
    self.verificationErrors = []
    app.app.testing = True
    self.maxDiff = None
    self.app = app.app.test_client()
    initialize_test()
    rv = self.app.post('/api/auth', json = { "username": normal_user, "password": normal_password})
    self.headers = {'Content-Type': 'application/json', 'Authorization': "JWT " + json.loads(rv.data.decode("utf-8"))['access_token']}

  def tearDown(self):
    self.assertEqual([], self.verificationErrors)

      # GET  /api/v1/project/1
      # GET  /api/v1/projects
      # POST /api/v1/new_project
  def test_030_project(self):
    prj_name = "my_prj_name"
    prj_name2 = "my_prj_name2"
    # Try to get project
    rv = self.app.get("/api/v1/project/1", headers = self.headers)
    try: self.assertEqual(rv.status, '404 NOT FOUND')
    except AssertionError as e: self.verificationErrors.append(str(e) + "Line: " + str(sys.exc_info()[2].tb_lineno)) 
    # Create project
    rv = self.app.post("/api/v1/new_project", json = {"name" : prj_name}, headers = self.headers)
    try: self.assertEqual(rv.status, '201 CREATED')
    except AssertionError as e: self.verificationErrors.append(str(e) + "Line: " + str(sys.exc_info()[2].tb_lineno)) 
    try: self.assertEqual(json.loads(rv.data.decode("utf-8")), { "name": prj_name, "id": 1})
    except AssertionError as e: self.verificationErrors.append(str(e) + "Line: " + str(sys.exc_info()[2].tb_lineno)) 
    # Get project
    rv = self.app.get("/api/v1/project/1", json = {}, headers = self.headers)
    try: self.assertEqual(rv.status, '200 OK')
    except AssertionError as e: self.verificationErrors.append(str(e) + "Line: " + str(sys.exc_info()[2].tb_lineno)) 
    try: self.assertEqual(json.loads(rv.data.decode("utf-8")), { "project": { "name": prj_name, "id": 1 } } )
    except AssertionError as e: self.verificationErrors.append(str(e) + "Line: " + str(sys.exc_info()[2].tb_lineno)) 
    # Create second project as second user
    rv = self.app.post('/api/auth', json = { "username": second_user, "password": second_password})
    self.headers = {'Content-Type': 'application/json', 'Authorization': "JWT " + json.loads(rv.data.decode("utf-8"))['access_token']}
    rv = self.app.post("/api/v1/new_project", json = {"name": prj_name2}, headers = self.headers)
    try: self.assertEqual(rv.status, '201 CREATED')
    except AssertionError as e: self.verificationErrors.append(str(e) + "Line: " + str(sys.exc_info()[2].tb_lineno)) 
    try: self.assertEqual(json.loads(rv.data.decode("utf-8")), { "name": prj_name2, "id": 2})
    except AssertionError as e: self.verificationErrors.append(str(e) + "Line: " + str(sys.exc_info()[2].tb_lineno)) 
    # Get all project of the second user
    rv = self.app.get("/api/v1/projects", json = {}, headers = self.headers)
    try: self.assertEqual(rv.status, '200 OK')
    except AssertionError as e: self.verificationErrors.append(str(e) + "Line: " + str(sys.exc_info()[2].tb_lineno)) 
    try: self.assertEqual(json.loads(rv.data.decode("utf-8")), { "projects": [ { "name": prj_name2, "id": 2 } ] })
    except AssertionError as e: self.verificationErrors.append(str(e) + "Line: " + str(sys.exc_info()[2].tb_lineno)) 

      # GET/PUT /api/v1/project/1/setting/scm_url
  def test_050_get_setting_project(self):
    prj_name = "my_prj_name"
    prj_name2 = "my_prj_name2"
    value = "http://github.com/fdfdsfsd/dfsddf.git"
    value2 = "http://github.com/fdfd/dff.git"
    # Try to get setting of a non-existant project
    rv = self.app.get('/api/v1/project/1/setting/scm_url', headers = self.headers)
    try: self.assertEqual(rv.status, '404 NOT FOUND')
    except AssertionError as e: self.verificationErrors.append(str(e) + "Line: " + str(sys.exc_info()[2].tb_lineno)) 
    # Create first project
    rv = self.app.post("/api/v1/new_project", json = {"name": prj_name}, headers = self.headers)
    try: self.assertEqual(rv.status, '201 CREATED')
    except AssertionError as e: self.verificationErrors.append(str(e) + "Line: " + str(sys.exc_info()[2].tb_lineno)) 
    try: self.assertEqual(json.loads(rv.data.decode("utf-8")), { "name": prj_name, "id": 1})
    except AssertionError as e: self.verificationErrors.append(str(e)) 
    # Get setting of first project
    rv = self.app.get('/api/v1/project/1/setting/scm_url', headers = self.headers)
    try: self.assertEqual(rv.status, '200 OK')
    except AssertionError as e: self.verificationErrors.append(str(e) + "Line: " + str(sys.exc_info()[2].tb_lineno)) 
    try: self.assertEqual(json.loads(rv.data.decode("utf-8")), { "setting": {
      "name": "scm_url", 
      "id": None, 
      "description": "Source Control Manager URL of the project",
      "value": None,
      "default_value": None }})
    except AssertionError as e: self.verificationErrors.append(str(e) + "Line: " + str(sys.exc_info()[2].tb_lineno)) 
    # Put setting of first project
    payload = { "value": value }
    rv = self.app.put('/api/v1/project/1/setting/scm_url', json = payload, headers = self.headers)
    try: self.assertEqual(rv.status, '200 OK')
    except AssertionError as e: self.verificationErrors.append(str(e) + "Line: " + str(sys.exc_info()[2].tb_lineno)) 
    try: self.assertEqual(json.loads(rv.data.decode("utf-8")), { "setting":
      {"name": "scm_url", 
      "id": 5, 
      "description": "Source Control Manager URL of the project",
      "value": value,
      "default_value": None }})
    except AssertionError as e: self.verificationErrors.append(str(e) + "Line: " + str(sys.exc_info()[2].tb_lineno)) 
    # Create second project as second user
    rv = self.app.post('/api/auth', json = { "username": second_user, "password": second_password})
    self.headers = {'Content-Type': 'application/json', 'Authorization': "JWT " + json.loads(rv.data.decode("utf-8"))['access_token']}
    rv = self.app.post("/api/v1/new_project", json = {"name": prj_name2}, headers = self.headers)
    try: self.assertEqual(rv.status, '201 CREATED')
    except AssertionError as e: self.verificationErrors.append(str(e) + "Line: " + str(sys.exc_info()[2].tb_lineno)) 
    try: self.assertEqual(json.loads(rv.data.decode("utf-8")), { "name": prj_name2, "id": 2})
    except AssertionError as e: self.verificationErrors.append(str(e) + "Line: " + str(sys.exc_info()[2].tb_lineno)) 
    # Try to get setting of the first as second user
    rv = self.app.get('/api/v1/project/1/setting/scm_url', headers = self.headers)
    try: self.assertEqual(rv.status, '403 FORBIDDEN')
    except AssertionError as e: self.verificationErrors.append(str(e) + "Line: " + str(sys.exc_info()[2].tb_lineno)) 
    # Try to modify as second user
    rv = self.app.post('/api/auth', json = { "username": second_user, "password": second_password})
    self.headers = {'Content-Type': 'application/json', 'Authorization': "JWT " + json.loads(rv.data.decode("utf-8"))['access_token']}
    payload = { "value": value2 }
    rv = self.app.put('/api/v1/project/1/setting/scm_url', json = payload, headers = self.headers)
    try: self.assertEqual(rv.status, '403 FORBIDDEN')
    except AssertionError as e: self.verificationErrors.append(str(e) + "Line: " + str(sys.exc_info()[2].tb_lineno)) 
    # Put non existantsetting of second project as second user
    payload = { "value": value }
    rv = self.app.put('/api/v1/project/2/setting/non_existent_setting_Dfdsfs', json = payload, headers = self.headers)
    try: self.assertEqual(rv.status, '404 NOT FOUND')
    except AssertionError as e: self.verificationErrors.append(str(e) + "Line: " + str(sys.exc_info()[2].tb_lineno)) 
    # Try to get setting of the second as second user
    rv = self.app.get('/api/v1/project/2/setting/non_existent_setting_Dfdsfs', headers = self.headers)
    try: self.assertEqual(rv.status, '404 NOT FOUND')
    except AssertionError as e: self.verificationErrors.append(str(e) + "Line: " + str(sys.exc_info()[2].tb_lineno)) 

      # PUT /api/v1/project/1/setting/scm_url
  def test_035_put_project_setting(self):
    prj_name = "my_prj_name"
    value = "http://github.com/fdfdsfsd/dfsddf.git"
    value2 = "http://github.com/fdfd/dff.git"
    # Create project
    rv = self.app.post("/api/v1/new_project", json = {"name": prj_name}, headers = self.headers)
    try: self.assertEqual(rv.status, '201 CREATED')
    except AssertionError as e: self.verificationErrors.append(str(e) + "Line: " + str(sys.exc_info()[2].tb_lineno)) 
    try: self.assertEqual(json.loads(rv.data.decode("utf-8")), { "name": prj_name, "id": 1})
    except AssertionError as e: self.verificationErrors.append(str(e) + "Line: " + str(sys.exc_info()[2].tb_lineno)) 
    # Put setting
    payload = { "value": value }
    rv = self.app.put('/api/v1/project/1/setting/scm_url', json = payload, headers = self.headers)
    try: self.assertEqual(rv.status, '200 OK')
    except AssertionError as e: self.verificationErrors.append(str(e) + "Line: " + str(sys.exc_info()[2].tb_lineno)) 
    try: self.assertEqual(json.loads(rv.data.decode("utf-8")), { "setting":
      {"name": "scm_url", 
      "id": 5, 
      "description": "Source Control Manager URL of the project",
      "value": value,
      "default_value": None }})
    except AssertionError as e: self.verificationErrors.append(str(e) + "Line: " + str(sys.exc_info()[2].tb_lineno)) 
    # Try to modify as second user
    rv = self.app.post('/api/auth', json = { "username": second_user, "password": second_password})
    self.headers = {'Content-Type': 'application/json', 'Authorization': "JWT " + json.loads(rv.data.decode("utf-8"))['access_token']}
    payload = { "value": value2 }
    rv = self.app.put('/api/v1/project/1/setting/scm_url', json = payload, headers = self.headers)
    try: self.assertEqual(rv.status, '403 FORBIDDEN')
    except AssertionError as e: self.verificationErrors.append(str(e) + "Line: " + str(sys.exc_info()[2].tb_lineno)) 
    # Get setting as first user
    rv = self.app.post('/api/auth', json = { "username": normal_user, "password": normal_password})
    self.headers = {'Content-Type': 'application/json', 'Authorization': "JWT " + json.loads(rv.data.decode("utf-8"))['access_token']}
    rv = self.app.get('/api/v1/project/1/setting/scm_url', headers = self.headers)
    try: self.assertEqual(rv.status, '200 OK')
    except AssertionError as e: self.verificationErrors.append(str(e) + "Line: " + str(sys.exc_info()[2].tb_lineno)) 
    try: self.assertEqual(json.loads(rv.data.decode("utf-8")), { "setting":
      {"name": "scm_url", 
      "id": 5, 
      "description": "Source Control Manager URL of the project",
      "value": value,
      "default_value": None }})
    except AssertionError as e: self.verificationErrors.append(str(e) + "Line: " + str(sys.exc_info()[2].tb_lineno)) 

##############################
# Environment

class TestAPI_EnvironmentAsUser(unittest.TestCase):

  def setUp(self):
    self.verificationErrors = []
    app.app.testing = True
    self.app = app.app.test_client()
    initialize_test()
    rv = self.app.post('/api/auth', json = { "username": normal_user, "password": normal_password})
    self.headers = {'Content-Type': 'application/json', 'Authorization': "JWT " + json.loads(rv.data.decode("utf-8"))['access_token']}

  def tearDown(self):
    self.assertEqual([], self.verificationErrors)

      # GET             /api/v1/project/1/environments
      # POST            /api/v1/project/1/new_environment
  def test_050_get_environment_project(self):
    prj_name = "my_prj_name"
    prj_name2 = "my_prj_name2"
    name = "MY_VARIABLE"
    value = "MY_BEAUTIFUL_VARIABLE"
    name2 = "MY_VARIABLE2"
    value2 = "MY_BEAUTIFUL_VARIABLE2"
    # Try to get environment of a non-existant project
    rv = self.app.get('/api/v1/project/1/environments', headers = self.headers)
    try: self.assertEqual(rv.status, '404 NOT FOUND')
    except AssertionError as e: self.verificationErrors.append(str(e) + "Line: " + str(sys.exc_info()[2].tb_lineno)) 
    # Create first project
    rv = self.app.post("/api/v1/new_project", json = {"name": prj_name}, headers = self.headers)
    try: self.assertEqual(rv.status, '201 CREATED')
    except AssertionError as e: self.verificationErrors.append(str(e) + "Line: " + str(sys.exc_info()[2].tb_lineno)) 
    try: self.assertEqual(json.loads(rv.data.decode("utf-8")), { "name": prj_name, "id": 1})
    except AssertionError as e: self.verificationErrors.append(str(e)) 
    # Get all enviroments of the project
    rv = self.app.get('/api/v1/project/1/environments', headers = self.headers)
    try: self.assertEqual(rv.status, '200 OK')
    except AssertionError as e: self.verificationErrors.append(str(e) + "Line: " + str(sys.exc_info()[2].tb_lineno)) 
    try: self.assertEqual(json.loads(rv.data.decode("utf-8")), { "environments" : [ ] })
    except AssertionError as e: self.verificationErrors.append(str(e) + "Line: " + str(sys.exc_info()[2].tb_lineno)) 
    # Post environment of first project
    payload = { "name": name, "value": value }
    rv = self.app.post('/api/v1/project/1/new_environment', json = payload, headers = self.headers)
    try: self.assertEqual(rv.status, '201 CREATED')
    except AssertionError as e: self.verificationErrors.append(str(e) + "Line: " + str(sys.exc_info()[2].tb_lineno)) 
    try: self.assertEqual(json.loads(rv.data.decode("utf-8")), { "environment": {
      "id": 1, 
      "description": None,
      "value": value,
      "name": name }})
    except AssertionError as e: self.verificationErrors.append(str(e) + "Line: " + str(sys.exc_info()[2].tb_lineno)) 
    # Create second project 
    rv = self.app.post("/api/v1/new_project", json = {"name": prj_name2}, headers = self.headers)
    try: self.assertEqual(rv.status, '201 CREATED')
    except AssertionError as e: self.verificationErrors.append(str(e) + "Line: " + str(sys.exc_info()[2].tb_lineno)) 
    try: self.assertEqual(json.loads(rv.data.decode("utf-8")), { "name": prj_name2, "id": 2})
    except AssertionError as e: self.verificationErrors.append(str(e) + "Line: " + str(sys.exc_info()[2].tb_lineno)) 
    # Get all enviroments of the first project
    rv = self.app.get('/api/v1/project/1/environments', headers = self.headers)
    try: self.assertEqual(rv.status, '200 OK')
    except AssertionError as e: self.verificationErrors.append(str(e) + "Line: " + str(sys.exc_info()[2].tb_lineno)) 
    try: self.assertEqual(json.loads(rv.data.decode("utf-8")), { "environments" : [
    { "id": 1, 
      "description": None,
      "value": value,
      "name": name } ] })
    except AssertionError as e: self.verificationErrors.append(str(e) + "Line: " + str(sys.exc_info()[2].tb_lineno)) 
    # Get all enviroments of the second project
    rv = self.app.get('/api/v1/project/2/environments', headers = self.headers)
    try: self.assertEqual(rv.status, '200 OK')
    except AssertionError as e: self.verificationErrors.append(str(e) + "Line: " + str(sys.exc_info()[2].tb_lineno)) 
    try: self.assertEqual(json.loads(rv.data.decode("utf-8")), { "environments" : [ ] })
    except AssertionError as e: self.verificationErrors.append(str(e) + "Line: " + str(sys.exc_info()[2].tb_lineno)) 

      # GET/PUT/DELETE  /api/v1/project/environment/1
  def test_035_put_environment(self):
    prj_name = "my_prj_name"
    name = "MY_VARIABLE"
    value = "MY_BEAUTIFUL_VARIABLE"
    name2 = "MY_VARIABLE2"
    value2 = "MY_BEAUTIFUL_VARIABLE2"
    # Create project
    rv = self.app.post("/api/v1/new_project", json = {"name": prj_name}, headers = self.headers)
    try: self.assertEqual(rv.status, '201 CREATED')
    except AssertionError as e: self.verificationErrors.append(str(e) + "Line: " + str(sys.exc_info()[2].tb_lineno)) 
    try: self.assertEqual(json.loads(rv.data.decode("utf-8")), {"name": prj_name, "id": 1})
    except AssertionError as e: self.verificationErrors.append(str(e) + "Line: " + str(sys.exc_info()[2].tb_lineno)) 
    # Post environment
    payload = { "name": name, "value": value }
    rv = self.app.post('/api/v1/project/1/new_environment', json = payload, headers = self.headers)
    try: self.assertEqual(rv.status, '201 CREATED')
    except AssertionError as e: self.verificationErrors.append(str(e) + "Line: " + str(sys.exc_info()[2].tb_lineno)) 
    try: self.assertEqual(json.loads(rv.data.decode("utf-8")), { "environment": {
      "name": name, 
      "id": 1, 
      "description": None,
      "value": value, }})
    except AssertionError as e: self.verificationErrors.append(str(e) + "Line: " + str(sys.exc_info()[2].tb_lineno)) 
    # Try to modify as second user
    rv = self.app.post('/api/auth', json = { "username": second_user, "password": second_password})
    self.headers = {'Content-Type': 'application/json', 'Authorization': "JWT " + json.loads(rv.data.decode("utf-8"))['access_token']}
    payload = { "value": value2 }
    rv = self.app.put('/api/v1/project/environment/1', json = payload, headers = self.headers)
    try: self.assertEqual(rv.status, '403 FORBIDDEN')
    except AssertionError as e: self.verificationErrors.append(str(e) + "Line: " + str(sys.exc_info()[2].tb_lineno)) 
    # Delete as first user
    rv = self.app.post('/api/auth', json = { "username": normal_user, "password": normal_password})
    self.headers = {'Content-Type': 'application/json', 'Authorization': "JWT " + json.loads(rv.data.decode("utf-8"))['access_token']}
    rv = self.app.delete('/api/v1/project/environment/1', headers = self.headers)
    try: self.assertEqual(rv.status, '200 OK')
    except AssertionError as e: self.verificationErrors.append(str(e) + "Line: " + str(sys.exc_info()[2].tb_lineno)) 
    # Get environments of first project as fisrt user
    rv = self.app.get('/api/v1/project/environment/1', headers = self.headers)
    try: self.assertEqual(rv.status, '404 NOT FOUND')
    except AssertionError as e: self.verificationErrors.append(str(e) + "Line: " + str(sys.exc_info()[2].tb_lineno)) 
    # Put non existant environment of second project as first user
    payload = { "value": value }
    rv = self.app.put('/api/v1/project/environment/1', json = payload, headers = self.headers)
    try: self.assertEqual(rv.status, '404 NOT FOUND')
    except AssertionError as e: self.verificationErrors.append(str(e) + "Line: " + str(sys.exc_info()[2].tb_lineno)) 
    # Try to get environments of the second as first user
    rv = self.app.get('/api/v1/project/environemnt/2', headers = self.headers)
    try: self.assertEqual(rv.status, '404 NOT FOUND')
    except AssertionError as e: self.verificationErrors.append(str(e) + "Line: " + str(sys.exc_info()[2].tb_lineno)) 


##############################
# Executions

class TestAPI_ExecutionAsUser(unittest.TestCase):

  def setUp(self):
    self.verificationErrors = []
    app.app.testing = True
    self.app = app.app.test_client()
    # Creating a subshell and closing the application creates a ResourceWarning
    warnings.simplefilter("ignore", ResourceWarning)
    initialize_test()
    rv = self.app.post('/api/auth', json = { "username": normal_user, "password": normal_password})
    self.headers = {'Content-Type': 'application/json', 'Authorization': "JWT " + json.loads(rv.data.decode("utf-8"))['access_token']}

  def tearDown(self):
    warnings.simplefilter("default", ResourceWarning)
    self.assertEqual([], self.verificationErrors)

      # GET             /api/v1/project/1/executions
      # POST            /api/v1/project/1/new_execution
  def test_050_get_execution_project(self):
    prj_name = "my_prj_name"
    prj_name2 = "my_prj_name2"
    name = "MY_VARIABLE"
    value = "https://github.com/fabrizio2210/web-infrared.git"
    name2 = "MY_VARIABLE2"
    value2 = "MY_BEAUTIFUL_VARIABLE2"
    # Try to get execution of a non-existant project
    rv = self.app.get('/api/v1/project/1/executions', headers = self.headers)
    try: self.assertEqual(rv.status, '404 NOT FOUND')
    except AssertionError as e: self.verificationErrors.append(str(e) + "Line: " + str(sys.exc_info()[2].tb_lineno)) 
    # Try to post execution of first project
    rv = self.app.post('/api/v1/project/1/new_execution', json = {}, headers = self.headers)
    try: self.assertEqual(rv.status, '404 NOT FOUND')
    except AssertionError as e: self.verificationErrors.append(str(e) + "Line: " + str(sys.exc_info()[2].tb_lineno)) 
    # Create first project
    rv = self.app.post("/api/v1/new_project", json = {"name": prj_name}, headers = self.headers)
    try: self.assertEqual(rv.status, '201 CREATED')
    except AssertionError as e: self.verificationErrors.append(str(e) + "Line: " + str(sys.exc_info()[2].tb_lineno)) 
    try: self.assertEqual(json.loads(rv.data.decode("utf-8")), { "name": prj_name, "id": 1})
    except AssertionError as e: self.verificationErrors.append(str(e)) 
    # Put setting of first project
    payload = { "value": value }
    rv = self.app.put('/api/v1/project/1/setting/scm_url', json = payload, headers = self.headers)
    try: self.assertEqual(rv.status, '200 OK')
    except AssertionError as e: self.verificationErrors.append(str(e) + "Line: " + str(sys.exc_info()[2].tb_lineno)) 
    # Get all executions of the project
    rv = self.app.get('/api/v1/project/1/executions', headers = self.headers)
    try: self.assertEqual(rv.status, '200 OK')
    except AssertionError as e: self.verificationErrors.append(str(e) + "Line: " + str(sys.exc_info()[2].tb_lineno)) 
    try: self.assertEqual(json.loads(rv.data.decode("utf-8")), { "executions" : [ ] })
    except AssertionError as e: self.verificationErrors.append(str(e) + "Line: " + str(sys.exc_info()[2].tb_lineno)) 
    # Try to get a non-existant execution of the new project
    rv = self.app.get('/api/v1/project/1/execution/1234567891234', headers = self.headers)
    try: self.assertEqual(rv.status, '404 NOT FOUND')
    except AssertionError as e: self.verificationErrors.append(str(e) + "Line: " + str(sys.exc_info()[2].tb_lineno)) 
    # Post execution of first project
    rv = self.app.post('/api/v1/project/1/new_execution', json = {}, headers = self.headers)
    try: self.assertEqual(rv.status, '201 CREATED')
    except AssertionError as e: self.verificationErrors.append(str(e) + "Line: " + str(sys.exc_info()[2].tb_lineno)) 
    response = json.loads(rv.data.decode("utf-8"))
    try: self.assertTrue(isinstance(response['execution']['id'], int))
    except AssertionError as e: self.verificationErrors.append(str(e) + "Line: " + str(sys.exc_info()[2].tb_lineno)) 
    try: self.assertTrue(response['execution']['project_id'] == 1)
    except AssertionError as e: self.verificationErrors.append(str(e) + "Line: " + str(sys.exc_info()[2].tb_lineno)) 
    try: self.assertTrue(response['execution']['settings'])
    except AssertionError as e: self.verificationErrors.append(str(e) + "Line: " + str(sys.exc_info()[2].tb_lineno)) 
    try: self.assertTrue(isinstance(response['execution']['commandline'], str))
    except AssertionError as e: self.verificationErrors.append(str(e) + "Line: " + str(sys.exc_info()[2].tb_lineno)) 
    try: self.assertTrue(response['execution']['rc'] is None)
    except AssertionError as e: self.verificationErrors.append(str(e) + "Line: " + str(sys.exc_info()[2].tb_lineno)) 
    try: self.assertTrue(isinstance(response['execution']['start_time'], int))
    except AssertionError as e: self.verificationErrors.append(str(e) + "Line: " + str(sys.exc_info()[2].tb_lineno)) 
    try: self.assertTrue(response['execution']['stop_time'] is None)
    except AssertionError as e: self.verificationErrors.append(str(e) + "Line: " + str(sys.exc_info()[2].tb_lineno)) 
    execution_id = response['execution']['id']
    # Get a execution of the new project
    rv = self.app.get('/api/v1/project/1/execution/' + str(execution_id), headers = self.headers)
    try: self.assertEqual(rv.status, '200 OK')
    except AssertionError as e: self.verificationErrors.append(str(e) + "Line: " + str(sys.exc_info()[2].tb_lineno)) 
    response = json.loads(rv.data.decode("utf-8"))
    try: self.assertEqual(response['execution']['id'], execution_id)
    except AssertionError as e: self.verificationErrors.append(str(e) + "Line: " + str(sys.exc_info()[2].tb_lineno)) 
    try: self.assertTrue(response['execution']['project_id'] == 1)
    except AssertionError as e: self.verificationErrors.append(str(e) + "Line: " + str(sys.exc_info()[2].tb_lineno)) 
    try: self.assertTrue(response['execution']['settings'])
    except AssertionError as e: self.verificationErrors.append(str(e) + "Line: " + str(sys.exc_info()[2].tb_lineno)) 
    try: self.assertTrue(isinstance(response['execution']['commandline'], str))
    except AssertionError as e: self.verificationErrors.append(str(e) + "Line: " + str(sys.exc_info()[2].tb_lineno)) 
    try: self.assertTrue(response['execution']['rc'] is None)
    except AssertionError as e: self.verificationErrors.append(str(e) + "Line: " + str(sys.exc_info()[2].tb_lineno)) 
    try: self.assertTrue(isinstance(response['execution']['start_time'], int))
    except AssertionError as e: self.verificationErrors.append(str(e) + "Line: " + str(sys.exc_info()[2].tb_lineno)) 
    try: self.assertTrue(response['execution']['stop_time'] is None)
    except AssertionError as e: self.verificationErrors.append(str(e) + "Line: " + str(sys.exc_info()[2].tb_lineno)) 
    # Get a execution output of the new project
    time.sleep(1)
    first_byte = 0
    last_byte = 200
    payload = "first_requested_byte=" + str(first_byte) + "&last_requested_byte=" + str(last_byte)
    rv = self.app.get('/api/v1/project/1/execution/' + str(execution_id) + "/output?" + payload, headers = self.headers)
    try: self.assertEqual(rv.status, '200 OK')
    except AssertionError as e: self.verificationErrors.append(str(e) + "Line: " + str(sys.exc_info()[2].tb_lineno)) 
    response = json.loads(rv.data.decode("utf-8"))
    try: self.assertEqual(response['output']['first_transmitted_byte'], first_byte)
    except AssertionError as e: self.verificationErrors.append(str(e) + "Line: " + str(sys.exc_info()[2].tb_lineno)) 
    try: self.assertEqual(response['output']['last_transmitted_byte'], last_byte)
    except AssertionError as e: self.verificationErrors.append(str(e) + "Line: " + str(sys.exc_info()[2].tb_lineno)) 
    try: self.assertTrue(isinstance(response['output']['file_bytes'], int))
    except AssertionError as e: self.verificationErrors.append(str(e) + "Line: " + str(sys.exc_info()[2].tb_lineno)) 
    try: self.assertTrue(response['output']['data'])
    except AssertionError as e: self.verificationErrors.append(str(e) + "Line: " + str(sys.exc_info()[2].tb_lineno)) 
    # Delete a execution of the new project
    rv = self.app.delete('/api/v1/project/1/execution/' + str(execution_id), headers = self.headers)
    try: self.assertEqual(rv.status, '200 OK')
    except AssertionError as e: self.verificationErrors.append(str(e) + "Line: " + str(sys.exc_info()[2].tb_lineno)) 
    # Get all executions of the project
    rv = self.app.get('/api/v1/project/1/executions', headers = self.headers)
    try: self.assertEqual(rv.status, '200 OK')
    except AssertionError as e: self.verificationErrors.append(str(e) + "Line: " + str(sys.exc_info()[2].tb_lineno)) 
    try: self.assertEqual(json.loads(rv.data.decode("utf-8")), { "executions" : [ ] })
    except AssertionError as e: self.verificationErrors.append(str(e) + "Line: " + str(sys.exc_info()[2].tb_lineno)) 



if __name__ == '__main__':
  logging.basicConfig(level=logging.INFO)
  unittest.main()
