import time
import uuid
import io
import os
import json
import werkzeug
import logging
import go.proto.photo_in_pb2
from flask_restful import Resource, reqparse
from models.photo import PhotoModel
from redis import Redis
from utility.filemanager import FileManager
from redis_wrapper import RedisWrapper
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
        FileManager.delete_photo(photo[0].photo_id)
        new_photo = PhotoModel(id=str(uuid.uuid4()),
                           event="deletion",
                           photo_id=photo[0].photo_id,
                           description=photo[0].description,
                           author=photo[0].author,
                           author_id=photo[0].author_id,
                           order=photo[0].order,
                           timestamp=RedisWrapper.get_counter("events_count"))

        new_photo.save_to_db()
        # Notify other clients
        RedisWrapper.publish(json.dumps(
          FileManager.photo_to_client(
            new_photo.public_json()
          ))
        )
        return {'photo': FileManager.photo_to_client(photo[0].json())}, 201
      return {'message': 'Not authorized'}, 403
    return {'message': 'Item not found.'}, 404

  def get(self, id):
    id = str(id)
    photo = PhotoModel.find_by_id(id)
    if photo:
      return {'photo': FileManager.photo_to_client(photo[0].public_json())}, 200
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

        new_photo = PhotoModel(id=str(uuid.uuid4()),
                           photo_id=photo[0].photo_id,
                           event="edit",
                           description=photo[0].description,
                           author=photo[0].author,
                           author_id=photo[0].author_id,
                           order=photo[0].order,
                           timestamp=RedisWrapper.get_counter("events_count"))

        new_photo.save_to_db()
        # Notify other clients
        RedisWrapper.publish(json.dumps(
          FileManager.photo_to_client(
            new_photo.public_json()
          ))
        )
      else:
        return {'message': 'Not authorized'}, 403
    else:
      return {'message': 'Item not found.'}, 404
    return {'photo': FileManager.photo_to_client(photo[0].json())}, 201

class PhotoList(Resource):
  parser = reqparse.RequestParser()
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

    # Create a new photo structure, but not saving in DB
    photo = PhotoModel(id=str(uuid.uuid4()),
                       description=data.get('description', ''),
                       event='creation',
                       photo_id=str(uuid.uuid4()),
                       author=data.get('author', ''),
                       author_id=str(data['author_id']),
                       order=RedisWrapper.get_counter("photos_count"),
                       timestamp=RedisWrapper.get_counter("events_count"))

    # Image processing
    image = Image.open(image_file)
    image = ImageOps.exif_transpose(image)

    # Save full quality photo
    image.save(FileManager.path_to_full_quality_folder(photo.photo_id))

    image.thumbnail((900,600))

    # Save photo on filesystem
    image.save(FileManager.path_to_upload_folder(photo.photo_id))

    # Enque photo to check
    with io.BytesIO() as output:
      image.save(output, format="JPEG")
      logging.info("Author_id:%s" % data['author_id'])
      RedisWrapper.enque_photo(
                                go.proto.photo_in_pb2.PhotoIn(
                                  id=photo.id,
                                  photo_id=photo.photo_id,
                                  photo=output.getvalue(),
                                  description=photo.description,
                                  author_id=photo.author_id,
                                  author=photo.author,
                                  timestamp=photo.timestamp,
                                  order=photo.order,
                                  location=FileManager.location_for_client(
                                      photo.photo_id
                                    )
                                )
                              )

    return FileManager.photo_to_client(photo.json()), 201
