from db import db


class PhotoModel(db.Model):
  id = db.Field(_type = "integer", primary_key =  True)
  description = db.Field(_type = "string")
  author = db.Field(_type = "string")
  author_id = db.Field(_type = "string")
  timestamp = db.Field(_type = "integer")
  __tablename__ = "photos"

  def __init__(self, id, description, author_id, author, timestamp):
    self.description = description
    self.author = author
    self.author_id = author_id
    self.timestamp = timestamp
    self.id = id

  def json(self):
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
    return cls.find()

  @classmethod
  def find_by_id(cls, id):
    return cls.find(id=id)


  @classmethod
  def find_by_timestamp(cls, timestamp):
    return cls.greaterThan(timestamp=timestamp)
