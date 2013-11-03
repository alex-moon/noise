object Noise {
  def main(args: Array[String]) {
    val consumer = new Thread(new Consumer("noise"))
    consumer.start

    val messages = List(
        "PLAYER THREE HAS JOINED THE GAME",
        "Scala is now speaking to go and erlang",
        "We're going to do four messages of course just to keep it you know consistent",
        "More like Kayfabe Lincoln lol"
    )

    val publisher = new Publisher("noise")
    messages.foreach(message => publisher.publish(message))

    while(true){}
  }
}
