import socket
import numpy as np
import pandas as pd
import nltk
import re
from colorama import Fore, Back, Style
print("Start download")
nltk.download('stopwords')
print("Finished downloading stopwords")
nltk.download('wordnet')
print("Finished downloading wordnet")
import torch
import torch.nn as nn
from transformers import DistilBertTokenizer
device = torch.device("cuda" if torch.cuda.is_available() else "cpu")



print("Connecting to socket")
sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
sock.bind(('localhost', 8080))
sock.listen(1)
conn, addr = sock.accept()



def utils_preprocess_text(text, flg_stemm=False, flg_lemm=True, lst_stopwords=None):
    text = re.sub(r"(?i)\b((?:https?://|www\d{0,3}[.]|[a-z0-9.\-]+[.][a-z]{2,4}/)(?:[^\s()<>]+|\(([^\s()<>]+|(\([^\s()<>]+\)))*\))+(?:\(([^\s()<>]+|(\([^\s()<>]+\)))*\)|[^\s`!()\[\]{};:'\".,<>?«»“”‘’]))"
         ,'__link__',string=text)
    text = re.sub(r'[^\w\s]', '', str(text).lower().strip())

    
    
    
    lst_text = text.split()    ## remove Stopwords
    if lst_stopwords is not None:
        lst_text = [word for word in lst_text if word not in 
                    lst_stopwords]
    if flg_stemm == True:
        ps = nltk.stem.porter.PorterStemmer()
        lst_text = [ps.stem(word) for word in lst_text]
    if flg_lemm == True:
        lem = nltk.stem.wordnet.WordNetLemmatizer()
        lst_text = [lem.lemmatize(word) for word in lst_text]
    text = " ".join(lst_text)
    
    
    
def textCleaner(df = None , src = 'comment_text' ,dst = 'text_clean',stop_words = 'english'):
    
    if df is None:
        raise TypeError("Data Frame cannot be type 'None'")
    
    try:
        lst_stopwords = nltk.corpus.stopwords.words(stop_words)
    except:
        raise Exception( "'" + stop_words +"'"+" is not a valid type.")
    df[dst] = df[src].apply(lambda x: 
          utils_preprocess_text(x, flg_stemm=False, flg_lemm=True, 
          lst_stopwords=lst_stopwords))

    return text

pretrained = "distilbert-base-uncased"
print("starting tokenizer")
tokenizer = DistilBertTokenizer.from_pretrained(pretrained)
print("finished tokenizer")
maxlen = 2000

class LSTMModel(nn.Module):
    
    def __init__(self,n_vocab,d_model,maxlen,bidirectional = True):
        
        super().__init__()
        self.d_model = d_model
        self.times = (2 if bidirectional  else 1)
        
        
        
        self.embedding = nn.Embedding(n_vocab,self.d_model,)
        self.lstm = nn.LSTM(self.d_model,self.d_model,1,batch_first = True,
                          bidirectional = bidirectional)
        self.attention =  nn.Linear(self.d_model * self.times,1)
        self.softmax = nn.Softmax(dim=1)
        self.classifier = nn.Sequential(
            nn.Linear(self.times * d_model,256),
            nn.LayerNorm(256),
            nn.Dropout(.15),
            nn.ReLU()    ,
            nn.Linear(256,128),
            nn.LayerNorm(128),
            nn.Dropout(.15),
            nn.ReLU(),
            nn.Linear(128,2))
    
    
    def _generate_initial_hidden(self,batch_size):
        return (torch.zeros((self.times,batch_size,self.d_model ),device = device,
                           dtype=torch.float32),
                torch.zeros((self.times,batch_size,self.d_model ),device = device,
                           dtype=torch.float32))

    def forward(self,x,mask):
        embeds = self.embedding(x)
        batch_size = embeds.shape[0]
        h0,c0 = self._generate_initial_hidden(batch_size)
        h,c = self.lstm(embeds,(h0,c0))
        mask = mask.permute(0,2,1)
        attentions = self.softmax(self.attention(h * mask))
        h = attentions * h
        outs = self.classifier(h[:,-1,:])
        #outs = self.classifier(torch.rand((batch_size,h.shape[-1])))
        return outs,attentions
        
model = LSTMModel(n_vocab = tokenizer.vocab_size,
              d_model= 256,
              maxlen = maxlen,
              bidirectional = False).to(device)
              
              
def attention2word(tokenizer,tokens,attention,top = 5):
    items = attention[0,:,0].topk(top).indices.detach().cpu()
    resp_words = [tokenizer.ids_to_tokens[tokens[0][item].detach().cpu().item()] for item in items]
    return resp_words
    
class WordHighlight():
    
    def __init__(self,highlight = True):
        
        self.highlight = highlight
        
    def __call__(self,original_txt,res):    
        color = Back.RED
        out_sent = ''
        
        if type(original_txt) != str:
            original_txt = str(original_txt).split()
        
        
        for word in original_txt.split():

            checker = sum([1 if resp_word in word else 0 for resp_word in res])
            if checker:
                if self.highlight:
                    out_sent += color+  word + Style.RESET_ALL+' '
                else:
                    out_sent += '*** '
                    
            else:
                out_sent += word + ' '

        return out_sent

model.load_state_dict(torch.load("params.pth"))





def process_data(data):
    t = None
    df = pd.DataFrame({"is":[0],"comment_text":data})
    textCleaner(df = df,stop_words = 'english',src='comment_text')

    sent = df.text_clean[0]

    batch_test = tokenizer(sent,max_length=maxlen,padding='max_length',
                          return_tensors="pt").to(device)

    tokens,mask = batch_test['input_ids'],batch_test['attention_mask'].unsqueeze(0)
    model.eval()

    with torch.no_grad():
        pred,att = model(tokens,mask)

    tokens = tokens.detach().cpu()
    label = torch.argmax(pred)
    if label == 0:
        t = "safe"
    else:
        t = "not  safe"
    
    #resp_words = attention2word(tokenizer,tokens,att,top = 5)
    #a = WordHighlight(highlight=True)
    #print(a(sentence,resp_words))
    
    return t
    
while True:
    data = conn.recv(1024)
    if not data:
        pass
    else: 
        label = process_data(data)
        send_data = label.encode()
        conn.send(send_data)
    

conn.close()
