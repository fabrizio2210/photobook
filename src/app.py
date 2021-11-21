from flask import Flask, Response
from flask_restful import Api
from flask_cors import CORS
import logging
import os

from db import db
from resources.photo import Photo, NewPhoto, PhotoList
from utility.networking import get_my_ip
from utility.data import bootstrap
from utility.messageannouncer import MAInstance


# Initialise from envrironment variables
db.set_db_filename(os.getenv('DB_PATH', '/tmp/data.db'))

# Initiate Pub/Sub
announcer = MAInstance.message_announcer

# Initiate Flask
app = Flask(__name__)
app.config['PROPAGATE_EXCEPTIONS'] = True
api = Api(app)

# Server-Sent Events
@app.route('/api/events', methods=['GET'])
def listen():
  def stream():
    messages = announcer.listen()  # returns a queue.Queue
    while True:
      msg = messages.get()  # blocks until a new message arrives
      yield msg
  return Response(stream(), mimetype='text/event-stream')


# API
api.add_resource(Photo,     '/api/photo/<int:id>')
api.add_resource(PhotoList,     '/api/photos')
api.add_resource(NewPhoto,     '/api/new_photo')

if __name__ == '__main__':
  logging.basicConfig(level=logging.DEBUG, format='%(asctime)s %(message)s', datefmt='%m/%d/%Y %I:%M:%S %p')
  logging.info('Started')
  bootstrap(force=True, dev=True)
  my_ip = get_my_ip()
  # enable CORS
  CORS(app, resources={r'/*': {'origins': '*'}})
  logging.info("Connect to http://{}:5000/".format(my_ip))
  app.run(host="0.0.0.0", port=5000, debug=True)
