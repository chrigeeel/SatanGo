from selenium import webdriver
import base64
import time


driver = webdriver.Chrome('chromedriver')  

driver.get("https://dashboard.satanbots.com/purchase?password=du-XAQoBk-nx")

time.sleep(3)

print(base64.b64decode(driver.execute_script("return window.__stripe_xid")))