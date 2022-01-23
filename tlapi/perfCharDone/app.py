from flask import Flask, request, jsonify
import cv2
from PIL import Image
import random
import os
from os import listdir
from os.path import isfile, join
import shutil
import base64
import numpy as np
import time
import math
import json
import requests
import re
from discord_webhook import DiscordWebhook, DiscordEmbed
from requests_futures.sessions import FuturesSession

app = Flask(__name__)
global allowedKeys
allowedKeys = {}
with open('parseddata2.json', 'r') as parsed:
    parsedData = json.load(parsed)
    for key in parsedData:
        allowedKeys[key] = 50


for key in allowedKeys:
    with open('./ratelimits/' + key +'.json', 'w') as parsed:
        key1 = {}
        key1[key] = 50
        json.dump(key1, parsed, indent=4)


@app.after_request
def response_processor(response):
    try:
        sol1 = sol
        processingTime1 = processingTime
        data = request.json
        captchaID = data['b64']
        key = data['key']
        username = data['username']
        webhookUrl = data['webhookUrl']
        print(username, key)
        done = True

    except:
        done = False

    @response.call_on_close
    def process_after_request():

        if done == True:
            session = FuturesSession()

            headers = {
            'content-type': 'application/json'
            }

            data1 = json.dumps({'b64': captchaID,
                               'username': username,
                               'solution': sol1,
                               'processingTime': processingTime1,
                               'webhookUrl': webhookUrl
                               })
            
            #r = session.post('http://50.18.189.206:5069/webhook', headers=headers,data=data1)

    return response

@app.route('/', methods=['GET'])
def health_check():
    return jsonify({'status': 'I\'m doing fine!'}), 200


@app.route('/pvr/update/userdata', methods=['POST'])
def updatedata():
    def write_json(data, filename='parseddata2.json'): 
        with open(filename,'w') as f:
            json.dump(data, f, indent=4)

    def parse():
        data = request.json
        data1 = data['data']
        auth = data1[0]['auth']
        if auth == 'TLAPI-T86G-4CRQ-F8SJ-7KZJ':
            del data1[0]
            with open('parseddata2.json', 'w') as parsed:
                json.dump({}, parsed, indent=4)
            with open('parseddata2.json', 'r') as parsed2:
                parsedData = json.load(parsed2)
                for user in data1:
                    s = user['key']
                    parsedData[s] = 50

            newKeys = list()
            for key in parsedData:
                newKeys.append(key)
                try:
                    with open('./ratelimits/' + key +'.json', 'r') as parsed:
                        pass

                except Exception as e:
                    print(e)
                    with open('./ratelimits/' + key +'.json', 'w') as parsed:
                        key1 = {}
                        key1[key] = 50
                        json.dump(key1, parsed, indent=4)

            rateLimits = os.listdir('./ratelimits')
            for key in rateLimits:
                key = key[:-5]
                print(key)
                if key not in newKeys:
                    os.remove('./ratelimits/' + key +'.json')

        else:
            return jsonify({'error': 'wrong tlapi key'}), 401

        write_json(parsedData)
        return jsonify({'success': 'successfully updated keys'}), 200

    s = parse()

    return s


@app.route('/pvr/update/ratelimit', methods=['POST'])
def resetlimit():
    def parse():
        data = request.json
        auth = data['auth']
        amt = data['amount']

        if auth == 'TLAPI-T86G-4CRQ-F8SJ-7KZJ':
            global allowedKeys
            allowedKeys = {}

            with open('parseddata2.json', 'r') as parsed:
                parsedData = json.load(parsed)
                for key in parsedData:
                    allowedKeys[key] = amt + 1

            for key in allowedKeys:
                with open('./ratelimits/' + key +'.json', 'w') as parsed:
                    key1 = {}
                    key1[key] = amt + 1
                    json.dump(key1, parsed, indent=4)

        else:
            return jsonify({'error': 'wrong tlapi key'}), 401

        return jsonify({'success': 'successfully updated ratelimit to ' + str(amt)}), 200

    s = parse()

    return s


