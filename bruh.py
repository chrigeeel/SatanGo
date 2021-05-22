from selenium import webdriver
import base64

driver = webdriver.Chrome('chromedriver')  

driver.get("https://dashboard.satanbots.com/purchase?password=du-XAQoBk-nx")

print(base64.b64decode(driver.execute_script("return window.__stripe_xid")))