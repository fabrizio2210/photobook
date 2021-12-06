import os
from models.photo import PhotoModel

ALLOWED_EXTENSIONS = {'png', 'jpg', 'jpeg', 'gif'}

class FileManager():
  upload_folder = '/tmp/'
  static_path_url = '/static/'

  @classmethod
  def set_upload_folder(cls, path):
    cls.upload_folder = path.rstrip('/') + '/'

  @classmethod
  def allowed_file(cls, filename):
      return '.' in filename and \
             filename.rsplit('.', 1)[1].lower() in ALLOWED_EXTENSIONS

  @classmethod
  def path_to_upload_folder(cls, id):
    return os.path.join(cls.upload_folder, cls.get_file_name(id))

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
    if os.path.isfile(os.path.join(cls.upload_folder, cls.get_file_name(id))):
      os.unlink(os.path.join(cls.upload_folder, cls.get_file_name(id)))

  @classmethod
  def photo_to_client(cls, photo):
    return {
      'author': photo['author'],
      'description': photo['description'],
      'id': photo['id'],
      'location': cls.static_path_url + cls.get_file_name(photo['id']),
      'timestamp': photo['timestamp']
      }

