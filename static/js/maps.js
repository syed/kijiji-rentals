// globals
var curr_info_window = false;
var map = false;
var markersArray = [];
function initialize_map() {
    var mapOptions = {
        center: new google.maps.LatLng(45.5081, -73.5550), //hardcoded to montreal's centre
        zoom: 12,
        mapTypeId: google.maps.MapTypeId.ROADMAP
    };

    var directionsService = new google.maps.DirectionsService;
    var directionsDisplay = new google.maps.DirectionsRenderer;


    map = new google.maps.Map(document.getElementById("map_canvas"),
        mapOptions);
}

function addMarkers(results) {
    deleteOverlays(); //remove the old ones
    for (i = 0; i < results.length; i++) {
        res = results[i];
        if (res.map_location == null) {
            continue;
        }
        var mylat = new google.maps.LatLng(res.map_location.lat, res.map_location.lng);

        var marker = new google.maps.Marker({
            position: mylat,
            title: res.price.toString()
        });

        marker.setMap(map);
        markersArray.push(marker);

        // Content string
        var info_data = `<b> ${res.title} </b><br>
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
                </div>`;

        listenMarker(map, marker, info_data);
    }
    console.log("Plotted " + markersArray.length)
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
    if (markersArray) {
        for (i in markersArray) {
            markersArray[i].setMap(null);
        }
        markersArray.length = 0;
    }
}
