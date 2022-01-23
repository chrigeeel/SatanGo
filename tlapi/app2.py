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

blck_lo=np.array([0,0,0])
blck_hi=np.array([10,10,10])
perfDict = {}
def addToDict(char):
    perfDict[char] = []
    onlyfiles = [f for f in listdir('./perfCharDone/' + char) if isfile(join('./perfCharDone/' + char, f))]
    for pr in onlyfiles:
        perf = cv2.imread('./perfCharDone/' + char + '/' + pr)
        hsvp=cv2.cvtColor(perf,cv2.COLOR_BGR2HSV)
        blck_lo=np.array([0,0,0])
        blck_hi=np.array([10,10,10])
        maskp=cv2.inRange(hsvp,blck_lo,blck_hi)
        perf[maskp>0]=(255,0,0)
        perfDict[char].append(perf)


strings = '@#ABCDEFGHJKLMNPQRSTUVWXYZ23456789'

i = 0
while i < 34:
    addToDict(strings[i])
    i += 1


pink = [128, 0, 128]
blue = [255, 128, 128]

@app.route('/', methods=['GET'])
def health_check():
    return jsonify({'status': 'I\'m doing fine!'}), 200
    
def getProb(input_image, char):

    lenfinalboss = 0
    i = 0
    inpt = cv2.imread(input_image)
    hsvi=cv2.cvtColor(inpt,cv2.COLOR_BGR2HSV)

    # mask images to black
    maski=cv2.inRange(hsvi,blck_lo,blck_hi)

    # change image to red/blue where black
    inpt[maski>0]=(0,0,255)
    for perf in perfDict[char]:
        # combine/mix images
        added_image = cv2.addWeighted(inpt,0.5,perf,0.5,0)
        Y, X = np.where(np.all(added_image == pink, axis=2))
        Y2, X2 = np.where(np.all(added_image == blue, axis=2))

        len1 = len(Y) + len(X)
        len2 = len(Y2) + len(X2)

        lenfinal = len1 - len2 * 1

        if lenfinal > lenfinalboss:
            lenfinalboss = lenfinal
        i += 1

    return lenfinalboss


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
    cv2.imwrite('./raws/captchaOG' + tmpFolder + '.png', imdef)

    lower = np.array([0,0,0,111])
    upper = np.array([255,255,255,255])

    mask = cv2.inRange(imdef, lower, upper)
    nobg = cv2.bitwise_and(imdef, imdef, mask= mask)

    cv2.imwrite('./' + tmpFolder + '/captcha.png', nobg)

    imdef = nobg

    grayImage = cv2.cvtColor(imdef, cv2.COLOR_BGR2GRAY)
    (thresh, imdefblw) = cv2.threshold(grayImage, 1, 255, cv2.THRESH_BINARY)

    cv2.imwrite('./' + tmpFolder + '/captcha2.png', imdefblw)

    imblw = cv2.imread('./' + tmpFolder + '/captcha2.png')
    imblw = cv2.bitwise_not(imblw)


    # finds all white pixels in image
    white = [0, 0, 0]
    Y, X = np.where(np.all(imblw == white, axis=2))
    Y2 = Y.tolist()
    X2 = X.tolist()
    i = -1
    listWhite = list()
    for cord in X2:
        i += 1
        xcord = cord
        ycord = Y2[i]
        listWhite.append([ycord, xcord])


    # checks opacity of pixels that are white in blw image
    # gets line through checking opacity and creates a listline
    listLine = list()

    i = 0
    for cord in X:
        if imdef[Y[i], cord][3] > 237:
            listLine.append([Y[i], cord])
            listLine.append([Y[i]+1, cord])
            listLine.append([Y[i]-1, cord])
        i+= 1

    # erases line from image
    for cords in listLine:
        imblw[cords[0], cords[1]] = [255, 255, 255]


    whitenp = np.array([255, 255, 255])

    for cords in listWhite:
        allPxl = [imblw[cords[0]+1, cords[1]],imblw[cords[0]-1, cords[1]],imblw[cords[0], cords[1]+1],imblw[cords[0], cords[1]-1], whitenp]
        if (np.diff(np.vstack(allPxl).reshape(len(allPxl),-1),axis=0)==0).all():
            imblw[cords[0], cords[1]] = [255, 255, 255]
    for cords in listLine:
        if np.any(imblw[cords[0]-1, cords[1]] == 0):
            imblw[cords[0], cords[1]] = [0, 0, 0]

        elif np.any(imblw[cords[0]+1, cords[1]] == 0):
            imblw[cords[0], cords[1]] = [0, 0, 0]

    # another conversion to black and white since we used 1, 0, 0 before
    cv2.imwrite('./' + tmpFolder + '/output3.png', imblw)

    # crops final output into six slices
    #img = Image.open('./' + tmpFolder + '/output3.png')
    left = 4
    top = 17
    right = left + 30
    bottom = 83

    #img_cropped_1 = img.crop((left, top, right, bottom))
    img_cropped_1 = imblw[top:bottom, left:right].copy()
    left = 35
    right = left + 30
    #img_cropped_2 = img.crop((left, top, right, bottom))
    img_cropped_2 = imblw[top:bottom, left:right].copy()
    left = 68
    right = left + 30
    #img_cropped_3 = img.crop((left, top, right, bottom))
    img_cropped_3 = imblw[top:bottom, left:right].copy()
    left = 101
    right = left + 30
    #img_cropped_4 = img.crop((left, top, right, bottom))
    img_cropped_4 = imblw[top:bottom, left:right].copy()
    left = 134
    right = left + 30
    #img_cropped_5 = img.crop((left, top, right, bottom))
    img_cropped_5 = imblw[top:bottom, left:right].copy()
    left = 167
    right = left + 30
    #img_cropped_6 = img.crop((left, top, right, bottom))
    img_cropped_6 = imblw[top:bottom, left:right].copy()

    # saves the slices
    #img_cropped_1.save('./' + tmpFolder + '/img_cropped_1.png')
    #img_cropped_2.save('./' + tmpFolder + '/img_cropped_2.png')
    #img_cropped_3.save('./' + tmpFolder + '/img_cropped_3.png')
    #img_cropped_4.save('./' + tmpFolder + '/img_cropped_4.png')
    #img_cropped_5.save('./' + tmpFolder + '/img_cropped_5.png')
    #img_cropped_6.save('./' + tmpFolder + '/img_cropped_6.png')
    cv2.imwrite('./' + tmpFolder + '/img_cropped_1.png', img_cropped_1)
    cv2.imwrite('./' + tmpFolder + '/img_cropped_2.png', img_cropped_2)
    cv2.imwrite('./' + tmpFolder + '/img_cropped_3.png', img_cropped_3)
    cv2.imwrite('./' + tmpFolder + '/img_cropped_4.png', img_cropped_4)
    cv2.imwrite('./' + tmpFolder + '/img_cropped_5.png', img_cropped_5)
    cv2.imwrite('./' + tmpFolder + '/img_cropped_6.png', img_cropped_6)

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
    shutil.rmtree(tmpFolder)
    return sol, tmpFolder


@app.route('/v1/solve', methods=['POST'])
def solver():

    a = time.time() * 1000
    print('New Request!')
    data = request.json
    try:
        b64 = data["b64"]
    except:
        return jsonify({'error': 'bad captcha', 'solution': 'AAAAAA'})

    try:
        sol, tmp = tlCaptcha(data['b64'])
        processingTime = str(math.floor(time.time() * 1000 - a)) + 'ms'
        return jsonify({'solution': sol, 'tmp': tmp, 'processing_time': processingTime}), 200
    except Exception as e:
        print(e)
        return jsonify({'error': 'bad captcha', 'solution': 'AAAAAA'})


app.debug = True
if __name__ == '__main__':
    app.run()