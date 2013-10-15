Codename "Noise"
==============

**Officially**: a generator of a new kind of metadata (roughly speaking this "sentiment analysis" malarchy)

**Off the Record**: an excuse to learn `erlang`, `rust` and `go`


Mission Statement
-----------------

As waste by-products have been put to use and made saleable in heavy industry, so too do we believe:

> <big>**Even the trash of Social Media opinion is useful.**</big>

Typical user interaction:


    Noise: There seems to be a lot of noise relating to these keywords:
           "miley cyrus" "twerk" "rape culture"
           Do you know what we are talking about?
           Yes / No
    User:  Yes
    Noise: Tell us a bit about it.
    User:  And so on and so forth.


Behind the (thin) Web app layer (`twisted` perhaps?) that takes input and sends output, the app's split into three modules. The precise responsibilities of these modules is, I think, still to be determined, but a rough description of the app runs thus: 


We want to correlate ngrams in a feedback loop which builds a deep semantic lexicon based on social media (perhaps the Web in general). It seems reasonable to expect that underneath the noise there is a voice, the voice of whatever human universal the Web can express - the idea being that the voice being the voice of a machine, it takes a machine to hear it.


Does this make sense? Three modules, connected by `redis` - preferably with developer-friendly web APIs, at any rate accessible from `twisted` (or whatever we'll use for our Web app). Will it work? Only one way to find out.
