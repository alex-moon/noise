import uuid
import json
import requests
from bs4 import BeautifulSoup

from tweepy.streaming import StreamListener
from tweepy import OAuthHandler
from tweepy import Stream

# Go to http://dev.twitter.com and create an app.
# The consumer key and secret will be generated for you after
consumer_key="CtTL0PExNLbVpauatFpSD4pPL"
consumer_secret="GJqYbfAiHnV9POWHZfdRgiueziv1BGFQ6SJBhOUr04jYiGM7RW"

# After the step above, you will be redirected to your app's page.
# Create an access token under the the "Your access token" section
access_token="511062632-o9TBpczVGGWL9BrYFiIyd4OJtyG0yw9c8vQtHzKr"
access_token_secret="GmQNxY33FiK6N8iGFjGgFLeQr8HOR4yChZnVt55IiOXpN"

class StdOutListener(StreamListener):
    """ A listener handles tweets are the received from the stream.
    This is a basic listener that just prints received tweets to stdout.

    """
    def on_data(self, data):
        data = json.loads(data)
        if 'text' in data:
            text = data['text']
            for word in text.split(' '):
                if word.startswith('http'):
                    response = requests.get(word)
                    soup = BeautifulSoup(response.content)
                    for paragraph in soup.select("p"):
                        with open("/home/moona/work/noise/text/%s" % uuid.uuid4(), 'w') as text:
                            text.write(paragraph.get_text().encode('ascii', 'replace'))
                            print paragraph.get_text()
        return True

    def on_error(self, status):
        print status


l = StdOutListener()
auth = OAuthHandler(consumer_key, consumer_secret)
auth.set_access_token(access_token, access_token_secret)

stream = Stream(auth, l)
stream.filter(track=['the'])
