import scala.actors._
import com.redis._

class Subscriber(channel: String) {
    val conn = new RedisClient

    def subscribe(notifier: Channel[Int]) {
        conn.subscribe(channel) { received =>
            received match {
                case M(channel, _) =>
                    notifier ! 1
                //case E(exception) =>  // unrecognised type E!?
                //    ()  // could do something here I don't know
                case other =>  // subscribe, unsubscribe, etc.
                    println("SUBSCRIBER - received something that is not a message " + other)
            }
        }
    }


}