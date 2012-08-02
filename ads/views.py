# Create your views here.
from django.shortcuts import render_to_response
from django.template import RequestContext
from django.core import serializers
from ads.models import Ad

import datetime
import urllib2
import re
from BeautifulSoup import BeautifulSoup as bs
from cookielib import CookieJar

cj = CookieJar()
opener = urllib2.build_opener(urllib2.HTTPCookieProcessor(cj))


def home(request):

    if request.POST.get('query') :
        #results is an array which contains Ad objets
        results = search_kijiji(request.POST.get('query'))
        centre = calculate_centre(results)
        json_results = serializers.serialize('json' ,results,ensure_ascii=False)
        return render_to_response('ads/home.html',
                {'results' :  results , 'centre_lat' : centre[0] , 'centre_lng' : centre[1] , 'json_results' : json_results},
                context_instance=RequestContext(request))
        

    #defautl response 
    return render_to_response('ads/home.html',
                        context_instance=RequestContext(request) )

def search_kijiji(query):
    query_url='http://montreal.kijiji.ca/f-SearchAdRedirect?isSearchForm=true&Keyword=%s&CatId=37&lang=en' % query
    ret = []

    try : 
        query_soup  = bs(opener.open(query_url).read())
    except : 
        return []


    ad_urls = [] 
    try : 
        for row in query_soup.find('table', id='SNB_Results').findAll('tr', id = re.compile('resultRow.*') ) :
            ad_url = (row.find('a').get('href'))
            if not Ad.objects.filter(url=ad_url) : 
                ad_urls.append(ad_url)
            else : 
                print "[FOUND]",ad_url
                ret.append(Ad.objects.get(url=ad_url))
    except : 
        return []
    
    return ret + extract_data(ad_urls)


def extract_data(ad_urls):
    result = []
    attribute_dict = { 'Date Listed' : 'date' ,
                       'Price' : 'price' ,
                       'Address' : 'address',
                       'Bathrooms (#)' : 'bathrooms' ,
                       'Furnished' : 'furnished' ,
                       'Pet Friendly' : 'pet_friendly' 
                     }

    for ad_url in ad_urls : 
        try: 
            items = dict() #empty dict used to create the db object 
            items['url'] = ad_url
            #get ad data
            ad_soup = bs(opener.open(ad_url).read())
            map_link = ''

            #title 
            title = ad_soup.find('h1',id='preview-local-title').getText().replace('google_ad_section_start','').replace('google_ad_section_end','')
            items['title'] = title
           
            #table data
            for tr in ad_soup.find('table' , id='attributeTable').findAll('tr') :
                for td in tr.findAll('td') :
                    key =  td.getText()
                    if key in attribute_dict : 
                        value = td.findNext('td').getText()
                        if attribute_dict.get(key) == 'date' :
                            items['pub_date'] =  datetime.datetime.strptime(value,"%d-%b-%y")

                        if attribute_dict.get(key) == 'bathrooms' :
                            items['bathrooms'] = float(re.match('(\d+.?\d?) bathroom.*',value).group(1))

                        if attribute_dict.get(key) == 'price' :
                            if value.find('contact') > 0 :
                                items['rent'] = 0
                            else : 
                                items['rent'] = int(float(value[1:].replace(',',''))) 

                        if attribute_dict.get(key) == 'address' :
                            items['address'] = value.replace('View map','')

                        if attribute_dict.get(key) == 'furnished' or attribute_dict.get(key) == 'pet_friendly' :
                            if value == 'No':
                                items[attribute_dict.get(key)] = False
                            else:
                                items[attribute_dict.get(key)]=True
                        
            # map coordinates
            map_url = 'http://montreal.kijiji.ca' + ad_soup.find('a', attrs = { 'class' : 'viewmap-link' } ).get('href')
            map_soup = bs(opener.open(map_url).read())
            for noscript in map_soup.findAll('noscript') :
                if noscript.find('img') :
                    map_link =  noscript.find('img').get('src')
            coords = urllib2.urlparse.parse_qs(urllib2.urlparse.urlparse(map_link).query)
            lat_lng = coords['center'][0].split(',')
            items['lat'] = float(lat_lng[0])
            items['lng'] = float(lat_lng[1])

            print items
            ad = Ad(**items)
            ad.save()
            result.append(ad)
        
        except Exception as e: 
            print "[FAILED]" , ad_url
            print e
            pass #skip to the next one 


    return result

def calculate_centre(res) :
    lat_c = 0
    lng_c = 0
    count = len(res)
    for r in res : 
        lat_c += r.lat
        lng_c += r.lng
    
    return [ 45.5081 , -73.5550 ] 
