My entry to Telegram Data Clustering content.
https://entry1144-dcround1.usercontent.dev/categories/en/

Basically, my first attempt to write something of use in GoLang. Here's some things I implemented and general approach.

 - FastText integration. I've used FastText C++ library for model construction and inference, so I've CC'ed C wrapper for FastText and extended it to my needs, so I could call it from GoLang. It works smoothly. 
 - For language detection I've used existing fasttext model.
 - For news/non-news classification I've looked thorght the list of domain names and cherrypicked those into 2 categories - definitely trustworthy and definitly spam/fake. Then I assinged appropriate classes to all articles coming from those domain and used them exclusively to train classifier.
 - For categories detection I've made labels for a few (5k) english articles using Google Text Classification API and manually created a conversion rules for G.Categories into Telegram categories. Using those I trained English categories model. As for Russian, I've cheated even twice - I've translated 5k or Russian texts into Engish with Google Translation API, assigned them categories using model created for English texts and trained Russian model using those pseudo-labels, that had confidence over significant threshold.
 
I didn't had time to do last parts of the contenst - Top and Threads and my few starting experiments on threads were of low quality results. Naive sentence embeddings and cosine distances produced low-quality results, I've spent too much time on implementing LSH, as I was targeting super-performance to practice a bit. Surely, different approach was required here.
