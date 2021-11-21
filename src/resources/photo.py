import time
import os
import werkzeug
from flask_restful import Resource, reqparse
from flask_sse import sse
from models.photo import PhotoModel
from utility.filemanager import FileManager
from PIL import Image


class Photo(Resource):
  parser = reqparse.RequestParser()
  parser.add_argument('author_id',
                      type=int,
                      required=True,
                      help="Author identifier."
                      )
  parser.add_argument('author',
                      type=str,
                      required=False,
                      help="Author name."
                      )
  parser.add_argument('description',
                      type=str,
                      required=False,
                      help="Photo description"
                      )

  def delete(self, id):
    photo = PhotoModel.find_by_id(id)
    if photo:
      photo.delete_from_db()
      return {'message': 'Item deleted.'}
    return {'message': 'Item not found.'}, 404

  def put(self, id):
    data = Photo.parser.parse_args()
    photo = PhotoModel.find_by_id(id)
    if photo:
      photo.name = data['name']
    else:
      photo = PhotoModel(None, **data)
    photo.save_to_db()
    return photo.json()

class PhotoList(Resource):
  parser = reqparse.RequestParser()
  parser.add_argument('timestamp',
                      type=int,
                      required=False,
                      help="If timestamp is provided, get latest photos after timestamp."
                      )
  def get(self):
    data = PhotoList.parser.parse_args()
    if data.get('timestamp', None):
      return {'photos': list(map(lambda x: 
                             FileManager.photo_to_client(x.json()),
                             PhotoModel.find_by_timestamp(data['timestamp'])
                         ))
             }
    return {'photos': list(map(lambda x: 
                           FileManager.photo_to_client(x.json()),
                           PhotoModel.get_all_photos()
                       ))
           }


class NewPhoto(Resource):
  parser = reqparse.RequestParser()
  parser.add_argument('author_id',
                      type=int,
                      required=True,
                      help="Author identifier."
                      )
  parser.add_argument('file',
                      type=werkzeug.datastructures.FileStorage,
                      required=True,
                      location='files')
  parser.add_argument('author',
                      type=str,
                      required=False,
                      help="Author name."
                      )
  parser.add_argument('description',
                      type=str,
                      required=False,
                      help="Photo description"
                      )

  def post(self):
    data = NewPhoto.parser.parse_args()
    image_file = data['file']
    if not FileManager.allowed_file(image_file.filename):
      return {"message": "File name extension not allowed."}, 401

    # Create a new photo in DB
    photo = PhotoModel(id=None,
                       description=data.get('description', ''),
                       author=data.get('author', ''),
                       author_id=data['author_id'],
                       timestamp=int(time.time()*1000))
    try:
      photo.save_to_db()
    except Exception as e:
      print(e)
      return {"message": "An error occurred inserting the photo."}, 500

    # Image processing
    image = Image.open(image_file)
    image.thumbnail((900,600))

    # Save photo on filesystem
    image.save(FileManager.path_to_upload_folder(photo.id))

    # Notify other clients
    sse.publish('new_image')

    return FileManager.photo_to_client(photo.json()), 201
