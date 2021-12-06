import uuid
from flask_restful import Resource


class Uid(Resource):

  def get(self):
    return {'uid': uuid.uuid4().int },200

