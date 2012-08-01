from django.db import models

class Ad(models.Model) :
    BEDROOM_CHOICES = (
        ('1BR' , 'One Bedroom'),
        ('2BR' , 'Two Bedroom'),
        ('3BR' , 'Three Bedroom'),
        ('4BR' , 'Four Bedroom'),
    )

    ad_id = models.IntegerField()
    pub_date = models.DateField()
    rent = models.IntegerField()
    lat = models.FloatField()
    lng = models.FloatField()
    rooms = models.CharField(choices=BEDROOM_CHOICES,max_length=20)
    furnished = models.BooleanField()
    content =  models.TextField()
    url =  models.TextField()
    

