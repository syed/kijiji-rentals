
function objectifyForm(formArray) {//serialize data function

    var returnArray = {};
    for (var i = 0; i < formArray.length; i++){
        returnArray[formArray[i]['name']] = formArray[i]['value'];
    }
    return returnArray;
}

$(function () {
    var form = $("#searchform");


    form.submit(function (e) {
        //prevent double submit
        //make spinner visible and hide the search button
        $("#submitbutton input").css("visibility", "hidden");
        $("#submitbutton i").css("visibility", "visible");

        var query_data = objectifyForm(form.serializeArray());

        query_data.posted_after = Date.parse(query_data.posted_after);

        bounds  = map.getBounds();
        ne = bounds.getNorthEast();
        sw = bounds.getSouthWest();

        query_data.bounds = {};
        query_data.bounds.ne = {};
        query_data.bounds.sw = {};

        query_data.bounds.ne.lat = ne.lat();
        query_data.bounds.ne.lng = ne.lng();

        query_data.bounds.sw.lat = sw.lat();
        query_data.bounds.sw.lng = sw.lng();

        console.log(query_data);

        $.getJSON(form.attr('action'), //url
            query_data,//post data
            function (responseText, responseStatus) {
                console.log(responseText);

                $("#submitbutton input").css("visibility", "visible");
                $("#submitbutton i").css("visibility", "hidden");

                addMarkers(responseText);
            }
        );
        // prevent actual submit from happening
        e.preventDefault();
    });
});

$(function () {
    $("#datepicker").datepicker({minDate: -20, maxDate: "+1M +10D"});
});