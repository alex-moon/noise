import uuid
import requests
from bs4 import BeautifulSoup

ps = []
for i in range(2,11):
    response = requests.get("http://gasto.cc/new-life-%d" % i)
    soup = BeautifulSoup(response.content)
    for paragraph in soup.select("div.content-body p"):
        with open("/home/moona/work/noise/text/%s" % uuid.uuid4(), 'w') as text:
            text.write(paragraph.get_text().encode('ascii', 'replace'))
