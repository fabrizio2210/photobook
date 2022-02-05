from flask import Flask, Response
from flask_restful import Api
from flask_cors import CORS
from flask_sse import sse
import logging
import os

from db import db
from resources.photo import Photo, NewPhoto, PhotoList
from resources.uid import Uid
from utility.networking import get_my_ip
from utility.data import bootstrap
from utility.filemanager import FileManager

if __name__ == '__main__' or os.getenv('DEBUG', 0) == '1':
  logging.basicConfig(level=logging.DEBUG, format='%(asctime)s %(message)s', datefmt='%m/%d/%Y %I:%M:%S %p')

# Initiate Flask
app = Flask(__name__)
app.config['PROPAGATE_EXCEPTIONS'] = True
app.config['REDIS_URL'] = os.getenv('REDIS_URL','redis://localhost')
app.register_blueprint(sse, url_prefix='/api/events')
api = Api(app)

# API
api.add_resource(Photo,     '/api/photo/<uuid:id>')
api.add_resource(PhotoList,     '/api/photos')
api.add_resource(NewPhoto,     '/api/new_photo')
api.add_resource(Uid,     '/api/uid')

# Initialise from envrironment variables
db.set_db_url(os.getenv('DB_URL', 'mongodb://root:develop@mongo:27017/'))
db.set_db_name(os.getenv('DB_NAME', 'photobook'))
FileManager.set_upload_folder(os.getenv('STATIC_FILES_PATH', '/tmp'))

# Initialise data
bootstrap(force=False)

if __name__ == '__main__':
  logging.info('Started')
  bootstrap(force=True, dev=True)
  my_ip = get_my_ip()
  # enable CORS
  CORS(app, resources={r'/*': {'origins': '*'}})
  logging.info("Connect to http://{}:5000/".format(my_ip))
  app.run(host="0.0.0.0", port=5000, debug=True)
