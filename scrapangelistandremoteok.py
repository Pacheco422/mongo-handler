import re
from selenium import webdriver
from bs4 import BeautifulSoup
import time
from selenium.webdriver.common.keys import Keys
from selenium.webdriver.chrome.options import Options
from selenium.webdriver.chrome.service import Service
#from webdriver_manager.chrome import ChromeDriverManager
 

def driversetup():
    options = webdriver.ChromeOptions()
    #run Selenium in headless mode
    options.add_argument('--headless')
    options.add_argument('--no-sandbox')
    #overcome limited resource problems
    options.add_argument('--disable-dev-shm-usage')
    options.add_argument("lang=en")
    #open Browser in maximized mode
    options.add_argument("start-maximized")
    #disable infobars
    options.add_argument("disable-infobars")
    #disable extension
    options.add_argument("--disable-extensions")
    options.add_argument("--incognito")
    options.add_argument("--disable-blink-features=AutomationControlled")
    
    driver = webdriver.Chrome("C:\webdrivers\chromedriver.exe")
    #webdriver.Chrome("/usr/bin/chromium-browser", options=options)

    driver.execute_script("Object.defineProperty(navigator, 'webdriver', {get: () => undefined});")

    return driver


def pagesource(url, driver):
    driver = driver
    driver.get(url)
    print('cheguei', url)
    soup_len = BeautifulSoup(driver.page_source)
    text = soup_len.find('div', class_='action-remove-latest-filter').text
    match = re.search(r'\d+', text)
    number = match.group()
    number = int(number)
    roll = number/20
    roll = int(roll)
    for item in range(roll):
      driver.execute_script("window.scrollTo(0, document.body.scrollHeight);")
      print('cheguei')
      time.sleep(5)

    global soupa
    soupa = BeautifulSoup(driver.page_source)
    driver.close()
    return soupa 


