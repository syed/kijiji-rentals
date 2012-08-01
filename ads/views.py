# Create your views here.
from django.shortcuts import render_to_response
from django.template import RequestContext

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
        return render_to_response('ads/home.html',
                {'results' :  results},
                context_instance=RequestContext(request))
        

    #defautl response 
    return render_to_response('ads/home.html',
                        context_instance=RequestContext(request) )

def search_kijiji(query):
    query_url='http://montreal.kijiji.ca/f-SearchAdRedirect?isSearchForm=true&Keyword=%s&CatId=37&lang=en' % query

    try : 
        query_soup  = bs(opener.open(query_url).read())
    except : 
        return []


    ad_urls = [] 
    try : 
        for row in query_soup.find('table', id='SNB_Results').findAll('tr', id = re.compile('resultRow.*') ) :
            ad_urls.append(row.find('a').get('href'))
    except : 
        return []

   
    result = []
    attribute_dict = { 'Date Listed' : 'date' ,
                       'Price' : 'price' ,
                       'Address' : 'address',
                       'Bathrooms (#)' : 'bathrooms' ,
                       'Furnished' : 'furnished' ,
                       'Pet Friendly' : 'pet_friendly' 
                     }

    for ad_url in ad_urls : 
        print "url" , ad_url 
        try: 
            #get ad data
            ad_soup = bs(opener.open(ad_url).read())
            map_link = ''

            #title 
            title = ad_soup.find('h1',id='preview-local-title').getText()
            print "title: " , title
            
            for tr in ad_soup.find('table' , id='attributeTable').findAll('tr') :
                for td in tr.findAll('td') :
                    key =  td.getText()
                    if key in attribute_dict : 
                        value = td.findNext('td').getText()
                        print key , ":" ,  value

            # map coordinates
            map_url = 'http://montreal.kijiji.ca' + ad_soup.find('a', attrs = { 'class' : 'viewmap-link' } ).get('href')
            map_soup = bs(opener.open(map_url).read())
            for noscript in map_soup.findAll('noscript') :
                if noscript.find('img') :
                    map_link =  noscript.find('img').get('src')
            print map_link
            coords = urllib2.urlparse.parse_qs(urllib2.urlparse.urlparse(map_link).query)
            print coords['center'][0].split(',')
        
        except: 
            pass #skip to the next one 


    return result
    
    


    