@app.route('/v1/solve', methods=['POST'])
def solver():

    def getProb(input_image, char):

        Xboss = []
        Yboss = []
        onlyfiles = [f for f in listdir('./perfCharDone/' + char) if isfile(join('./perfCharDone/' + char, f))]

        for pr in onlyfiles:
            inpt = cv2.imread(input_image)
            perf = cv2.imread('./perfCharDone/' + char + '/' + pr)
            hsvi=cv2.cvtColor(inpt,cv2.COLOR_BGR2HSV)
            hsvp=cv2.cvtColor(perf,cv2.COLOR_BGR2HSV)

            # upper/Lower limits black
            blck_lo=np.array([0,0,0])
            blck_hi=np.array([10,10,10])

            # mask images to black
            maski=cv2.inRange(hsvi,blck_lo,blck_hi)
            maskp=cv2.inRange(hsvp,blck_lo,blck_hi)

            # change image to red/blue where black
            inpt[maski>0]=(0,0,255)
            perf[maskp>0]=(255,0,0)

            # combine/mix images
            added_image = cv2.addWeighted(inpt,1.0,perf,1.0,0)

            pink = [255, 0, 255]

            # find intersections (pink)
            Y, X = np.where(np.all(added_image == pink, axis=2))

            if len(X) + len(Y) > len(Xboss) + len(Yboss):
                Xboss = X
                Yboss = Y

        return len(Xboss) + len(Yboss)


    def ocr(input_image):
        strings = '@#ABCDEFGHJKLMNPQRSTUVWXYZ23456789'
        highscore = 0
        highchar = ''
        i = 0
        while i < 34:
            score = getProb(input_image, strings[i])
            if score > highscore:
                highscore = score
                highchar = strings[i]
                #print(highchar)
            i += 1
        return highchar


    def splitChars(b64string):

        # gets random name for folder and files
        tmpFolder = str(random.randint(1000000000, 9999999999))

        # creates folder
        os.mkdir(tmpFolder)
        # decodes base64 and saves image
        with open('./' + tmpFolder + '/captcha.png', 'wb') as fl:
            fl.write(base64.b64decode(b64string))


        imdef = cv2.imread('./' + tmpFolder + '/captcha.png', cv2.IMREAD_UNCHANGED)
        cv2.imwrite('./' + tmpFolder + '/captchaOG.png', imdef)

        lower = np.array([0,0,0,111])
        upper = np.array([255,255,255,255])

        mask = cv2.inRange(imdef, lower, upper)
        nobg = cv2.bitwise_and(imdef, imdef, mask= mask)

        cv2.imwrite('./' + tmpFolder + '/captcha.png', nobg)


        imdef = cv2.imread('./' + tmpFolder + '/captcha.png', cv2.IMREAD_UNCHANGED)

        grayImage = cv2.cvtColor(imdef, cv2.COLOR_BGR2GRAY)
        (thresh, imdefblw) = cv2.threshold(grayImage, 1, 255, cv2.THRESH_BINARY)

        cv2.imwrite('./' + tmpFolder + '/captcha2.png', imdefblw)


        imdef = cv2.imread('./' + tmpFolder + '/captcha.png', cv2.IMREAD_UNCHANGED)
        imblw = cv2.imread('./' + tmpFolder + '/captcha2.png')
        imblw = cv2.bitwise_not(imblw)

        # finds all white pixels in image
        white = [0, 0, 0]
        Y, X = np.where(np.all(imblw == white, axis=2))
        Y = Y.tolist()
        X = X.tolist()
        i = -1
        listWhite = list()
        for cord in X:
            i += 1
            xcord = cord
            ycord = Y[i]
            listWhite.append([ycord, xcord])

        # checks opacity of pixels that are white in blw image
        # gets line through checking opacity and creates a listline
        listLine = list()
        listLine2 = list()
        for cords in listWhite:
            if imdef[cords[0], cords[1]][3] > 237:
                listLine.append([cords[0], cords[1]])
                listLine.append([cords[0]+1, cords[1]])
                listLine.append([cords[0]-1, cords[1]])

        # erases line from image
        for cords in listLine:
            imblw[cords[0], cords[1]] = [255, 255, 255]


        whitenp = np.array([255, 255, 255])
        blacknp = np.array([0, 0, 0])


        for cords in listWhite:

            allPxl = [imblw[cords[0]+1, cords[1]],imblw[cords[0]-1, cords[1]],imblw[cords[0], cords[1]+1],imblw[cords[0], cords[1]-1], whitenp]
            lonely = (np.diff(np.vstack(allPxl).reshape(len(allPxl),-1),axis=0)==0).all()
            if lonely == True:
                imblw[cords[0], cords[1]] = [255, 255, 255]



        for cords in listLine:
            if str(imblw[cords[0]-1, cords[1]]) == '[0 0 0]':
                imblw[cords[0], cords[1]] = [0, 0, 0]

            elif str(imblw[cords[0]+1, cords[1]]) == '[0 0 0]':
                imblw[cords[0], cords[1]] = [0, 0, 0]


        # another conversion to black and white since we used 1, 0, 0 before
        cv2.imwrite('./' + tmpFolder + '/output3.png', imblw)

        # crops final output into six slices
        img = Image.open('./' + tmpFolder + '/output3.png')
        left = 4
        top = 17
        right = left + 30
        bottom = 83

        img_cropped_1 = img.crop((left, top, right, bottom))
        left = 35
        right = left + 30
        img_cropped_2 = img.crop((left, top, right, bottom))
        left = 68
        right = left + 30
        img_cropped_3 = img.crop((left, top, right, bottom))
        left = 101
        right = left + 30
        img_cropped_4 = img.crop((left, top, right, bottom))
        left = 134
        right = left + 30
        img_cropped_5 = img.crop((left, top, right, bottom))
        left = 167
        right = left + 30
        img_cropped_6 = img.crop((left, top, right, bottom))

        # saves the slices
        img_cropped_1.save('./' + tmpFolder + '/img_cropped_1.png')
        img_cropped_2.save('./' + tmpFolder + '/img_cropped_2.png')
        img_cropped_3.save('./' + tmpFolder + '/img_cropped_3.png')
        img_cropped_4.save('./' + tmpFolder + '/img_cropped_4.png')
        img_cropped_5.save('./' + tmpFolder + '/img_cropped_5.png')
        img_cropped_6.save('./' + tmpFolder + '/img_cropped_6.png')
        

        return tmpFolder


    def align(input_image, tmpFolder):

        # creates new image for aligned character
        output_img = 255 * np.ones((50, 30, 3), dtype=np.uint8)
        im = cv2.imread(input_image)

        black = [0, 0, 0]

        # find all black pixes (pixels in character)
        Y, X = np.where(np.all(im == black, axis=2))
        Y = Y.tolist()
        X = X.tolist()
        i = -1
        listBlack = list()
        if 'img_cropped_1' in input_image:
            for cord in X:
                i += 1
                listBlack.append([cord+2, Y[i]-1])
        else:
            for cord in X:
                i += 1
                listBlack.append([cord, Y[i]])

        # since 90 is a bit above the middle, get difference to 90
        moveMargin = 19-listBlack[0][1]

        # draws character on different position on empty canvas output_img
        for cords in listBlack:
            output_img[cords[1] + moveMargin, cords[0]] = [0, 0, 0]

        cv2.imwrite(input_image, output_img)


    def alignToImage(tmpFolder):

        i = 0
        while i < 6:
            i += 1
            inp = './' + tmpFolder + '/img_cropped_' + str(i) + '.png'
            align(inp, tmpFolder)


    def tlCaptcha(captchaID):
        tmpFolder = splitChars(captchaID)
        alignToImage(tmpFolder)
        sol = (ocr('./' + tmpFolder + '/img_cropped_1.png') +
        (ocr('./' + tmpFolder + '/img_cropped_2.png')) +
        (ocr('./' + tmpFolder + '/img_cropped_3.png')) +
        (ocr('./' + tmpFolder + '/img_cropped_4.png')) +
        (ocr('./' + tmpFolder + '/img_cropped_5.png')) +
        (ocr('./' + tmpFolder + '/img_cropped_6.png')))
        sol = sol.replace('Q','?')
        #shutil.rmtree(tmpFolder)
        return sol, tmpFolder


    def rateLimit(key):
        global allowedKeys
        try:
            with open('./ratelimits/' + key +'.json', 'r') as parsed:
                parsedData = json.load(parsed)
                reqsleft = parsedData[key]
                print(reqsleft)
                parsed.close()
                if reqsleft > -1:
                    with open('./ratelimits/' + key +'.json', 'w') as parsed2:
                        parsedData[key] = parsedData[key]-1
                        json.dump(parsedData, parsed2, indent=4)
                        parsed2.close()
                        return parsedData[key]
                else:
                    return 0

        except Exception as e:
            print(e)
            return 0

    a = time.time() * 1000
    print('New Request!')
    try:
        data = request.json
        captchaID = data['b64']
        key = data['key']
    except:
        try:
            print('Invalid key/Captcha')
            print(captchaID)
        except Exception as e:
            print(e)
        return jsonify({'error': 'you didn\'t provided an invalid key and/or invalid captchaid'}), 400

    allowed = rateLimit(key)
    if allowed > 0:
        global sol
        global processingTime
        try:
            sol, tmp = tlCaptcha(data['b64'])
            processingTime = str(math.floor(time.time() * 1000 - a)) + 'ms'
        except Exception as e:
            print(e)
            return jsonify({'error': 'bad captcha'})
    else:
        print('Keys ratelimit has been reached')
        return jsonify({'error': 'your key\'s rate limit has been reached. please ask staff for reset'}), 401

    print('solution: ' + sol)
    return jsonify({'solution': sol, 'tmp': tmp, 'processing_time': processingTime, 'ratelimit': allowed}), 200



app.debug = True
if __name__ == '__main__':
    app.run()