
from redis import Redis

class RedisWrapper():
  client = None

  @classmethod
  def init(cls, url):
    cls.client = Redis.from_url(url)

  @classmethod
  def publish(cls, msg):
    cls.client.publish('sse', msg)

  @classmethod
  def enque_photo(cls, photo_pb):
    cls.client.lpush('in_photos', photo_pb.SerializeToString())
