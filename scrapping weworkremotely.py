import requests 
from bs4 import BeautifulSoup 
from bs4 import BeautifulSoup
import pandas as pd

#get the endpoints
url = 'https://weworkremotely.com/'
headers = {'user-agent': 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36'}
r = requests.get(url, headers = headers)
sleeptime = random.randrange(3, 9)
time.sleep(sleeptime)
soup = BeautifulSoup(r.content, 'html.parser')
container = soup.find('div', class_='dropdown-container')
endpoints = []
for item in container:
  endpoints.append(item.find('a')['href'])
endpoints

def jobscrape(endpoint):
  url = f'https://weworkremotely.com/{endpoint}'
  #headers são os headers do request para passar para o servidor
  headers = {'user-agent': 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36'}
  r = requests.get(url, headers = headers)
  # função que pega a url, faz o request nela, passando os headers selecionados
  global soup
  soup = BeautifulSoup(r.content,'html.parser')
  # todo o html da pagina
  # r.content é o html, e o parser é o interpretador
  return soup




import time
import random 
company = []
title = []
jobtype = []
region = []
joburl = []
jtime = []
def scrapendpoints(endpoints):
  for item in endpoints:
    jobscrape(item)
    sleeptime = random.randrange(3, 9)

    card = soup.findAll('li', class_='feature')
    print(len(card))
    time.sleep(sleeptime)
    print(item)
    return card
def scrapitems(card):
  company = []  
  title = []
  jobtype = []
  region = []
  joburl = []
  jtime = []    
  for item in card:
    sleeptime = random.randrange(3, 9)
    title.append(item.find('span', class_='title').text)
   #append title
    time.sleep(sleeptime)
    sleeptime = random.randrange(3, 9)
    joburl.append(item.find('a')['href'])
  
    time.sleep(sleeptime)
    sleeptime = random.randrange(3, 9)
    company.append(item.find('span', class_='company').text)
    time.sleep(sleeptime)
    sleeptime = random.randrange(3, 9)
    jobtype.append(item.find('span', class_='region company'))
    test = item.find('span', class_='featured')
    if test == None:
      test = item.find('span', class_='date').text
      jtime.append(test)
    else: 
      jtime.append(test)
  df = pd.DataFrame({'company':company, 'title':title,"jobtype":jobtype,'region':region, 'joburl':joburl, 'time':jtime})      
  return df



