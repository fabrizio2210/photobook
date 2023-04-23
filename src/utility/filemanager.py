import os
import logging
from models.photo import PhotoModel

ALLOWED_EXTENSIONS = {'png', 'jpg', 'jpeg', 'gif'}

class FileManager():
  upload_folder = '/tmp/'
  static_path_url = '/static/resized/'
  full_quality_folder = '/tmp/'

  @classmethod
  def set_upload_folder(cls, path):
    cls.upload_folder = path.rstrip('/') + '/'
    try:
      os.mkdir(cls.upload_folder)
    except OSError as error:
     logging.info(error)

  @classmethod
  def set_full_quality_folder(cls, path):
    cls.full_quality_folder = path.rstrip('/') + '/'
    try:
      os.mkdir(cls.full_quality_folder)
    except OSError as error:
      logging.info(error)

  @classmethod
  def allowed_file(cls, filename):
      return '.' in filename and \
             filename.rsplit('.', 1)[1].lower() in ALLOWED_EXTENSIONS

  @classmethod
  def path_to_upload_folder(cls, id):
    return os.path.join(cls.upload_folder, cls.get_file_name(id))

  @classmethod
  def path_to_full_quality_folder(cls, id):
    return os.path.join(cls.full_quality_folder, cls.get_file_name(id))

  @classmethod
  def get_file_name(cls, id):
    return str(id) + ".jpg"

  @classmethod
  def delete_all_photos(cls):
    photos = PhotoModel.get_all_photos()
    for photo in photos:
      cls.delete_photo(photo.id)

  @classmethod
  def delete_photo(cls, id):
    file_path = os.path.join(cls.upload_folder, cls.get_file_name(id))
    if os.path.isfile(file_path):
      logging.debug("Deleting file: %s", file_path)
      os.unlink(file_path)
    file_path = os.path.join(cls.full_quality_folder, cls.get_file_name(id))
    if os.path.isfile(file_path):
      logging.debug("Deleting file: %s", file_path)
      os.unlink(file_path)

  @classmethod
  def photo_to_client(cls, photo):
    photo.update({
      'location': cls.static_path_url + cls.get_file_name(photo['photo_id']),
    })
    return photo
