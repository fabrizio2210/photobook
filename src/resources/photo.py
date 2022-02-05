import time
import os
import werkzeug
import logging
from flask_restful import Resource, reqparse
from flask_sse import sse
from models.photo import PhotoModel
from utility.filemanager import FileManager
from PIL import Image, ImageOps


class Photo(Resource):
  parser = reqparse.RequestParser()
  parser.add_argument('author_id',
                      type=str,
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
    id = str(id)
    data = Photo.parser.parse_args()
    if os.getenv('BLOCK_UPLOAD', False):
      return { 'message':
        os.getenv('BLOCK_UPLOAD_MSG', 'The upload is blocked by admin.')}, 403
    photo = PhotoModel.find_by_id(id)
    if photo:
      if photo[0].author_id == data.get('author_id', None):
        FileManager.delete_photo(photo[0].id)
        photo[0].delete_from_db()
        # Notify other clients
        sse.publish('deleted ' + str(photo[0].id))
        return {'photo': FileManager.photo_to_client(photo[0].json())}, 201
      return {'message': 'Not authorized'}, 403
    return {'message': 'Item not found.'}, 404

  def get(self, id):
    id = str(id)
    photo = PhotoModel.find_by_id(id)
    if photo:
      return {'photo': FileManager.photo_to_client(photo[0].json())}, 200
    return {'message': 'Item not found.'}, 404

  def put(self, id):
    id = str(id)
    if os.getenv('BLOCK_UPLOAD', False):
      return { 'message':
        os.getenv('BLOCK_UPLOAD_MSG', 'The upload is blocked by admin.')}, 403
    data = Photo.parser.parse_args()
    photo = PhotoModel.find_by_id(id)
    if photo:
      if photo[0].author_id == data.get('author_id', None):
        if data.get('description') is not None:
          photo[0].description = data.get('description')
        if data.get('author') is not None:
          photo[0].author = data.get('author')
        photo[0].save_to_db()
        # Notify other clients
        sse.publish('changed ' + str(photo[0].id))
      else:
        return {'message': 'Not authorized'}, 403
    else:
      return {'message': 'Item not found.'}, 404
    return {'photo': FileManager.photo_to_client(photo[0].json())}, 201

class PhotoList(Resource):
  parser = reqparse.RequestParser()
  parser.add_argument('timestamp',
                      type=int,
                      required=False,
                      help="If timestamp is provided, get latest photos after timestamp."
                      )
  parser.add_argument('author_id',
                      type=str,
                      required=False,
                      help="If author_id is provided, get all the photos of that author."
                      )
  def get(self):
    data = PhotoList.parser.parse_args()
    if data.get('author_id', None):
      return {'photos': list(map(lambda x:
                             FileManager.photo_to_client(x.json()),
                             PhotoModel.get_photos_by_author_id(data['author_id'])
                             ))
             }
    if data.get('timestamp', None):
      return {'photos': list(map(lambda x: 
                             FileManager.photo_to_client(x.json()),
                             PhotoModel.find_by_timestamp(data['timestamp'])
                         ))
             }
    logging.debug("In the resources: %s", PhotoModel.get_all_photos())
    return {'photos': list(map(lambda x: 
                           FileManager.photo_to_client(x.public_json()),
                           PhotoModel.get_all_photos()
                       ))
           }


class NewPhoto(Resource):
  parser = reqparse.RequestParser()
  parser.add_argument('author_id',
                      type=str,
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
    if os.getenv('BLOCK_UPLOAD', False):
      return { 'message':
        os.getenv('BLOCK_UPLOAD_MSG', 'The upload is blocked by admin.')}, 403
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
    image = ImageOps.exif_transpose(image)

    image.thumbnail((900,600))

    # Save photo on filesystem
    image.save(FileManager.path_to_upload_folder(photo.id))

    # Notify other clients
    sse.publish('new_image')

    return FileManager.photo_to_client(photo.json()), 201
