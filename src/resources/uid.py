import uuid
from flask_restful import Resource


class Uid(Resource):

  def get(self):
    return {'uid': str(uuid.uuid4()) },200

