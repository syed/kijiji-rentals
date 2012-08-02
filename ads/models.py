from django.db import models

class Ad(models.Model) :
    BEDROOM_CHOICES = (
        ('1BR' , 'One Bedroom'),
        ('2BR' , 'Two Bedroom'),
        ('3BR' , 'Three Bedroom'),
        ('4BR' , 'Four Bedroom'),
    )

    pub_date = models.DateField()
    rent = models.IntegerField(default=0)

    lat = models.FloatField(default=0)
    lng = models.FloatField(default=0)

    rooms = models.CharField(choices=BEDROOM_CHOICES,max_length=20)

    furnished = models.BooleanField(default=False)
    pet_friendly = models.BooleanField(default=False)
    content =  models.TextField(default='')
    
    bathrooms = models.FloatField(default=0)
    bedrooms = models.FloatField(default=0)
    
    title =  models.CharField(max_length=1000,default='')
    address =  models.CharField(max_length=1000,default='')
    url =  models.CharField(max_length=1000,default='')

    def __unicode__(self):
        return str(self.url)

    

