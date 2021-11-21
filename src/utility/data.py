from db import db
from models.photo import PhotoModel
from utility.filemanager import FileManager


def bootstrap(force=False, dev=False, quiet=False):
  if not force:
    # verify and skip if it is not necessary
    if db.present():
      return
  FileManager.delete_all_photos()
  db.delete_all()
  PhotoModel.create_table()
