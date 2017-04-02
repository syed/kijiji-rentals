// globals
var curr_info_window = false;
var map = false;

var rentalMarkersArray = [];
var directionsArray = []
var nearbyPlacesMarkersArray = [];

var currentLocation;



function customMarker(latLng) {

    var pinColor = "6495ED"; //blueish

    var pinImage = new google.maps.MarkerImage("http://chart.apis.google.com/chart?chst=d_map_pin_letter&chld=%E2%80%A2|" + pinColor,
    new google.maps.Size(21, 34),
    new google.maps.Point(0,0),
    new google.maps.Point(10, 34));

    var pinShadow = new google.maps.MarkerImage("http://chart.apis.google.com/chart?chst=d_map_pin_shadow",
    new google.maps.Size(40, 37),
    new google.maps.Point(0, 0),
    new google.maps.Point(12, 35));

    return new google.maps.Marker({
            position: latLng,
            map: map,
            icon: pinImage,
            shadow: pinShadow
    });
}

function initialize_map() {
    var mapOptions = {
        center: new google.maps.LatLng(45.5081, -73.5550), //hardcoded to montreal's centre
        zoom: 12,
        mapTypeId: google.maps.MapTypeId.ROADMAP
    };

    map = new google.maps.Map(document.getElementById("map_canvas"),
        mapOptions);
}


function addMarkerToMap(latLng, title, infoWindowContent, custom) {

     var marker = new google.maps.Marker({
        position: latLng,
        title: title
    });

    if (custom) {
        marker = customMarker(latLng)
    }

    marker.setMap(map);
    listenMarker(map, marker, infoWindowContent);
    return marker;
}

function addMarkers(results) {
    deleteOverlays(); //remove the old ones

    for (i = 0; i < results.length; i++) {
        res = results[i];
        if (res.lat === null) {
            continue;
        }
        var mylat = new google.maps.LatLng(res.lat, res.lng);

        // Content string
        var infoData = `<div>
                <input type="hidden" name="lat" value="${res.lat}">
                <input type="hidden" name="lng" value="${res.lng}">
                <b> ${res.title} </b><br>
                rent: $${res.price} <br>
                date: ${res.date_listed} <br>
                link: <a href="${res.url}" target=_blank >Here </a>
                <hr>

                <div >
                  <b> Directions From Here: </b>
                    <a href="#" onclick="addDirection(this)" >
                        <i class="fa fa-plus-circle"></i>
                    </a>
                  <br />
                </div> <hr />
               <div>
                    <b> Places Nearby: </b>
                    <a href="#" onclick="addPlace(this)" >
                        <i class="fa fa-plus-circle"></i>
                    </a>
                <br />
                </div>
            </div>`;

        var marker = addMarkerToMap(mylat, res.price.toString(), infoData, false);
        rentalMarkersArray.push(marker);

    }
    console.log("Plotted " + rentalMarkersArray.length)
}

function listenMarker(map, marker, info_data) {
    var info_window = new google.maps.InfoWindow({content: info_data});
    //click on marker will open popup
    google.maps.event.addListener(marker, 'click', function () {
        if (curr_info_window) {
            curr_info_window.close();
        }
        info_window.open(map, marker);
        //update curr_info_window
        curr_info_window = info_window;
    });
}
function deleteOverlays() {
    if (rentalMarkersArray) {
        for (i in rentalMarkersArray) {
            rentalMarkersArray[i].setMap(null);
        }
        rentalMarkersArray.length = 0;
    }
}

function addDirectionOverlay(srcLatLng, destStr, travelMode) {

    console.log(srcLatLng);
    console.log(destStr);

    var directionsDisplay = new google.maps.DirectionsRenderer({
          map: map
    });

    var request = {
        origin: srcLatLng,
        destination: destStr,
        travelMode: travelMode
    };


    var directionsService = new google.maps.DirectionsService();
    directionsService.route(request, function(response, status) {
      if (status == 'OK') {
        // Display the route on the map.
        directionsDisplay.setDirections(response);
      }
    });
}

function addNearbyPlaces(query, src) {
    placeService = new google.maps.places.PlacesService(map);
    var location = new google.maps.LatLng(src.lat, src.lng);
    var request = {
        location: location,
        radius: '500', //meters
        query: query
    };

    var resultCallback = function(results, status) {

        if (status !== google.maps.places.PlacesServiceStatus.OK) {
            console.error(status);
            return;
        }

        for (var i = 0, result; result = results[i]; i++) {
            console.log(result);
            addMarkerToMap(result.geometry.location, result.name, result.name, true);
        }
    };

    placeService.textSearch(request, resultCallback);
}