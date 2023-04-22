from db import db
import logging


class PhotoModel(db.Model):
  __tablename__ = "events"

  def __init__(self, **kwargs):
    for arg in kwargs.keys():
      setattr(self, arg, kwargs[arg])

  def public_json(self):
    return {'id': self.id,
            'description': self.description,
            'author': self.author,
            'timestamp': self.timestamp}

  def save_to_db(self):
    self.save()

  def delete_from_db(self):
    self.delete()

  @classmethod
  def get_photos_by_author_id(cls, author_id):
    return cls.find(author_id=author_id)

  @classmethod
  def get_all_photos(cls):
    logging.debug("Retrieving all the photos.")
    return cls.find()

  @classmethod
  def find_by_id(cls, id):
    return cls.find(id=id)

  @classmethod
  def find_by_timestamp(cls, timestamp):
    logging.debug("Retrieving phots since \"%d\" timestamp.", timestamp)
    return cls.greaterThan(timestamp=timestamp)
