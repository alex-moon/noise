import scala.actors._
import scala.actors.Actor._
import com.redis._

object Noise {
  def main(args: Array[String]) {
    val r = new RedisClient("localhost", 6379)
    r.set("key", "some value")
    val butt = r.get("key")
    println(butt)
  }
}