import scala.actors._
import Actor._
import com.redis._

class Consumer(channel: String) extends Runnable {
    val conn = new RedisClient
    val subscriber = new Subscriber(channel)

    class Notifier extends Actor {
        def act {
            while (true) {
                receive {
                    case _ =>
                        rpop_loop
                }
            }
        }
    }

    def rpop_loop {
        while (true) {
            val v = conn.rpop(channel)
            v match {
                case None => return
                case message => println(message.getOrElse("<empty message>"))
            }
        }
    }

    def run {
        val notifier = new Notifier
        notifier.start
        subscriber.subscribe(new Channel[Int](notifier))
    }
}