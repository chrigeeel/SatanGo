import requests

url = "https://dashboard.purepings.eu/password"

payload = 'authenticity_token=sEmbJAWHTES9gxnpcPokjKcFDmEPaQjN934MLkJC9EjqN%2BaMLFBPhsp%2BFJXEqgxgYO9QlJlfHRscenxpRpxcJg%3D%3D&password=PurePings2021'
headers = {
  'sec-ch-ua': '" Not;A Brand";v="99", "Google Chrome";v="91", "Chromium";v="91"',
  'sec-ch-ua-mobile': '?0',
  'Upgrade-Insecure-Requests': '1',
  'Content-Type': 'application/x-www-form-urlencoded',
  'User-Agent': 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36',
  'Accept': 'text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9',
  'Sec-Fetch-Site': 'cross-site',
  'Sec-Fetch-Mode': 'navigate',
  'Sec-Fetch-User': '?1',
  'Sec-Fetch-Dest': 'document',
  'Cookie': '_shreyauth_session=nv1ObbtEO3GHveXt0dDZxH84HxEkVPkL5IxqqneuCCwIK5IwAYaX1a9LqBpLsM9WMYEH5aVOi1mYnnm9dJWF8UfzvZSK7UeJVhcyGRSc1%2BHTAPC9q4fQUY0%2FVUZWh1kRe9PjDVStZhcwwbcnyT5PzKYwCJJb5xHeSYcTnMyJeE6ZrSmCkaM8wnjK4jnWoP1KGP5IyhfPbVUSamD5G0RyCCMgUAIU30SEJVuUY21ubq274xZI763Ch%2F47VarvVKJYNWqeHBsWCgQ7tANi1vgNsCiegw%3D%3D--FssyjyXwR8QII6BY--qMQ2Qx2AiELns8grZ9PqNA%3D%3D'
}

response = requests.request("POST", url, headers=headers, data=payload)

print(response.text)