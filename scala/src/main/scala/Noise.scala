import com.redis._

object Noise {
  def main(args: Array[String]) {
    val toput = Prover.hello
    val r = new RedisClient("localhost", 6379)
    r.set("key", toput)
    val butt = r.get("key")
    println(toput + " " + butt)
  }
}
