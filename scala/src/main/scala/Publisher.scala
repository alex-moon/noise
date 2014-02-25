import com.redis._

class Publisher(channel: String) {
  val conn = new RedisClient

  def publish(value: String) {
    conn.lpush(channel, value)
    conn.publish(channel, "1")
  }
}