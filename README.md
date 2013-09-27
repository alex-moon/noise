Codename "Noise"
==============

**Officially**: a generator of a new kind of metadata (roughly speaking this "sentiment analysis" malarchy)

**Off the Record**: an excuse to learn `erlang`, `rust` and `go`


Typical user interaction:


    Noise: There seems to be a lot of noise relating to these keywords:
           "miley cyrus" "twerk" "rape culture"
           Do you know what we are talking about?
           Yes / No
    User:  Yes
    Noise: Tell us a bit about it.
    User:  And so on and so forth.


Behind the (thin) web app layer (`twisted` perhaps?) that takes input and sends output, the app's split into three modules:

2. A word-by-word analyser that links 1-3grams with keywords (noSQL could work here? Some kind of distributed database...)
3. A fulltext analyser that feeds back into the system
4. The system: we're gonna scour the Web for trending keywords
   - raw counts
   - weighted by length of time in top 100 (1000? 10000?)
     - this should account for stopwords too
   - in correlations (need to look up some statistical analysis algorithms)


Each of these will be written in a different language.


Will it work? Only one way to find out.
